package driversvc

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/cache"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/clients"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/config"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/mqttclient"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/transform"
	"go.uber.org/atomic"
	"google.golang.org/grpc/status"
)

type BaseService struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	logger      commons.TedgeLogger
	gwInfo      commons.GatewayInfo
	devCache    *cache.DeviceCache
	resourceCli *clients.ResourceClient
	MqttClient  *mqttclient.MqttClient

	options       commons.Options
	configMgr     *ConfigManger
	ossManager    *ossManager
	appCliManager *appCliManager
	libraryId     string
	cloudStatus   atomic.Bool
	readyChan     chan struct{}
}

func NewBaseService(l commons.TedgeLogger) *BaseService {
	var lc commons.TedgeLogger = l
	if lc == nil {
		driverName := "tedge-driver"
		lc = commons.DefaultLogger(commons.InfoLevel, driverName)
	}
	lc.Infof("NewBaseService tedge-driver-sdk-go sdk version:%s", commons.SDKVersion)

	//1. get config
	var configMgr ConfigManger
	flag.StringVar(&configMgr.configFile, "c", common.DefaultConfiguration, "./driver -c /etc/driver/res/configuration.toml")
	flag.Parse()
	err := config.ParseConfig(configMgr.configFile, &configMgr.configuration)
	if err != nil {
		lc.Errorf("NewBaseService read cfg file error:%s", err)
		os.Exit(-1)
	}

	if err = configMgr.configuration.ValidateClientConfig(); err != nil {
		lc.Errorf("NewBaseService validate client config err:%s", err)
		os.Exit(-1)
	}

	//2. new client
	cli, err := clients.NewResourceClient(configMgr.GetClient()[common.Resource])
	if err != nil {
		lc.Errorf("NewBaseService new resource client error:%s", err)
		os.Exit(-1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	bds := &BaseService{
		ctx:     ctx,
		cancel:  cancel,
		wg:      &wg,
		options: commons.DefaultOptions(),

		logger:      lc,
		resourceCli: cli,
		configMgr:   &configMgr,
		devCache:    cache.NewDeviceCache(),
		readyChan:   make(chan struct{}),

		ossManager:    newOssManager(),
		appCliManager: newAppCliManager(lc, cli),
	}

	bds.updateGatewayInfo()
	bds.initDeviceCache()
	go bds.syncingGatewayStatus()
	go bds.waitServerStart()

	return bds
}

////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) waitServerStart() {
	lc := bds.logger
	var i int
	for {
		select {
		case <-bds.ctx.Done():
			return
		case <-time.After(time.Second):
			i++
			lc.Warnf("NewBaseService driver is not running, wait i:%d, Start() shouldn't blocked!", i)
			if i == 10 {
				subCtx, SubCancel := context.WithTimeout(context.Background(), time.Second)
				bds.ReportAlert(subCtx, commons.ERROR, "driver is not running")
				SubCancel()
			}
		case <-bds.readyChan:
			lc.Infof("NewBaseService driver is running, ready !!!")
			return
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) initDeviceCache() {
	var (
		err error
		mdr *proto.MultiDeviceResponse
	)

	c, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancel()

	// sync cid list
	if mdr, err = bds.resourceCli.DevicesSearch(c, &proto.DeviceSearchQueryRequest{
		BaseSearchConditionQuery: &proto.BaseSearchConditionQuery{
			Page:     0,
			PageSize: 0,
		},
		ServiceId: bds.GetServiceId(),
	}); err != nil {
		bds.logger.Errorf("initDeviceCache error: %s", err)
		return
	}

	for _, dev := range mdr.Devices {
		if _, ok := bds.devCache.ById(dev.Id); ok {
			bds.devCache.Update(transform.ToDeviceModel(dev))
		} else {
			bds.devCache.Add(transform.ToDeviceModel(dev))
		}

		devStr, _ := json.Marshal(dev)
		bds.logger.Debugf("initDeviceCache dev:%s", devStr)
	}

	return
}

func (bds *BaseService) initMqttClient() error {
	if bds.options.MqttDriver == nil {
		return nil
	}

	lc := bds.GetLogger()
	lc.Infof("initMqttClient sdk will run with mqtt client")
	mqttServer, ok := bds.configMgr.GetClient()[common.MQTTBroker]
	if !ok {
		lc.Errorf("required mqtt broker config")
		return errors.New("required mqtt broker config")
	}

	username := bds.options.MqttUsername
	var err error
	if bds.MqttClient, err = mqttclient.NewMqttClient(mqttServer.Address, username, lc); err != nil {
		lc.Errorf("new mqtt client error: %s", err)
		return fmt.Errorf("new mqtt client error: %s", err)
	}

	bds.MqttClient.SetDriver(bds.options.MqttDriver)
	mqopts := bds.MqttClient.GetOpts()
	mqopts.SetOnConnectHandler(bds.MqttClient.OnConnectHandler(bds.options.ConnHandler))
	bds.MqttClient.SetClient(mqtt.NewClient(mqopts))
	if err = bds.MqttClient.Connect(); err != nil {
		lc.Errorf("mqtt client connect error: %s", err)
		return err
	}

	bds.wg.Add(1)
	go func() {
		defer bds.wg.Done()
		<-bds.ctx.Done()
		lc.Infof("mqtt client disconnected")
		bds.MqttClient.Disconnect()
	}()

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) syncingGatewayStatus() {
	timerA := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-bds.ctx.Done():
			timerA.Stop()
			return
		case <-timerA.C:
			bds.updateGatewayInfo()
		}
	}
}

func (bds *BaseService) updateGatewayInfo() {
	lc := bds.logger
	timeout, cancelF := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancelF()

	resp, err := bds.resourceCli.GetGatewayInfo(timeout, common.EmptyPb)
	if err != nil {
		lc.Errorf("BaseService GetGatewayInfo error:%s", err)
		bds.SetCloudStatus(false)
		return
	}

	bds.gwInfo = transform.ToGatewayModel(resp)
	bds.SetCloudStatus(bds.gwInfo.CloudState)

	return
}

////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) GetServiceId() string {
	return bds.configMgr.GetService().ID
}

func (bds *BaseService) GetServiceName() string {
	return bds.configMgr.GetService().Name
}

// 获取驱动自定义配置
func (bds *BaseService) GetDriverConfig() map[string]interface{} {
	return bds.configMgr.GetCustomConfig()
}

func (bds *BaseService) GetGatewayId() string {
	return bds.gwInfo.GwId
}

func (bds *BaseService) GetGatewayMeta() commons.GatewayInfo {
	return bds.gwInfo
}

func (bds *BaseService) GetTEdgeModel() commons.RunningModel {
	return bds.gwInfo.Mode
}

//////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) AllDevices() map[string]commons.DeviceInfo {
	return bds.devCache.All()
}

func (bds *BaseService) GetDeviceById(cid string) (commons.DeviceInfo, bool) {
	return bds.devCache.ById(cid)
}

func (bds *BaseService) GetActiveDeviceById(cid string) (commons.DeviceInfo, bool) {
	dev, ok := bds.devCache.ById(cid)
	if !ok {
		return dev, ok
	}
	return dev, dev.ActiveStatus == commons.DeviceActiveStatusActivated
}

func (bds *BaseService) GetActiveDevices() map[string]commons.DeviceInfo {
	var dMap = make(map[string]commons.DeviceInfo)
	for k, v := range bds.devCache.All() {
		if v.ActiveStatus != commons.DeviceActiveStatusActivated {
			continue
		}
		dMap[k] = v
	}
	return dMap
}

////////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) ActiveDevice(device commons.DeviceMeta, isIPC bool) error {
	if err := bds.AddDevice(device, isIPC); err != nil {
		bds.logger.Errorf("active device(%+v) error: %s", device, err)
		return err
	}

	devActive := proto.DeviceActive{
		Id:    []string{device.Cid},
		IsAll: false,
	}
	bds.logger.Infof("driverItf active device, cid: %s", device.Cid)
	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancel()
	if activeRsp, err := bds.resourceCli.ActivateDevice(ctx, &devActive); err != nil {
		errmsg := status.Convert(err).Message()
		bds.logger.Errorf("ActivateDevice activeRsp, err:%s", activeRsp, errmsg)
		return errors.New(errmsg)
	}
	return nil
}

func (bds *BaseService) AddDevice(device commons.DeviceMeta, isIPC bool) error {
	if len(device.Cid) == 0 || len(device.ProductId) == 0 {
		return errors.New("required device id and productid")
	}
	bds.logger.Debugf("AddDevice info: %+v", device)

	// 查询此cid是否填加过
	if _, ok := bds.devCache.ById(device.Cid); ok {
		bds.logger.Infof("addDevice cid:%s is exist, can't add again", device.Cid)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancel()
	dev, err := transform.FromDeviceActiveInfoModelToProto(device, bds.GetServiceId(), isIPC)
	if err != nil {
		return err
	}

	devReq := transform.NewAddDeviceRequest(dev)
	devReq.Device.Source = common.DeviceFromDriver
	_, err = bds.resourceCli.AddDevice(ctx, devReq)
	if err != nil {
		bds.logger.Errorf("active device(%+v) error: %s", device, err)
		return errors.New(status.Convert(err).Message())
	}
	// add to resource success so add to local cache now
	bds.devCache.Add(commons.DeviceInfo{
		DeviceMeta:   device,
		ActiveStatus: commons.DeviceActiveStatusInactivated,
		OnLineStatus: string(commons.Offline),
	})
	return nil
}

func (bds *BaseService) SetDeviceProperty(cid string, property commons.ExtendedProperty) error {
	// 查询此cid是否存在
	devInfo, ok := bds.devCache.ById(cid)
	if !ok {
		return fmt.Errorf("setDeviceProperty error, device(%s) not exists", cid)
	}

	extendDataByte, err := json.Marshal(property.ExtendData)
	if err != nil {
		return err
	}
	extendDataStr := string(extendDataByte)

	devUpdate := proto.UpdateDeviceRequest{
		UpdateDevice: &proto.DeviceUpdateInfo{
			Id: cid,
			//VendorCode:      &property.VendorCode,
			InstallLocation: &property.InstallLocation,
			ExtendData:      &extendDataStr,
			Source:          common.DeviceFromDriver, // 添加设备更新请求来源
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancel()
	if updateRsp, err := bds.resourceCli.UpdateDevice(ctx, &devUpdate); err != nil {
		bds.logger.Errorf("UpdateDevice updateRsp:%+v, err:%v", updateRsp, err)
		return errors.New(status.Convert(err).Message())
	}

	devInfo.ExtendedAttr = property
	bds.devCache.Update(devInfo)

	return nil
}

func (bds *BaseService) SetDeviceBaseAttr(cid string, baseAttr commons.BaseProperty) error {
	// 查询此cid是否存在
	devInfo, ok := bds.devCache.ById(cid)
	if !ok {
		return fmt.Errorf("setDeviceBaseAttr error, device(%s) not exists", cid)
	}
	devUpdate := proto.UpdateDeviceRequest{
		UpdateDevice: &proto.DeviceUpdateInfo{
			Id:     cid,
			Name:   &baseAttr.Name,
			Ip:     &baseAttr.Ip,
			Lat:    &baseAttr.Lat,
			Lon:    &baseAttr.Lon,
			Source: common.DeviceFromDriver, // 添加设备更新请求来源
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancel()
	if updateRsp, err := bds.resourceCli.UpdateDevice(ctx, &devUpdate); err != nil {
		bds.logger.Errorf("UpdateDevice updateRsp:%+v, err:%v", updateRsp, err)
		return errors.New(status.Convert(err).Message())
	}

	devInfo.BaseAttr = baseAttr
	bds.devCache.Update(devInfo)

	return nil
}

func (bds *BaseService) DeleteDevice(cid string) error {
	if len(cid) == 0 {
		return errors.New("required device cid")
	}

	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancel()
	deleteRequest := proto.DeleteDeviceByIdRequest{
		Id: cid,
	}

	if _, err := bds.resourceCli.DeleteDeviceById(ctx, &deleteRequest); err != nil {
		return errors.New(status.Convert(err).Message())
	}

	bds.devCache.RemoveById(cid)
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) ReportDeviceStatus(data *commons.DeviceStatus) error {
	in, err := transform.DeviceStatusToProto(data)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancel()
	if _, err = bds.resourceCli.ReportDevicesOnlineAndOffline(ctx, in); err != nil {
		return errors.New(status.Convert(err).Message())
	}

	for _, online := range data.Online {
		info, exist := bds.devCache.ById(online)
		if exist {
			info.OnLineStatus = string(commons.Online)
			bds.devCache.Update(info)
		}
	}

	for _, offline := range data.Offline {
		info, exist := bds.devCache.ById(offline)
		if exist {
			info.OnLineStatus = string(commons.Offline)
			bds.devCache.Update(info)
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) ReportThroughHttp(api, version string, payload map[string]interface{}) (string, error) {
	in, err := transform.HTTPDataReportToProto(api, version, payload)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCHttpTimeout)
	defer cancel()
	if bds.resourceCli == nil {
		//bds.logger.Error("resourceCli is nil")
		return "", errors.New("report client is nil")
	}
	resp, err := bds.resourceCli.HttpReportData(ctx, in)
	if err != nil {
		//bds.logger.Errorf("reportThroughHttp api:%s, version:%s err:%s", api, version, err)
		return "", errors.New(status.Convert(err).Message())
	}
	return resp.GetMessage(), nil
}

func (bds *BaseService) CmdResultUpload(sn int64, success int, message string) error {
	requestMap := map[string]interface{}{
		"sn":      sn,
		"success": success,
		"message": message,
	}
	resp, err := bds.ReportThroughHttp(common.HTTP_API_EG_SYNCDATA_RESULT, "1.0", requestMap)
	if err != nil {
		bds.logger.Errorf("cmdResultUpload error, api:%s, requestMap:%+v, resp:%+v, err:%v",
			common.HTTP_API_EG_SYNCDATA_RESULT, requestMap, resp, err)
		return err
	}
	return nil
}

func (bds *BaseService) ReportAlert(ctx context.Context, level commons.AlertLevel, content string) error {
	if len(content) > common.ContentMaxLen {
		return fmt.Errorf("alert.Content length is too long, max length is %d bytes", common.ContentMaxLen)
	}

	if ctx == nil {
		c, cf := context.WithTimeout(context.Background(), common.GRPCTimeout)
		ctx = c
		defer cf()
	}
	if level != commons.ERROR && level != commons.WARN && level != commons.NOTIFY {
		bds.logger.Errorf("report alert level error, %s", level)
		return errors.New("report alert type or level error")
	}

	if resp, err := bds.resourceCli.AlertReport(ctx, transform.ToAlertReportProto(bds.GetServiceId(), time.Now().Unix(), level, content)); err != nil {
		return errors.New(status.Convert(err).Message())
	} else {
		if resp.GetStatusCode() != 0 {
			return errors.New(resp.GetMessage())
		}
		return nil
	}
}

func (bds *BaseService) ReportDevEvent(deviceId, event_type, deviceAddr, content string) error {
	ctx, cancelF := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancelF()

	if event_type != commons.DEVICE_REPORT_EVENT && event_type != commons.DEVICE_ALERT_EVENT {
		bds.logger.Errorf("report device deviceId:%s, event_type:%s", deviceId, event_type)
		return errors.New("device event_type error")
	}

	req := &proto.EventRequest{
		DeviceId:   deviceId,
		T:          time.Now().Unix(),
		EventType:  event_type,
		Message:    content,
		DeviceAddr: deviceAddr,
	}

	resp, err := bds.resourceCli.EventReport(ctx, req)
	if err != nil {
		return errors.New(status.Convert(err).Message())
	}

	if resp.GetStatusCode() != 0 {
		return errors.New(resp.GetMessage())
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) UploadFile(content []byte, fileName, subjectType string, timeout int32) (string, error) {
	token, err := bds.getUploadToken(subjectType)
	if err != nil {
		return "", err
	}

	if timeout == 0 {
		timeout = 5
	}
	result := token.Result

	// check env
	if strings.ToLower(bds.gwInfo.Env) == "pre" {
		if u, err := url.Parse(result.UploadUrl); err != nil {
			return "", err
		} else {
			if u.Port() != "7799" {
				result.UploadUrl = fmt.Sprintf("%s://%s:7799%s", u.Scheme, u.Host, u.Path)
			}
		}
	}

	response, err := httpUploadFile(content, result.UploadUrl, result.Token, fileName, timeout, bds.logger)
	bds.logger.Infof("UploadFile Url:%s, err:%v", result.UploadUrl, err)
	if err != nil {
		bds.logger.Errorf("DoAlarm UploadImage err:%v", err)
		return "", fmt.Errorf("upload fail:%s", err)
	}

	imageID := response.Result.TmpFileId
	bds.logger.Infof("UploadFile success, imageID:%s, resp:%+v", imageID, response)
	return imageID, nil
}

func (bds *BaseService) UploadFileV2(cid, fileName string, content []byte, timeout int32) (string, string, error) {
	var (
		err    error
		ok     bool
		cfg    *common.ResultNode
		client clients.OSSClientI
	)
	if cfg, ok = bds.ossManager.getSubDevStorConfig(cid); !ok {
		if cfg, err = bds.getSubDevStorConfigFromCloud(cid); err != nil {
			bds.logger.Errorf("uploadFileV2 err:%s", err)
			return "", "", err
		}
	} else {
		now := time.Now().UTC()
		exp, err := time.Parse("2006-01-02T15:04:05Z", cfg.Expiration)
		if err != nil {
			if cfg, err = bds.getSubDevStorConfigFromCloud(cid); err != nil {
				bds.logger.Errorf("uploadFileV2 err:%s", err)
				return "", "", err
			}
		} else {
			if exp.Sub(now) < 5 { // ugly
				bds.logger.Warnf("sub dev stor config expired, %s, update again", cfg.Expiration)
				if cfg, err = bds.getSubDevStorConfigFromCloud(cid); err != nil {
					bds.logger.Errorf("uploadFileV2 err:%s", err)
					return "", "", err
				}
			}
		}
	}

	if client, err = clients.NewOSSClient(cfg, bds.logger); err != nil {
		bds.logger.Errorf("new oss client with config(%s) error: %s", cfg, err)
		return "", "", err
	}
	if timeout == 0 {
		timeout = 5
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	return client.UploadFile(ctx, fileName, content)

}

/////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) SendToApp(ctx context.Context, appName string, req commons.AppDriverReq, cnnNum ...int) (commons.Response, error) {
	if cnnNum == nil {
		cnnNum = []int{0}
	}

	appRpcClient, err := bds.appCliManager.getAppServiceRpcClient(appName, cnnNum...)
	if err != nil {
		bds.logger.Errorf("SendToApp get appService Client error: %s", err.Error())
		return commons.Response{}, fmt.Errorf("app service(%s) rpc client is not created", appName)
	}

	var resp *proto.SendResponse
	rpcReq := req.ToRpc()
	rpcReq.Name = appName
	resp, err = appRpcClient.SendToAppService(ctx, rpcReq)
	if err != nil {
		bds.logger.Errorf("SendToApp send data failed: %v", err)
		return commons.Response{}, err
	}

	return commons.Response{
		Success: resp.Success,
		Message: resp.Message,
		Payload: resp.Payload,
	}, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) DriverProxyRegist(host string, port int) error {
	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancel()

	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("error host:%s", host)
	}

	if port <= 0 {
		return fmt.Errorf("error port:%d", port)
	}

	dPxyReq := &proto.DriverProxyRequest{
		Id:   bds.GetServiceId(),
		Name: bds.GetServiceName(),
		Host: host,
		Port: strconv.Itoa(port),
	}

	_, err := bds.resourceCli.DriverProxyRegister(ctx, dPxyReq)
	return err
}

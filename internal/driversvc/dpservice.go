package driversvc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/cache"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/interfaces"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/retrans"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/server"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/transform"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/utils"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/worker"
	"google.golang.org/grpc/status"
)

type DPModelService struct {
	*BaseService

	driverItf  dpmodel.DPModelDriver
	rpcServer  interfaces.RPCServerItf
	pdCache    *cache.DPModelProductCache
	reportPool *worker.DPWorkerPool
	RtsManager *retrans.ReTransfer //DP重传manager
}

func NewDPModelService(bds *BaseService) *DPModelService {
	lc := bds.GetLogger()
	var err error

	if bds.GetTEdgeModel() != commons.DPModel {
		lc.Errorf("NewDPModelService but TEdge run in mod:%s, exit", bds.GetTEdgeModel())
		os.Exit(-1)
	}

	var dds = &DPModelService{
		BaseService: bds,
		pdCache:     cache.NewDPProductCache(),
	}

	// sync product
	if dds.libraryId, err = dds.initDPProductCache(); err != nil {
		lc.Errorf("NewDPModelService init dpmodel product error:%s", err)
		os.Exit(-1)
	}

	// new pool
	clientInfo := bds.configMgr.GetClient()
	if dds.reportPool, err = worker.NewDPWorkerPool(clientInfo[common.Resource], lc); err != nil {
		lc.Errorf("NewDPModelService new worker pool error: %s", err)
		os.Exit(-1)
	}

	return dds
}

func (dds *DPModelService) Start(driver dpmodel.DPModelDriver, opts ...commons.Option) error {
	lc := dds.GetLogger()
	if driver == nil {
		lc.Errorf("Start: dpmodel.DPModelDriver unimplemented")
		os.Exit(-1)
	}
	dds.driverItf = driver

	for _, opt := range opts {
		opt.Apply(&dds.options)
	}
	dds.appCliManager.AppDataHandler = dds.options.AppDataHandler

	if err := dds.initMqttClient(); err != nil {
		lc.Infof("Start fail, initMqttClient err:%s", err)
		return err
	}

	go dds.waitSignalsExit()
	go dds.SyncDeviceAndProduct() //定时从TEdge更新设备和产品信息
	dds.initRtsManager() //初始化rtsManager
	close(dds.readyChan)

	dds.rpcServer = server.NewDPDriverRpcServer(dds)

	err := dds.rpcServer.Serve()
	dds.wg.Wait()

	lc.Warnf("DPModelService driver exit...")
	return err
}

func (dds *DPModelService) waitSignalsExit() {
	stopSignalCh := make(chan os.Signal, 1)
	signal.Notify(stopSignalCh, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	dds.wg.Add(1)

	for {
		select {
		case <-dds.ctx.Done():
			dds.logger.Infof("inner cancel executed, exit...")
			dds.reportPool.Stop()

			if dds.RtsManager != nil {
				dds.RtsManager.Exit()
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			if err := dds.driverItf.Stop(ctx); err != nil {
				dds.logger.Errorf("call protocol driverItf stop function error: %s", err)
			}
			cancel()
			dds.wg.Done()
			return
		case sig := <-stopSignalCh:
			dds.logger.Infof("DPModelService got signal:%v, exit...", sig)
			dds.cancel()
		}
	}
}

func (dds *DPModelService) initDPProductCache() (string, error) {
	lc := dds.GetLogger()
	serviceId := dds.GetServiceId()
	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancel()
	ds, err := dds.resourceCli.DeviceServiceById(ctx, &proto.DeviceServiceByIdRequest{Id: serviceId})
	if err != nil {
		lc.Errorf("initDPProductCache get deviceService Id:%s error:%s", serviceId, err)
		return "", err
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancel2()
	multiProducts, err := dds.resourceCli.RpcProductClient.ProductsSearch(ctx2, &proto.ProductSearchQueryRequest{
		BaseSearchConditionQuery: &proto.BaseSearchConditionQuery{},
		DeviceLibraryId:          ds.DeviceLibraryId,
	})
	if err != nil {
		lc.Errorf("initDPProductCache ProductsSearch libraryId:%s error:%s", ds.DeviceLibraryId, err)
		return ds.DeviceLibraryId, err
	}

	// sync product list
	for _, p := range multiProducts.Products {
		dpModel := transform.ToDPProductModel(p)
		if _, ok := dds.pdCache.ById(p.Id); ok {
			dds.pdCache.Update(dpModel)
		} else {
			dds.pdCache.Add(dpModel)
		}
		dpModelStr, _ := json.Marshal(dpModel)
		lc.Debugf("initDPProductCache pid:%s, dpModelStr:%s", p.Id, dpModelStr)
	}

	return ds.DeviceLibraryId, nil
}

func (dds *DPModelService) SyncDeviceAndProduct() {
	timerA := time.NewTicker(time.Minute * 3)
	for {
		select {
		case <-dds.ctx.Done():
			timerA.Stop()
			return
		case <-timerA.C:
			dds.initDeviceCache()
			dds.initDPProductCache()
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func (dds *DPModelService) GetPdCache() *cache.DPModelProductCache {
	return dds.pdCache
}

func (dds *DPModelService) GetDriver() dpmodel.DPModelDriver {
	return dds.driverItf
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
func (dds *DPModelService) ReportWithDPData(cid string, data []*dpmodel.WithDPValue) error {
	dev, ok := dds.devCache.ById(cid)
	if !ok {
		return fmt.Errorf("no cid:%s in local cache", cid)
	}

	d, err := transform.NewWithDPData(cid, data, dev, dds.pdCache)
	if err != nil {
		dds.logger.Errorf("new with dp data error: %s", err)
		return err
	}
	return dds.reportPool.PutData(&proto.Events{MultiEvents: []*proto.Event{d.ToProto()}})
}

func (dds *DPModelService) ReportWithoutDPData(topic string, data *dpmodel.WithoutDPValue) error {
	d, err := transform.WithoutDPDataToProto(&transform.WithoutDPReport{
		Topic: topic,
		Data:  data,
	})
	if err != nil {
		dds.logger.Errorf("without dp data transform error: %s", err)
		return err
	}
	return dds.reportPool.PutData(d)
}

func (dds *DPModelService) AddProduct(p dpmodel.DPModelProductAddInfo) error {
	if len(p.Id) == 0 {
		return errors.New("required productId")
	}
	_, ok := dds.pdCache.ById(p.Id)
	if ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancel()

	product, err := transform.FromAddProductModel2Proto(&p, dds.libraryId)
	if err != nil {
		return err
	}
	//添加产品
	if _, err = dds.resourceCli.ProductAdd(ctx, product); err != nil {
		errmsg := status.Convert(err).Message()
		dds.logger.Errorf("add product(%+v) error: %s", product, errmsg)
		return errors.New(status.Convert(err).Message())
	}
	dds.pdCache.Add(dpmodel.DPModelProduct{
		Id:          p.Id,
		Name:        p.Name,
		Model:       p.Model,
		Description: p.Description,
	})
	return nil
}

func (dds *DPModelService) AllProducts() map[string]dpmodel.DPModelProduct {
	return dds.pdCache.All()
}

// GetProductById 根据产品ID查询产品信息
func (dds *DPModelService) GetProductById(pid string) (dpmodel.DPModelProduct, bool) {
	return dds.pdCache.ById(pid)
}

//激活子设备后，单独更新IPC能力集时使用
func (dds *DPModelService) IpcDeviceUpdateSkill(version, cid string) (string, error) {
	ipcSkillSetStr, err := transform.GenIpcSkillSet()
	if err != nil {
		return "", err
	}

	reqBodyMap := make(map[string]interface{})
	reqBodyMap["subId"] = cid
	reqBodyMap["skill"] = ipcSkillSetStr

	resp, err := dds.ReportThroughHttp(common.HTTP_DEVICE_SKILL_UPDATE, version, reqBodyMap)
	if err != nil {
		dds.logger.Errorf("ipcDeviceUpdateSkill failed, cid:%s, resp:%s, err:%v", cid, resp, err)
		return resp, err
	}

	ipcSkillStr, _ := json.Marshal(reqBodyMap)
	dds.logger.Debugf("ipcDeviceUpdateSkill cid:%s, ipcSkillStr:%s", cid, ipcSkillStr)

	return resp, err
}

////////////////////////////////////////////////////////////////////////////////////////////////
func (dds *DPModelService) initRtsManager() {
	if !dds.options.EnableRTS {
		return
	}

	rtsManager, err := retrans.NewReTransfer(dds.logger, *retrans.DefaultRtOption())
	if err != nil {
		dds.logger.Errorf("initRtsManager failed, cant enable retrans err:%v", err)
		return
	}

	//初始化网关状态, 初始化 RtsManager
	dds.RtsManager = rtsManager
	dds.RtsManager.SetConnItf(dds)
	dds.RtsManager.SetMessageItf(dds)
	go dds.RtsManager.Running()

	dds.logger.Infof("initRtsManager done, cloudStatus:%v", dds.GetGatewayMeta().CloudState)
}

//DP断网续传接口
func (dds *DPModelService) DPRePublish(cid string, data []byte) error {
	lc := dds.logger

	//unserialize
	var values []*dpmodel.WithDPValue
	err := utils.JsonDecoder(data, &values)
	if err != nil {
		lc.Errorf("DPRePublish decode err:%v, data:%s", err, data)
		return err //1.不应该出现；若出现则是代码Bug
	}

	//类型强制转换
	for _, value := range values {
		dpId := value.DPId
		dpType := value.DPType
		if value.DPType == commons.RawType {
			rawData, _ := base64.StdEncoding.DecodeString(value.Value.(string))
			value.Value = rawData
		} else if dpType == commons.ValueType {
			numValue, ok := value.Value.(json.Number)
			if ok {
				intValue, _ := numValue.Int64()
				value.Value = intValue
			} else {
				lc.Errorf("2.DPRePublish dpId:%s, dpType:%s value:%v, valueType:%v not int64", dpId, dpType, value.Value, reflect.TypeOf(value.Value))
			}
		} else if dpType == commons.FaultType || dpType == commons.BitmapType {
			numValue, ok := value.Value.(json.Number)
			if ok {
				intValue, _ := numValue.Int64()
				value.Value = int32(intValue)
			} // else []string
		}
	}

	//fmt.Printf(">>>DPRePublish cid:%s, data type:%v value:%+v\n", cid, reflect.TypeOf(data), values[0])
	err = dds.ReportWithDPData(cid, values)
	lc.Infof("options DPRePublish cid:%s, data:%s, err:%v", cid, data, err)
	return err //2.不应该出现；若出现则是代码Bug
}

/////////////////////////////////////////////////////////////////////////////////////////////////
//Atop 失败重传
//key: api|version
func (dds *DPModelService) AtopReReport(key string, data []byte) error {
	lc := dds.logger
	//split atop key
	skey := strings.Split(key, retrans.AtopKeySep)
	if len(skey) != 2 {
		lc.Errorf("1.AtopReReport key:%s format err, cant happened!", key)
		return nil //1.不可能出现，所以不能返回错误，否则重传管理器无法正常工作; 若出现则是代码Bug
	}

	//unserialize payload
	var payload map[string]interface{}
	err := utils.JsonDecoder(data, &payload)
	if err != nil {
		lc.Errorf("2.AtopReReport key:%s payload unmarshal err:%v, cant happened!", key, err)
		return nil //2.不可能出现，所以不能返回错误，否则重传管理器无法正常工作; 若出现则是代码Bug
	}

	httpApi := skey[0]
	httpVersion := skey[1]
	resp, err := dds.ReportThroughHttp(httpApi, httpVersion, payload)
	if err == nil {
		return checkHttpResp([]byte(resp))
	}

	return err
}

func (dds *DPModelService) ReportWithDPDataV2(cid string, data []*dpmodel.WithDPValue) error {
	lc := dds.logger
	if !dds.options.EnableRTS {
		lc.Errorf("reportWithDPDataV2 cid:%s, rts manager not enable!!!", cid)
		return fmt.Errorf("rts manager not enable")
	}

	if dds.GetCloudStatus() {
		return dds.ReportWithDPData(cid, data)
	}

	value, err := json.Marshal(data)
	if err != nil {
		lc.Errorf("reportWithDPDataV2 cid:%s, marshal error: %s", cid, err)
		return err
	}

	if dds.RtsManager != nil {
		return dds.RtsManager.SaveDPKV([]byte(cid), value)
	}

	lc.Errorf("reportWithDPDataV2 cid:%s failed, rts manager not inited!!!", cid)
	return fmt.Errorf("RtsManager not init, check driverItf log")
}

//Atop 重传策略：
//1. timeout; http 网关返回非200?
//2. atop 网关出错，通过 errorCode 进行过滤?
func (dds *DPModelService) ReportHttpWithRetrans(api, version string, payload map[string]interface{}) (string, error) {
	in, err := transform.HTTPDataReportToProto(api, version, payload)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), common.GRPCHttpTimeout)
	defer cancel()
	if dds.resourceCli == nil {
		//dds.logger.Error("1.reportHttpWithRetrans resourceCli is nil")
		return "", errors.New("report client is nil")
	}

	resp, err := dds.resourceCli.HttpReportData(ctx, in)
	if err != nil {
		dds.logger.Errorf("2.reportHttpWithRetrans save atop data, report err:%s", err)
		if dds.RtsManager != nil {
			dds.RtsManager.SaveAtopKV(api, version, in.Payload)
		}
		return "", errors.New(status.Convert(err).Message())
	}

	if err := checkHttpResp([]byte(resp.Message)); err != nil {
		dds.logger.Errorf("3.reportHttpWithRetrans save atop data, checkHttpResp err:%s", err)
		if dds.RtsManager != nil {
			dds.RtsManager.SaveAtopKV(api, version, in.Payload)
		}
	}

	return resp.GetMessage(), nil
}

package driversvc

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/cache"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/interfaces"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/server"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/snowflake"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/transform"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/worker"
	"github.com/tuya/tuya-tedge-driver-sdk-go/thingmodel"
	"google.golang.org/grpc/status"
)

type TyModelService struct {
	*BaseService

	node       *snowflake.Node
	driver     thingmodel.ThingModelDriver
	reportPool *worker.TMWorkerPool
	pdCache    *cache.ThingModelProductCache
	rpcServer  interfaces.RPCServerItf
}

func NewTyModelService(bds *BaseService) *TyModelService {
	lc := bds.GetLogger()
	var err error

	if bds.GetTEdgeModel() != commons.ThingModel {
		lc.Errorf("NewTyModelService but TEdge run in mod:%s, exit", bds.GetTEdgeModel())
		os.Exit(-1)
	}

	var tmds = &TyModelService{
		BaseService: bds,
		pdCache:     cache.NewTyModelProduct(),
	}

	// sync product
	if tmds.libraryId, err = tmds.initTyProductCache(); err != nil {
		lc.Errorf("NewTyModelService sync thingmodel product error: %s", err)
		os.Exit(-1)
	}

	// new pool
	clientInfo := bds.configMgr.GetClient()
	if tmds.reportPool, err = worker.NewTMWorkerPool(clientInfo[common.Resource], lc); err != nil {
		lc.Errorf("NewTyModelService new worker pool error: %s", err)
		os.Exit(-1)
	}

	tmds.node, err = snowflake.NewNode(1)
	if err != nil {
		lc.Errorf("NewTyModelService new msg id generator error: %s", err)
		os.Exit(-1)
	}

	return tmds
}

func (tmds *TyModelService) initTyProductCache() (string, error) {
	lc := tmds.GetLogger()
	serviceId := tmds.GetServiceId()

	ctx1, cancelF := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancelF()
	ds, err := tmds.resourceCli.DeviceServiceById(ctx1, &proto.DeviceServiceByIdRequest{Id: serviceId})
	if err != nil {
		lc.Errorf("initTyProductCache deviceService Id:%s error:%s", serviceId, err)
		return "", err
	}

	ctx2, cancelF := context.WithTimeout(context.Background(), common.GRPCTimeout)
	defer cancelF()
	multiProducts, err := tmds.resourceCli.RpcThingModelClient.ProductsSearch(ctx2, &proto.TMProductSearchQueryRequest{
		BaseSearchConditionQuery: &proto.BaseSearchConditionQuery{},
		DeviceLibraryId:          ds.DeviceLibraryId,
	})
	if err != nil {
		lc.Errorf("initTyProductCache ProductsSearch libraryId:%s error:%s", ds.DeviceLibraryId, err)
		return ds.DeviceLibraryId, err
	}

	for _, p := range multiProducts.Products {
		tyModel := transform.ToTMProductModel(p)
		if _, ok := tmds.pdCache.ById(p.Id); ok {
			tmds.pdCache.Update(tyModel)
		} else {
			tmds.pdCache.Add(tyModel)
		}
		tyPStr, _ := json.Marshal(tyModel)
		lc.Debugf("initTyProductCache tyPStr:%s", tyPStr)
	}

	return ds.DeviceLibraryId, nil
}

func (tmds *TyModelService) Start(driver thingmodel.ThingModelDriver, opts ...commons.Option) error {
	lc := tmds.GetLogger()
	if driver == nil {
		lc.Errorf("thingsmodel.ThingModelDriver unimplemented")
		return errors.New("thingsmodel.ThingModelDriver unimplemented")
	}
	tmds.driver = driver

	for _, opt := range opts {
		opt.Apply(&tmds.options)
	}

	// with mqtt
	if err := tmds.initMqttClient(); err != nil {
		lc.Infof("Start fail, initMqttClient err:%s", err)
		return err
	}
	tmds.appCliManager.AppDataHandler = tmds.options.AppDataHandler

	go tmds.waitSignalsExit()
	go tmds.SyncDeviceAndProduct()

	close(tmds.readyChan)
	tmds.rpcServer = server.NewTyDriverRpcServer(tmds)

	err := tmds.rpcServer.Serve()
	tmds.wg.Wait()
	lc.Warnf("TyModelService driver exit...")

	return err
}

func (tmds *TyModelService) waitSignalsExit() {
	stopSignalCh := make(chan os.Signal, 1)
	signal.Notify(stopSignalCh, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	tmds.wg.Add(1)

	for {
		select {
		case <-tmds.ctx.Done():
			tmds.logger.Infof("inner cancel executed, exit...")
			tmds.reportPool.Stop()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			if err := tmds.driver.Stop(ctx); err != nil {
				tmds.logger.Errorf("call protocol driverItf stop function error: %s", err)
			}
			cancel()
			tmds.wg.Done()
			return
		case sig := <-stopSignalCh:
			tmds.logger.Infof("TyModelService got signal:%v, exit...", sig)
			tmds.cancel()
		}
	}
}

func (tmds *TyModelService) SyncDeviceAndProduct() {
	timerA := time.NewTicker(time.Minute * 3)
	for {
		select {
		case <-tmds.ctx.Done():
			timerA.Stop()
			return
		case <-timerA.C:
			tmds.initDeviceCache()
			tmds.initTyProductCache()
		}
	}
}

func (tmds *TyModelService) AddProduct(pr thingmodel.AddProductReq) error {
	if len(pr.Id) == 0 || len(pr.Name) == 0 {
		return errors.New("required productId")
	}
	_, ok := tmds.pdCache.ById(pr.Id)
	if ok {
		return nil
	}

	product, err := transform.FromAddTyProductModel2Proto(&pr, tmds.libraryId)
	if err != nil {
		return err
	}

	//添加产品
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if _, err = tmds.resourceCli.ProductAdd(ctx, product); err != nil {
		tmds.logger.Errorf("add product(%+v) error: %s", product, err)
		return errors.New(status.Convert(err).Message())
	}
	return nil
}

func (tmds *TyModelService) GetPdCache() *cache.ThingModelProductCache {
	return tmds.pdCache
}

func (tmds *TyModelService) GetDriver() thingmodel.ThingModelDriver {
	return tmds.driver
}

func (tmds *TyModelService) GetTMWorkerPool() *worker.TMWorkerPool {
	return tmds.reportPool
}

// AllProducts 获取当前实例下的所有产品
func (tmds *TyModelService) AllProducts() map[string]thingmodel.ThingModelProduct {
	return tmds.pdCache.All()
}

// GetProductById 根据产品ID获取产品信息
func (tmds *TyModelService) GetProductById(pid string) (thingmodel.ThingModelProduct, bool) {
	return tmds.pdCache.ById(pid)
}

func (tmds *TyModelService) GetEventByPid(pid string) (map[string]thingmodel.Event, bool) {
	return tmds.pdCache.GetEventsByPid(pid)
}

/////////////////////////////////////////////////////////////////////////////////////////////
func (tmds *TyModelService) AddLinkDevice(device commons.TMDeviceMeta) error {
	dev := transform.FromTMDeviceToDevice(device)
	return tmds.AddDevice(dev, false)
}

func (tmds *TyModelService) ActiveLinkDevice(device commons.TMDeviceMeta) error {
	dev := transform.FromTMDeviceToDevice(device)
	return tmds.ActiveDevice(dev, false)
}

func (tmds *TyModelService) AllTMDevices() map[string]commons.TMDeviceInfo {
	var devices = make(map[string]commons.TMDeviceInfo)
	for k, v := range tmds.devCache.All() {
		devices[k] = commons.TMDeviceInfo{
			TMDeviceMeta: commons.TMDeviceMeta{
				Cid:          v.Cid,
				ProductId:    v.ProductId,
				BaseAttr:     v.BaseAttr,
				ExtendedAttr: v.ExtendedAttr,
				Protocols:    v.Protocols,
			},
			DeviceId:     v.CloudDeviceId,
			ActiveStatus: v.ActiveStatus,
			OnLineStatus: v.OnLineStatus,
		}
	}
	return devices
}

func (tmds *TyModelService) GetLinkActiveDevices() map[string]commons.TMDeviceInfo {
	var devices = make(map[string]commons.TMDeviceInfo)
	for k, v := range tmds.GetActiveDevices() {
		devices[k] = commons.TMDeviceInfo{
			TMDeviceMeta: commons.TMDeviceMeta{
				Cid:          v.Cid,
				ProductId:    v.ProductId,
				BaseAttr:     v.BaseAttr,
				ExtendedAttr: v.ExtendedAttr,
				Protocols:    v.Protocols,
			},
			DeviceId:     v.CloudDeviceId,
			ActiveStatus: v.ActiveStatus,
			OnLineStatus: v.OnLineStatus,
		}
	}
	return devices
}

func (tmds *TyModelService) GetTMDeviceById(cid string) (commons.TMDeviceInfo, bool) {
	d, ok := tmds.devCache.ById(cid)
	if !ok {
		return commons.TMDeviceInfo{}, false
	}
	return commons.TMDeviceInfo{
		TMDeviceMeta: commons.TMDeviceMeta{
			Cid:          d.Cid,
			ProductId:    d.ProductId,
			BaseAttr:     d.BaseAttr,
			ExtendedAttr: d.ExtendedAttr,
			Protocols:    d.Protocols,
		},
		DeviceId:     d.CloudDeviceId,
		ActiveStatus: d.ActiveStatus,
		OnLineStatus: d.OnLineStatus,
	}, true
}

/////////////////////////////////////////////////////////////////////////////////////////////
func (tmds *TyModelService) PropertyReport(cid string, data thingmodel.PropertyReport) (thingmodel.CommonResponse, error) {
	msgId := tmds.GenRandomId()
	data.MsgId = msgId
	if data.Time == 0 {
		data.Time = time.Now().Unix()
	}
	msg, err := transform.TMDataToProto(cid, transform.PropertyReport, data)
	if err != nil {
		return thingmodel.CommonResponse{}, err
	}
	tmds.logger.Debugf(">>>>msg:%+v", msg)
	return tmds.withCommonResponse(msgId, data.Sys.Ack, msg)
}

func (tmds *TyModelService) EventReport(cid string, data thingmodel.EventReport) (thingmodel.CommonResponse, error) {
	msgId := tmds.GenRandomId()
	data.MsgId = msgId
	msg, err := transform.TMDataToProto(cid, transform.EventReport, data)
	if err != nil {
		return thingmodel.CommonResponse{}, err
	}

	return tmds.withCommonResponse(msgId, data.Sys.Ack, msg)
}

func (tmds *TyModelService) BatchReport(cid string, data thingmodel.BatchReport) (thingmodel.CommonResponse, error) {
	msgId := tmds.GenRandomId()
	data.MsgId = msgId
	msg, err := transform.TMDataToProto(cid, transform.BatchReport, data)
	if err != nil {
		return thingmodel.CommonResponse{}, err
	}

	return tmds.withCommonResponse(msgId, data.Sys.Ack, msg)
}

func (tmds *TyModelService) PropertyDesiredGet(cid string, data thingmodel.PropertyDesiredGet) (thingmodel.PropertyDesiredGetResponse, error) {
	msgId := tmds.GenRandomId()
	data.MsgId = msgId
	msg, err := transform.TMDataToProto(cid, transform.PropertyDesiredGet, data)
	if err != nil {
		return thingmodel.PropertyDesiredGetResponse{}, err
	}
	resp, err := tmds.reportPool.PutDataWithMsgId(msgId, msg)
	if err != nil {
		return thingmodel.PropertyDesiredGetResponse{}, err
	}
	r, ok := resp.(thingmodel.PropertyDesiredGetResponse)
	if !ok {
		return thingmodel.PropertyDesiredGetResponse{}, errors.New("type error")
	}
	return r, nil
}

func (tmds *TyModelService) PropertyDesiredDelete(cid string, data thingmodel.PropertyDesiredDelete) (thingmodel.PropertyDesiredDeleteResponse, error) {
	msgId := tmds.GenRandomId()
	data.MsgId = msgId
	msg, err := transform.TMDataToProto(cid, transform.PropertyDesiredDelete, data)
	if err != nil {
		return thingmodel.PropertyDesiredDeleteResponse{}, err
	}
	resp, err := tmds.reportPool.PutDataWithMsgId(msgId, msg)
	if err != nil {
		return thingmodel.PropertyDesiredDeleteResponse{}, err
	}
	r, ok := resp.(thingmodel.PropertyDesiredDeleteResponse)
	if !ok {
		return thingmodel.PropertyDesiredDeleteResponse{}, errors.New("type error")
	}
	return r, nil
}

func (tmds *TyModelService) PropertySetResponse(cid string, data thingmodel.CommonResponse) error {
	msg, err := transform.TMDataToProto(cid, transform.PropertySetResponse, data)
	if err != nil {
		return err
	}
	return tmds.reportPool.PutData(msg)
}

func (tmds *TyModelService) PropertyGetResponse(cid string, data thingmodel.PropertyGetResponse) error {
	msg, err := transform.TMDataToProto(cid, transform.PropertyGetResponse, data)
	if err != nil {
		return err
	}
	return tmds.reportPool.PutData(msg)
}

func (tmds *TyModelService) ActionExecuteResponse(cid string, data thingmodel.ActionExecuteResponse) error {
	msg, err := transform.TMDataToProto(cid, transform.ActionExecuteResponse, data)
	if err != nil {
		return err
	}
	return tmds.reportPool.PutData(msg)
}

func (tmds *TyModelService) withCommonResponse(msgId string, ack int8, data *proto.ThingModelMsg) (thingmodel.CommonResponse, error) {
	if ack == 1 {
		resp, err := tmds.reportPool.PutDataWithMsgId(msgId, data)
		if err != nil {
			return thingmodel.CommonResponse{}, err
		}
		r, ok := resp.(thingmodel.CommonResponse)
		if !ok {
			return thingmodel.CommonResponse{}, errors.New("type error")
		}
		return r, nil
	} else {
		if err := tmds.reportPool.PutData(data); err != nil {
			return thingmodel.CommonResponse{}, err
		}
		return thingmodel.CommonResponse{}, nil
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////
func (tmds *TyModelService) HttpRequestProxy(params thingmodel.HttpRequestParam, payload map[string]interface{}, timeout int) (string, error) {
	in, err := transform.HttpWithCustomParamsToProto(params, payload)
	if err != nil {
		return "", err
	}

	timeoutDuration := time.Duration(timeout) * time.Second
	if timeout <= 0 {
		timeoutDuration = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()
	if tmds.resourceCli == nil {
		tmds.logger.Errorf("httpRequestProxy resourceCli is nil")
		return "", errors.New("report client is nil")
	}
	resp, err := tmds.resourceCli.HttpRequestProxy(ctx, in)
	if err != nil {
		tmds.logger.Errorf("httpRequestProxy report data error: %s", err)
		return "", errors.New(status.Convert(err).Message())
	}
	return resp.GetMessage(), nil
}

func (tmds *TyModelService) GenRandomId() string {
	randomId := tmds.node.Generate().String()
	return randomId
}

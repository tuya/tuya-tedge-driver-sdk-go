package driversvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/cache"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/clients"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/config"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
)

//////////////////////////////////////////////////////////////////////////////////////////////
type ConfigManger struct {
	configFile    string
	configuration config.DriverConfig
}

func (cm *ConfigManger) GetClient() map[string]config.ClientInfo {
	return cm.configuration.Clients
}

func (cm *ConfigManger) GetService() config.ServiceInfo {
	return cm.configuration.Service
}

// 获取驱动自定义配置
func (cm *ConfigManger) GetCustomConfig() map[string]interface{} {
	return cm.configuration.CustomConfig
}

func (cm *ConfigManger) SetLogLevel(level string) {
	cm.configuration.Logger.LogLevel = level
}

func (cm *ConfigManger) WriteToFile() error {
	return cm.configuration.WriteToFile(cm.configFile)
}

////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) GetLogger() commons.TedgeLogger {
	return bds.logger
}

func (bds *BaseService) GetRPCConfig() config.RPCConfig {
	return bds.configMgr.GetService().Server
}

func (bds *BaseService) GetDevCache() *cache.DeviceCache {
	return bds.devCache
}

func (bds *BaseService) GetContext() (context.Context, context.CancelFunc) {
	return bds.ctx, bds.cancel
}

func (bds *BaseService) GetAppHandler() commons.AppCallBack {
	return bds.appCliManager.AppDataHandler
}

func (bds *BaseService) UpdateAppAddress(ctx context.Context, req *proto.AppBaseAddress) error {
	return bds.appCliManager.updateAppAddress(ctx, req)
}

func (bds *BaseService) ChangeLogLevel(level string) error {
	lc := bds.GetLogger()
	lc.Infof("recv change log to level: %s", level)

	lc.SetLogLevel(level)

	// write to config file
	bds.configMgr.SetLogLevel(level)
	if err := bds.configMgr.WriteToFile(); err != nil {
		lc.Errorf("configuration write to file error: %s", err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (bds *BaseService) SetCloudStatus(status bool) {
	bds.cloudStatus.Store(status)
}

func (bds *BaseService) GetCloudStatus() bool {
	return bds.cloudStatus.Load()
}

//////////////////////////////////////////////////////////////////////////////////////////////
type ossManager struct {
	mu       sync.Mutex
	tokenMap map[string]*common.ResultNode
}

func newOssManager() *ossManager {
	return &ossManager{
		tokenMap: make(map[string]*common.ResultNode),
	}
}

func (ossm *ossManager) setSubDevStorConfig(cid string, cfg *common.ResultNode) {
	ossm.mu.Lock()
	defer ossm.mu.Unlock()
	ossm.tokenMap[cid] = cfg
}

func (ossm *ossManager) getSubDevStorConfig(cid string) (*common.ResultNode, bool) {
	ossm.mu.Lock()
	defer ossm.mu.Unlock()
	cfg, ok := ossm.tokenMap[cid]
	return cfg, ok
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
type appCliManager struct {
	logger      commons.TedgeLogger
	resourceCli *clients.ResourceClient

	mutex          sync.Mutex
	appClientMap   map[string]clients.AppClientCnnMap
	AppDataHandler commons.AppCallBack
}

func newAppCliManager(logger commons.TedgeLogger, resourceCli *clients.ResourceClient) *appCliManager {
	return &appCliManager{
		logger:       logger,
		resourceCli:  resourceCli,
		appClientMap: make(map[string]clients.AppClientCnnMap),
	}
}

func (appm *appCliManager) setAppHandler(handler commons.AppCallBack) {
	appm.AppDataHandler = handler
}

func (appm *appCliManager) getAppServiceRpcClient(appName string, cnnNum ...int) (*clients.AppClient, error) {
	appRpcCnnNum := cnnNum[0]
	appRpcClientCnnMap, ok := appm.appClientMap[appName]

	lc := appm.logger
	// 存在应用rpcClientMap
	if ok {
		if appRpcClient, exist := appRpcClientCnnMap[appRpcCnnNum]; exist {
			connStat := appRpcClient.Conn.GetState()
			if connStat == connectivity.Idle || connStat == connectivity.Ready {
				return appRpcClient, nil
			}

			lc.Warnf("app:%s conNum:%d conn status:%s", appName, appRpcCnnNum, connStat.String())
			appRpcClient.Close()
			appm.mutex.Lock()
			delete(appRpcClientCnnMap, appRpcCnnNum)
			appm.mutex.Unlock()
		}
	} else {
		appRpcClientCnnMap = make(clients.AppClientCnnMap)
	}

	appRpcClient, err := appm.getAppRpcCient(appName)
	if err != nil {
		return nil, err
	}

	appm.mutex.Lock()
	appRpcClientCnnMap[appRpcCnnNum] = appRpcClient
	appm.appClientMap[appName] = appRpcClientCnnMap
	appm.mutex.Unlock()
	return appm.appClientMap[appName][appRpcCnnNum], nil
}

func (appm *appCliManager) getAppRpcCient(appName string) (*clients.AppClient, error) {
	app, err := appm.resourceCli.GetAppRegisterName(context.Background(), &proto.AppByRegisterNameReq{
		Name: appName,
	})
	if err != nil {
		return nil, err
	}
	// 创建连接驱动的grpc client
	builder := &clients.AppResolverBuilder{
		AppSerivceName: app.Name,
		Addr:           app.BaseAddress,
		AddrChan:       make(chan string),
	}
	appRpcClient, err := clients.NewAppRpcClient(builder)
	if err != nil {
		return nil, err
	}
	appRpcClient.Builder = builder
	return appRpcClient, nil
}

func (appm *appCliManager) updateAppAddress(ctx context.Context, req *proto.AppBaseAddress) error {
	lc := appm.logger
	appName := req.Name
	lc.Debugf("UpdateAppAddress name:%s, addr:%s", appName, req.Addr)

	appm.mutex.Lock()
	defer appm.mutex.Unlock()

	appRpcClientMap, ok := appm.appClientMap[appName]
	if !ok {
		lc.Warnf("AppServiceAddress no appName:%s client map", appName)
		return nil
	}

	var wg sync.WaitGroup
	for _, client := range appRpcClientMap {
		wg.Add(1)
		go func(c *clients.AppClient) {
			ctx, cancel := context.WithTimeout(ctx, common.GRPCTimeout)
			defer func() {
				cancel()
				wg.Done()
			}()

			select {
			case <-ctx.Done():
				lc.Errorf("AppServiceAddress update app (%s) address failed", appName)
			case c.Builder.AddrChan <- req.Addr:
				lc.Infof("AppServiceAddress update app (%s) address success", appName)
			}
		}(client)
	}
	wg.Wait()

	return nil
}

//func (bds *BaseService) getAppServiceRpcClient(appName string, cnnNum ...int) (*clients.AppClient, error) {
//	appRpcCnnNum := cnnNum[0]
//	appRpcClientCnnMap, ok := bds.appClientMap[appName]
//
//	// 存在应用rpcClientMap
//	if ok {
//		if appRpcClient, exist := appRpcClientCnnMap[appRpcCnnNum]; exist {
//			connStat := appRpcClient.ConnMqtt.GetState()
//			if connStat == connectivity.Idle || connStat == connectivity.Ready {
//				return appRpcClient, nil
//			}
//
//			bds.logger.Warnf("app:%s conNum:%d conn status:%s", appName, appRpcCnnNum, connStat.String())
//			appRpcClient.Close()
//			bds.mutex.Lock()
//			delete(appRpcClientCnnMap, appRpcCnnNum)
//			bds.mutex.Unlock()
//		}
//	} else {
//		appRpcClientCnnMap = make(clients.AppClientCnnMap)
//	}
//
//	appRpcClient, err := getAppRpcCient(bds, appName)
//	if err != nil {
//		return nil, err
//	}
//
//	bds.mutex.Lock()
//	appRpcClientCnnMap[appRpcCnnNum] = appRpcClient
//	bds.appClientMap[appName] = appRpcClientCnnMap
//	bds.mutex.Unlock()
//	return bds.appClientMap[appName][appRpcCnnNum], nil
//}
//
//func getAppRpcCient(dds *BaseService, appName string) (*clients.AppClient, error) {
//	app, err := dds.resourceCli.GetAppRegisterName(context.Background(), &proto.AppByRegisterNameReq{
//		Name: appName,
//	})
//	if err != nil {
//		return nil, err
//	}
//	// 创建连接驱动的grpc client
//	builder := &clients.AppResolverBuilder{
//		AppSerivceName: app.Name,
//		Addr:           app.BaseAddress,
//		AddrChan:       make(chan string),
//	}
//	appRpcClient, err := clients.NewAppRpcClient(builder)
//	if err != nil {
//		return nil, err
//	}
//	appRpcClient.Builder = builder
//	return appRpcClient, nil
//}

//////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) getUploadToken(subjectType string) (*common.UploadTokenResp, error) {
	reqBodyMap := make(map[string]interface{})
	if subjectType == "" {
		reqBodyMap["subjectType"] = "edgegateway_alarmImg"
	} else {
		reqBodyMap["subjectType"] = subjectType
	}
	resp, err := bds.ReportThroughHttp(common.HTTP_API_EDGE_UPLOAD_TOKEN, "1.0", reqBodyMap)
	if err != nil {
		bds.logger.Errorf("GetUploadToken get token err:%s, response:%+v", err, resp)
		return nil, err
	}

	uploadResp := new(common.UploadTokenResp)
	if err = json.Unmarshal([]byte(resp), uploadResp); err != nil {
		bds.logger.Errorf("GetUploadToken Unmarshal token err:%s, response:%+v\n", err, resp)
		return nil, err
	}

	return uploadResp, nil
}

func httpUploadFile(data []byte, remoteURL, token, filename string, timeout int32, logger commons.TedgeLogger) (resp *common.UploadImageResp, err error) {
	dataReader := bytes.NewReader(data)
	values := map[string]io.Reader{
		"file": dataReader, // lets assume its this file
	}

	uploadRes, err := utils.HttpUpload(remoteURL, values, token, filename, timeout)
	if err != nil {
		logger.Errorf("UploadImage error:%s, token:%s, filename:%s", err, token, filename)
		return nil, err
	}

	logger.Debugf("uoload image resp: %s", string(uploadRes))
	uploadResp := new(common.UploadImageResp)
	err = json.Unmarshal(uploadRes, uploadResp)
	if err != nil {
		logger.Errorf("UploadImage unmarshal error:%s, token:%s, filename:%s", err, token, filename)
		return nil, err
	}

	return uploadResp, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////
func (bds *BaseService) getSubDevStorConfigFromCloud(cid string) (*common.ResultNode, error) {
	var (
		err     error
		resp    string
		cidList = make([]string, 0, 1)
	)
	cidList = append(cidList, cid)
	cidListStr, _ := json.Marshal(cidList)
	reqBodyMap := make(map[string]interface{})
	reqBodyMap["nodeIds"] = string(cidListStr)
	reqBodyMap["type"] = "Motion"

	if resp, err = bds.ReportThroughHttp(common.HTTP_DEVICE_TY_SUB_STOR_CONFIG_GET, "2.0", reqBodyMap); err != nil {
		bds.logger.Errorf("get sub device(%s) stor config error: %s", cid, err)
		return nil, err
	}

	var ossStruct common.OssStruct
	if err = json.Unmarshal([]byte(resp), &ossStruct); err != nil {
		bds.logger.Errorf("resp:%s, err:%v", resp, err)
		return nil, err
	}
	bds.logger.Infof("sub device stor config resp: %+v", ossStruct)
	if !ossStruct.Success {
		return nil, errors.New("get sub device stor config failure")
	}
	if len(ossStruct.Result) >= 1 {
		bds.ossManager.setSubDevStorConfig(cid, &ossStruct.Result[0])
		return &ossStruct.Result[0], nil
	}
	return nil, errors.New("get sub device stor config result is empty")
}

func checkHttpResp(resp []byte) error {
	var atopResp common.AtopResp
	json.Unmarshal(resp, &atopResp)
	if atopResp.Success == true {
		return nil
	}

	errorCode := atopResp.ErrorCode
	switch errorCode {
	//Atop网关通用错误码
	case common.ATOP_ILLEGAL_ERR, common.ATOP_UNKNOW_ERR, common.ATOP_LIMIT_ERR:
		return fmt.Errorf("errorCode:%s", atopResp.ErrorCode)
	case common.ATOP_INTERNAL_ERR: //Dubbo 业务提供者返回错误码:
		return fmt.Errorf("errorCode:%s", atopResp.ErrorCode)
	default:
		return nil
	}
}

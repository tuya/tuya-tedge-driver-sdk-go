package service

import (
	"context"
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/driversvc"
	"github.com/tuya/tuya-tedge-driver-sdk-go/thingmodel"
)

var _ TyModelApi = (*TyDriverService)(nil)

type TyDriverService struct {
	tyService *driversvc.TyModelService
}

func NewTyService(l commons.TedgeLogger) *TyDriverService {
	baseService := driversvc.NewBaseService(l)
	return &TyDriverService{
		tyService: driversvc.NewTyModelService(baseService),
	}
}

func NewTyServiceWithBase(bds *driversvc.BaseService) *TyDriverService {
	return &TyDriverService{
		tyService: driversvc.NewTyModelService(bds),
	}
}

func (tmds *TyDriverService) Start(driver thingmodel.ThingModelDriver, opts ...driversvc.Option) error {
	return tmds.tyService.Start(driver, opts...)
}

// GetLogger 获取sdk的logger
func (tmds *TyDriverService) GetLogger() commons.TedgeLogger {
	return tmds.tyService.GetLogger()
}

// GetCustomConfig 获取驱动自定义配置
func (tmds *TyDriverService) GetCustomConfig() map[string]interface{} {
	return tmds.tyService.GetDriverConfig()
}

// GetGatewayInfo 获取网关信息
func (tmds *TyDriverService) GetGatewayInfo() commons.GatewayInfo {
	return tmds.tyService.GetGatewayMeta()
}

// ActiveDevice 新增并激活一个子设备
func (tmds *TyDriverService) ActiveDevice(device commons.TMDeviceMeta) error {
	return tmds.tyService.ActiveLinkDevice(device)
}

// AddDevice 新增一个子设备，但不激活到云端
func (tmds *TyDriverService) AddDevice(device commons.TMDeviceMeta) error {
	return tmds.tyService.AddLinkDevice(device)
}

// DeleteDevice 删除一个子设备
func (tmds *TyDriverService) DeleteDevice(cid string) error {
	return tmds.tyService.DeleteDevice(cid)
}

// ReportDeviceStatus 上报子设备状态
func (tmds *TyDriverService) ReportDeviceStatus(data *commons.DeviceStatus) error {
	return tmds.tyService.ReportDeviceStatus(data)
}

// ReportDevEvent
func (tmds *TyDriverService) ReportDevEvent(cid, event_type, deviceAddr, content string) error {
	return tmds.tyService.ReportDevEvent(cid, event_type, deviceAddr, content)
}

// SetDeviceExtendProperty 更新子设备附加属性，重复调用会覆盖之前的值
// property 示例:
// extenstion := make(map[string]interface{})
// extenstion["password"] = "12345678"
// property1 := devicemodel.ExtendedProperty{
// 	  InstallLocation: "XX东门",
// 	  ExtendData: extenstion,
// }
func (tmds *TyDriverService) SetDeviceExtendProperty(cid string, property commons.ExtendedProperty) error {
	return tmds.tyService.SetDeviceProperty(cid, property)
}

// SetDeviceBaseAttr 修改子设备名、IP地址、坐标 //比如：子设备名变化时，调用该接口修改，重复调用会覆盖之前的值
func (tmds *TyDriverService) SetDeviceBaseAttr(cid string, baseAttr commons.BaseProperty) error {
	return tmds.tyService.SetDeviceBaseAttr(cid, baseAttr)
}

////////////////////////////////////////////////////////////////////////////////////////////////
// PropertySetResponse 设备属性下发响应  设备在处理属性设置之后默认不对属性设置消息进行响应，除非云端在设置消息中显式指定ack:1
func (tmds *TyDriverService) PropertySetResponse(cid string, data thingmodel.CommonResponse) error {
	return tmds.tyService.PropertySetResponse(cid, data)
}

// PropertyGetResponse 设备属性查询响应
func (tmds *TyDriverService) PropertyGetResponse(cid string, data thingmodel.PropertyGetResponse) error {
	return tmds.tyService.PropertyGetResponse(cid, data)
}

// ActionExecuteResponse 设备动作执行响应
func (tmds *TyDriverService) ActionExecuteResponse(cid string, data thingmodel.ActionExecuteResponse) error {
	return tmds.tyService.ActionExecuteResponse(cid, data)
}

// PropertyReport 物模型属性上报 如果data参数中的Sys.Ack设置为1，则该方法会同步阻塞等待云端返回结果。
// 如非必要，不建议设置Sys.Ack
func (tmds *TyDriverService) PropertyReport(cid string, data thingmodel.PropertyReport) (thingmodel.CommonResponse, error) {
	return tmds.tyService.PropertyReport(cid, data)
}

// EventReport 物模型事件上报 如果data参数中的Sys.Ack设置为1，则该方法会同步阻塞等待云端返回结果。
//// 如非必要，不建议设置Sys.Ack
func (tmds *TyDriverService) EventReport(cid string, data thingmodel.EventReport) (thingmodel.CommonResponse, error) {
	return tmds.tyService.EventReport(cid, data)
}

// BatchReport 设备批量上报属性和事件 如果data参数中的Sys.Ack设置为1，则该方法会同步阻塞等待云端返回结果。
//// 如非必要，不建议设置Sys.Ack
func (tmds *TyDriverService) BatchReport(cid string, data thingmodel.BatchReport) (thingmodel.CommonResponse, error) {
	return tmds.tyService.BatchReport(cid, data)
}

// PropertyDesiredGet 设备拉取属性期望值 如果data参数中的Sys.Ack设置为1，则该方法会同步阻塞等待云端返回结果。
//// 如非必要，不建议设置Sys.Ack
func (tmds *TyDriverService) PropertyDesiredGet(cid string, data thingmodel.PropertyDesiredGet) (thingmodel.PropertyDesiredGetResponse, error) {
	return tmds.tyService.PropertyDesiredGet(cid, data)
}

// PropertyDesiredDelete 设备删除属性期望值 如果data参数中的Sys.Ack设置为1，则该方法会同步阻塞等待云端返回结果。
//// 如非必要，不建议设置Sys.Ack
func (tmds *TyDriverService) PropertyDesiredDelete(cid string, data thingmodel.PropertyDesiredDelete) (thingmodel.PropertyDesiredDeleteResponse, error) {
	return tmds.tyService.PropertyDesiredDelete(cid, data)
}

////////////////////////////////////////////////////////////////////////////////////////////////
// AllDevices 获取该驱动下的设备
func (tmds *TyDriverService) AllDevices() map[string]commons.TMDeviceInfo {
	return tmds.tyService.AllTMDevices()
}

// GetActiveDevices 获取已激活的设备列表
func (tmds *TyDriverService) GetActiveDevices() map[string]commons.TMDeviceInfo {
	return tmds.tyService.GetLinkActiveDevices()
}

// GetDeviceById 通过设备ID查询设备信息
func (tmds *TyDriverService) GetDeviceById(cid string) (commons.TMDeviceInfo, bool) {
	return tmds.tyService.GetTMDeviceById(cid)
}

// AddProduct 添加物模型产品
func (tmds *TyDriverService) AddProduct(pr thingmodel.AddProductReq) error {
	return tmds.tyService.AddProduct(pr)
}

// AllProducts 获取当前实例下的所有产品
func (tmds *TyDriverService) AllProducts() map[string]thingmodel.ThingModelProduct {
	return tmds.tyService.AllProducts()
}

// GetProductById 根据产品ID获取产品信息
func (tmds *TyDriverService) GetProductById(pid string) (thingmodel.ThingModelProduct, bool) {
	return tmds.tyService.GetProductById(pid)
}

func (tmds *TyDriverService) GetEventByPid(pid string) (map[string]thingmodel.Event, bool) {
	return tmds.tyService.GetEventByPid(pid)
}

// ReportAlert 驱动上报告警信息
func (tmds *TyDriverService) ReportAlert(ctx context.Context, level commons.AlertLevel, content string) error {
	return tmds.tyService.ReportAlert(ctx, level, content)
}

//随机id生成，32位
func (tmds *TyDriverService) GenRandomId() string {
	return tmds.tyService.GenRandomId()
}

// GetServiceId 获取驱动实例Id
func (tmds *TyDriverService) GetServiceId() string {
	return tmds.tyService.GetServiceId()
}

////////////////////////////////////////////////////////////////////////////////////////////
// ReportThroughHttp 通过http上报设备数据
func (tmds *TyDriverService) ReportThroughHttp(api, version string, payload map[string]interface{}) (string, error) {
	return tmds.tyService.ReportThroughHttp(api, version, payload)
}

// CmdRespSuccess 上报执行结果成功
func (tmds *TyDriverService) CmdRespSuccess(sn int64) error {
	return tmds.tyService.CmdResultUpload(sn, 1, "ok")
}

// CmdRespFail 上报执行结果失败
func (tmds *TyDriverService) CmdRespFail(sn int64, message string) error {
	return tmds.tyService.CmdResultUpload(sn, 0, message)
}

// UploadImage 新版文件上传
// 返回值：id, {"bucket":"ty-cn-storage30","objectKey":"/88012e-34125598-194c9d2a7048680e/v","secretKey":"c86440d1912a413d8d90d13097391159","expireTime":1551083405}
func (tmds *TyDriverService) UploadFileV2(cid, fileName string, content []byte, timeout int32) (string, string, error) {
	return tmds.tyService.UploadFileV2(cid, fileName, content, timeout)
}

// UploadFile 图片上传
func (tmds *TyDriverService) UploadFile(content []byte, fileName, subjectType string, timeout int32) (string, error) {
	return tmds.tyService.UploadFile(content, fileName, subjectType, timeout)
}

// HttpRequestProxy 物模型下，直接调用atop接口 //d.json
func (tmds *TyDriverService) HttpRequestProxy(params thingmodel.HttpRequestParam, payload map[string]interface{}, timeout int) (string, error) {
	return tmds.tyService.HttpRequestProxy(params, payload, timeout)
}

func (tmds *TyDriverService) GenHttpProxyParam(url, api, version string) thingmodel.HttpRequestParam {
	proxyParam := thingmodel.HttpRequestParam{
		Url:     url,
		Api:     api,
		Version: version,
	}

	return proxyParam
}

////////////////////////////////////////////////////////////////////////////////////////////
//mqtt 接口
func (tmds *TyDriverService) Publish(topic string, qos byte, retained bool, message []byte) error {
	if tmds.tyService.MqttClient != nil {
		return tmds.tyService.MqttClient.Publish(topic, qos, retained, message)
	}
	return errors.New("mqtt client is nil")
}

func (tmds *TyDriverService) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
	if tmds.tyService.MqttClient != nil {
		return tmds.tyService.MqttClient.Subscribe(topic, qos, handler)
	}
	return errors.New("mqtt client is nil")
}

func (tmds *TyDriverService) UnSubscribe(topic string) error {
	if tmds.tyService.MqttClient != nil {
		return tmds.tyService.MqttClient.UnSubscribe(topic)
	}
	return errors.New("mqtt client is nil")
}

////////////////////////////////////////////////////////////////////////////////////////////
// GetBackupKV 根据key值获取驱动存储的自定义内容，自定义内容支持云端备份
func (tmds *TyDriverService) GetBackupKV(keys []string) (map[string][]byte, error) {
	if len(keys) <= 0 {
		return nil, errors.New("required keys")
	}
	return tmds.tyService.GetBackupKV(keys)
}

// GetBackupKV 根据key获取KV存储，支持云端备份
func (tmds *TyDriverService) GetBackupKVOne(key string) ([]byte, error) {
	kvs, err := tmds.tyService.GetBackupKV([]string{key})
	if err != nil {
		return nil, err
	}
	return kvs[key], nil
}

// PutBackupKV 存储驱动的自定义内容，自定义内容支持云端备份
func (tmds *TyDriverService) PutBackupKV(kvs map[string][]byte) error {
	if len(kvs) <= 0 {
		return errors.New("required key value")
	}
	return tmds.tyService.PutBackupKV(kvs)
}

// PutBackupKVOne 更新KV存储，支持云端备份
func (tmds *TyDriverService) PutBackupKVOne(key string, value []byte) error {
	kvs := map[string][]byte{
		key: value,
	}
	return tmds.tyService.PutBackupKV(kvs)
}

// DelBackupKV 根据key值删除驱动存储的自定义内容，删除的Key Value支持云端备份
func (tmds *TyDriverService) DelBackupKV(keys []string) error {
	if len(keys) <= 0 {
		return errors.New("required keys")
	}
	return tmds.tyService.DelBackupKV(keys)
}

// GetBackupKVKeys 根据前缀获取所有key，获取的key支持云端备份
func (tmds *TyDriverService) GetBackupKVKeys(prefix string) ([]string, error) {
	return tmds.tyService.GetBackupKVKeys(prefix)
}

// GetAllCustomStorage 获取所有驱动存储的自定义内容，内容支持云端备份
func (tmds *TyDriverService) QueryBackupKV(prefix string) (map[string][]byte, error) {
	return tmds.tyService.QueryBackupKV(prefix)
}

// GetKV 根据key值获取驱动存储的自定义内容，不支持云端备份
func (tmds *TyDriverService) GetKV(keys []string) (map[string][]byte, error) {
	if len(keys) <= 0 {
		return nil, errors.New("required keys")
	}
	return tmds.tyService.GetKV(keys)
}

// GetKVOne 根据key获取内容，不支持云端备份
func (tmds *TyDriverService) GetKVOne(key string) ([]byte, error) {
	kvs, err := tmds.tyService.GetKV([]string{key})
	if err != nil {
		return nil, err
	}
	return kvs[key], nil
}

// PutKV 存储驱动的自定义内容，不支持云端备份
// input：map[string][]byte
func (tmds *TyDriverService) PutKV(kvs map[string][]byte) error {
	if len(kvs) <= 0 {
		return errors.New("required key value")
	}
	return tmds.tyService.PutKv(kvs)
}

// PutKVOne 存储驱动的自定义内容，不支持云端备份
func (tmds *TyDriverService) PutKVOne(key string, value []byte) error {
	kvs := map[string][]byte{
		key: value,
	}
	return tmds.tyService.PutKv(kvs)
}

// DeleteKV 根据key值删除驱动存储的自定义内容，不支持云端备份
func (tmds *TyDriverService) DeleteKV(keys []string) error {
	if len(keys) <= 0 {
		return errors.New("required keys")
	}
	return tmds.tyService.DeleteKV(keys)
}

// QueryKV，根据前缀查询所有KV
func (tmds *TyDriverService) QueryKV(prefix string) (map[string][]byte, error) {
	return tmds.tyService.QueryKV(prefix)
}

// 根据前缀筛选key，传空则返回所有key，不支持云端备份
func (tmds *TyDriverService) GetKVKeys(prefix string) ([]string, error) {
	return tmds.tyService.GetKVKeys(prefix)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
func (tmds *TyDriverService) SendToApp(appName string, data commons.AppDriverReq, cnnNum ...int) (commons.Response, error) {
	ctx, c := context.WithTimeout(context.Background(), time.Second*5)
	defer c()
	return tmds.tyService.SendToApp(ctx, appName, data, cnnNum...)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
func (tmds *TyDriverService) RegistDriverProxy(proxyInfo commons.ProxyInfo) error {
	return tmds.tyService.DriverProxyRegist(proxyInfo.Host, proxyInfo.Port)
}

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/driversvc"
)

func NewBaseService(l commons.TedgeLogger) *driversvc.BaseService {
	return driversvc.NewBaseService(l)
}

var _ DpModelApi = (*DPDriverService)(nil)

type DPDriverService struct {
	dpService *driversvc.DPModelService
}

func NewDPService(l commons.TedgeLogger) *DPDriverService {
	baseService := driversvc.NewBaseService(l)
	return &DPDriverService{
		dpService: driversvc.NewDPModelService(baseService),
	}
}

func NewDPServiceWithBase(bds *driversvc.BaseService) *DPDriverService {
	return &DPDriverService{
		dpService: driversvc.NewDPModelService(bds),
	}
}

func (dds *DPDriverService) Start(driver dpmodel.DPModelDriver, opts ...driversvc.Option) error {
	return dds.dpService.Start(driver, opts...)
}

// GetTEdgeModel 获取TEdge运行模式
func (dds *DPDriverService) GetTEdgeModel() commons.RunningModel {
	return dds.GetTEdgeModel()
}

// GetLogger 获取sdk的logger
func (dds *DPDriverService) GetLogger() commons.TedgeLogger {
	return dds.dpService.BaseService.GetLogger()
}

// GetCustomConfig 获取驱动自定义配置
func (dds *DPDriverService) GetCustomConfig() map[string]interface{} {
	return dds.dpService.BaseService.GetDriverConfig()
}

// ActiveDevice 新增并激活一个子设备
func (dds *DPDriverService) ActiveDevice(device commons.DeviceMeta) error {
	return dds.dpService.ActiveDevice(device, false)
}

// ActiveIPCDevice IPC子设备激活专用接口，IPC子设备需要单独设置能力集，该接口默认会设置IPC通用能力集
func (dds *DPDriverService) ActiveIPCDevice(device commons.DeviceMeta) error {
	return dds.dpService.ActiveDevice(device, true)
}

// SetDeviceExtendProperty 更新子设备附加属性，重复调用会覆盖之前的值
// property 示例:
// extenstion := make(map[string]interface{})
// extenstion["password"] = "12345678"
// property1 := devicemodel.ExtendedProperty{
// 	  VendorCode: "yufan",
// 	  InstallLocation: "华策东门",
// 	  ExtendData: extenstion,
// }
func (dds *DPDriverService) SetDeviceExtendProperty(cid string, property commons.ExtendedProperty) error {
	return dds.dpService.SetDeviceProperty(cid, property)
}

// SetDeviceBaseAttr 修改子设备名、IP地址、坐标 //比如：子设备名变化时，调用该接口修改，重复调用会覆盖之前的值
func (dds *DPDriverService) SetDeviceBaseAttr(cid string, baseAttr commons.BaseProperty) error {
	return dds.dpService.SetDeviceBaseAttr(cid, baseAttr)
}

// AddDevice 新增一个子设备，但不激活到云端 //TODO:特殊场景需要，比如对接涂鸦自己的zigbee网关
func (dds *DPDriverService) AddDevice(device commons.DeviceMeta) error {
	return dds.dpService.AddDevice(device, false)
}

// DeleteDevice 删除一个子设备
func (dds *DPDriverService) DeleteDevice(cid string) error {
	return dds.dpService.DeleteDevice(cid)
}

// ReportDeviceStatus 上报子设备状态
func (dds *DPDriverService) ReportDeviceStatus(data *commons.DeviceStatus) error {
	return dds.dpService.ReportDeviceStatus(data)
}

// ReportWithDPData 上报带dp点的数据
// 支持多个dp点同时上报
// dataType为以下其中之一：
//const (
//	ValueType  = "value"
//	BoolType   = "bool"
//	StringType = "string"
//	RawType    = "raw"
//	EnumType   = "enum"
//	FaultType  = "fault"
//)
// 上报数值类型数据时，data中的value类型需要显式指定为int64
// 上报枚举类型数据时，data中的value为dp点定义的枚举列表中的一个值，类型：string
// 上报故障类型数据时，data中的value为dp点定义的故障列表的子集，类型为：[]string或者int32
func (dds *DPDriverService) ReportWithDPData(cid string, data []*dpmodel.WithDPValue) error {
	return dds.dpService.ReportWithDPData(cid, data)
}

// ReportWithoutDPData 上报不带dp点的数据
func (dds *DPDriverService) ReportWithoutDPData(data *dpmodel.WithoutDPValue) error {
	return dds.dpService.ReportWithoutDPData(fmt.Sprintf("%s%s", "smart/device/out/", dds.dpService.GetGatewayId()), data)
}

// ReportWithoutDPDataWithTopic 向指定的topic上报不带dp点的数据
func (dds *DPDriverService) ReportWithoutDPDataWithTopic(topic string, data *dpmodel.WithoutDPValue) error {
	return dds.dpService.ReportWithoutDPData(topic, data)
}

// ReportThroughHttp 通过http上报设备数据
func (dds *DPDriverService) ReportThroughHttp(api, version string, payload map[string]interface{}) (string, error) {
	return dds.dpService.ReportThroughHttp(api, version, payload)
}

// AllDevices 获取该驱动下的设备
func (dds *DPDriverService) AllDevices() map[string]commons.DeviceInfo {
	return dds.dpService.AllDevices()
}

// GetActiveDevices 获取已激活的设备列表
func (dds *DPDriverService) GetActiveDevices() map[string]commons.DeviceInfo {
	return dds.dpService.GetActiveDevices()
}

// GetDeviceById 通过设备ID查询设备信息
func (dds *DPDriverService) GetDeviceById(cid string) (commons.DeviceInfo, bool) {
	return dds.dpService.GetDeviceById(cid)
}

// GetActiveDeviceById 通过设备ID查询已激活的设备信息
// return true: 设备已激活，false: 设备未激活或者设备不存在
func (dds *DPDriverService) GetActiveDeviceById(cid string) (commons.DeviceInfo, bool) {
	return dds.dpService.GetActiveDeviceById(cid)
}

// AddProduct 添加一个产品。
func (dds *DPDriverService) AddProduct(product dpmodel.DPModelProductAddInfo) error {
	return dds.dpService.AddProduct(product)
}

// AllProducts 获取驱动下的产品信息
func (dds *DPDriverService) AllProducts() map[string]dpmodel.DPModelProduct {
	return dds.dpService.AllProducts()
}

// GetProductById 根据产品ID查询产品信息
func (dds *DPDriverService) GetProductById(pid string) (dpmodel.DPModelProduct, bool) {
	return dds.dpService.GetProductById(pid)
}

// ReportAlert 驱动告警，告警内容会在tedge web前端页面展示
func (dds *DPDriverService) ReportAlert(ctx context.Context, level commons.AlertLevel, content string) error {
	return dds.dpService.ReportAlert(ctx, level, content)
}

// GetGatewayInfo 获取网关信息
func (dds *DPDriverService) GetGatewayInfo() commons.GatewayInfo {
	return dds.dpService.GetGatewayMeta()
}

// GetServiceId 获取驱动实例Id
func (dds *DPDriverService) GetServiceId() string {
	return dds.dpService.GetServiceId()
}

////////////////////////////////////////////////////////////////////////////////////////////////
// CmdRespSuccess 上报执行结果成功
func (dds *DPDriverService) CmdRespSuccess(sn int64) error {
	return dds.dpService.CmdResultUpload(sn, 1, "ok")
}

// CmdRespFail 上报执行结果失败
func (dds *DPDriverService) CmdRespFail(sn int64, message string) error {
	return dds.dpService.CmdResultUpload(sn, 0, message)
}

// UploadFile 文件上传
func (dds *DPDriverService) UploadFile(content []byte, fileName, subjectType string, timeout int32) (string, error) {
	return dds.dpService.UploadFile(content, fileName, subjectType, timeout)
}

// UploadImage 新版文件上传
// 返回值：id, {"bucket":"ty-cn-storage30","objectKey":"/88012e-34125598-194c9d2a7048680e/v","secretKey":"c86440d1912a413d8d90d13097391159","expireTime":1551083405}
func (dds *DPDriverService) UploadFileV2(cid, fileName string, content []byte, timeout int32) (string, string, error) {
	return dds.dpService.UploadFileV2(cid, fileName, content, timeout)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// IpcDeviceUpdateSkill 激活IPC子设备后，单独更新IPC"通用"能力集场景时使用，比如IPC子设备已经激活，需要补充能力集
//Deprecated
func (dds *DPDriverService) IpcDeviceUpdateSkill(version string, cid string) (string, error) {
	return dds.dpService.IpcDeviceUpdateSkill(version, cid)
}

////////////////////////////////////////////////////////////////////////////////////////////////
//Mqtt api
func (dds *DPDriverService) Publish(topic string, qos byte, retained bool, message []byte) error {
	if dds.dpService.MqttClient != nil {
		return dds.dpService.MqttClient.Publish(topic, qos, retained, message)
	}
	return errors.New("mqtt client is nil")
}

func (dds *DPDriverService) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
	if dds.dpService.MqttClient != nil {
		return dds.dpService.MqttClient.Subscribe(topic, qos, handler)
	}
	return errors.New("mqtt client is nil")
}

func (dds *DPDriverService) UnSubscribe(topic string) error {
	if dds.dpService.MqttClient != nil {
		return dds.dpService.MqttClient.UnSubscribe(topic)
	}
	return errors.New("mqtt client is nil")
}

////////////////////////////////////////////////////////////////////////////////////////////////
// GetKV 根据key值获取驱动存储的自定义内容，不支持云端备份
func (dds *DPDriverService) GetKV(keys []string) (map[string][]byte, error) {
	if len(keys) <= 0 {
		return nil, errors.New("required keys")
	}
	return dds.dpService.GetKV(keys)
}

// GetKVOne 根据key获取内容，不支持云端备份
func (dds *DPDriverService) GetKVOne(key string) ([]byte, error) {
	kvs, err := dds.dpService.GetKV([]string{key})
	if err != nil {
		return nil, err
	}
	return kvs[key], nil
}

// PutKV 存储驱动的自定义内容，不支持云端备份
func (dds *DPDriverService) PutKV(kvs map[string][]byte) error {
	if len(kvs) <= 0 {
		return errors.New("required key value")
	}
	return dds.dpService.PutKv(kvs)
}

// PutKVOne 存储驱动的自定义内容，不支持云端备份
func (dds *DPDriverService) PutKVOne(key string, value []byte) error {
	kvs := map[string][]byte{
		key: value,
	}
	return dds.dpService.PutKv(kvs)
}

// DeleteKV 根据key值删除驱动存储的自定义内容，不支持云端备份
func (dds *DPDriverService) DeleteKV(keys []string) error {
	if len(keys) <= 0 {
		return errors.New("required keys")
	}
	return dds.dpService.DeleteKV(keys)
}

// GetKVKeys 根据前缀筛选key，传空则返回所有key，不支持云端备份
func (dds *DPDriverService) GetKVKeys(prefix string) ([]string, error) {
	return dds.dpService.GetKVKeys(prefix)
}

// QueryKV 根据前缀搜索KV存储，传空返回所有的结果，不支持云端备份
func (dds *DPDriverService) QueryKV(prefix string) (map[string][]byte, error) {
	return dds.dpService.QueryKV(prefix)
}

// QueryBackupKV 根据前缀搜索KV存储，传空返回所有的结果，支持云端备份
func (dds *DPDriverService) QueryBackupKV(prefix string) (map[string][]byte, error) {
	return dds.dpService.QueryBackupKV(prefix)
}

// GetBackupKV 根据key获取KV存储，支持云端备份
func (dds *DPDriverService) GetBackupKV(keys []string) (map[string][]byte, error) {
	if len(keys) <= 0 {
		return nil, errors.New("required keys")
	}
	return dds.dpService.GetBackupKV(keys)
}

// GetBackupKV 根据key获取KV存储，支持云端备份
func (dds *DPDriverService) GetBackupKVOne(key string) ([]byte, error) {
	kvs, err := dds.dpService.GetBackupKV([]string{key})
	if err != nil {
		return nil, err
	}
	return kvs[key], nil
}

// GetBackupKVKeys 根据前缀获取keys，传空返回所有的结果，支持云端备份
func (dds *DPDriverService) GetBackupKVKeys(prefix string) ([]string, error) {
	return dds.dpService.GetBackupKVKeys(prefix)
}

// PutBackupKV 更新KV存储，支持云端备份
func (dds *DPDriverService) PutBackupKV(kvs map[string][]byte) error {
	if len(kvs) <= 0 {
		return errors.New("required key value")
	}
	return dds.dpService.PutBackupKV(kvs)
}

// PutBackupKVOne 更新KV存储，支持云端备份
func (dds *DPDriverService) PutBackupKVOne(key string, value []byte) error {
	kvs := map[string][]byte{
		key: value,
	}
	return dds.dpService.PutBackupKV(kvs)
}

// DelBackupKV 删除KV存储，支持云端备份
func (dds *DPDriverService) DelBackupKV(keys []string) error {
	if len(keys) <= 0 {
		return errors.New("required keys")
	}
	return dds.dpService.DelBackupKV(keys)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// 向App应用发数据
func (bds *DPDriverService) SendToApp(appName string, req commons.AppDriverReq, cnnNum ...int) (commons.Response, error) {
	ctx, c := context.WithTimeout(context.Background(), time.Second*5)
	defer c()
	return bds.dpService.SendToApp(ctx, appName, req, cnnNum...)

}
//////////////////////////////////////////////////////////////////////////////////////////////
//DP点上报：支持断网续传
func (dds *DPDriverService) ReportWithDPDataV2(cid string, data []*dpmodel.WithDPValue) error {
	return dds.dpService.ReportWithDPDataV2(cid, data)
}

//打印因断网而存储在本地的DP点数据
func (dds *DPDriverService) PrintUnPublishedDPS() error {
	if dds.dpService.RtsManager != nil {
		return dds.dpService.RtsManager.PrintDPKeys()
	}

	return fmt.Errorf("rts not enable")
}

//Atop 失败自动重传
func (dds *DPDriverService) ReportHttpWithRetrans(api, version string, payload map[string]interface{}) (string, error) {
	return dds.dpService.ReportHttpWithRetrans(api, version, payload)
}

//打印上传失败存储在本地的Atop数据
func (dds *DPDriverService) PrintUnReportedAtops() error {
	if dds.dpService.RtsManager != nil {
		return dds.dpService.RtsManager.PrintAtopKeys()
	}

	return fmt.Errorf("rts not enable")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
//驱动Web服务注册
func (dds *DPDriverService) RegistDriverProxy(proxyInfo commons.ProxyInfo) error {
	return dds.dpService.DriverProxyRegist(proxyInfo.Host, proxyInfo.Port)
}

package service

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/driversvc"
	"github.com/tuya/tuya-tedge-driver-sdk-go/thingmodel"
)

type DpModelApi interface {
	DpBaseApi
	CommonBaseApi

	DriverMqtt
	KVDatabase
	DriverProxy

	DPRetransmit
}

type TyModelApi interface {
	TyBaseApi
	CommonBaseApi

	DriverMqtt
	KVDatabase
	DriverProxy
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
type DpBaseApi interface {
	Start(driver dpmodel.DPModelDriver, opts ...driversvc.Option) error

	// ActiveDevice 新增并激活一个子设备
	ActiveDevice(device commons.DeviceMeta) error

	// ActiveIPCDevice IPC子设备激活专用接口，IPC子设备需要单独设置能力集，该接口默认会设置IPC通用能力集
	ActiveIPCDevice(device commons.DeviceMeta) error

	// AddDevice 新增一个子设备，但不激活到云端
	AddDevice(device commons.DeviceMeta) error

	// ReportWithDPData 上报带dp点的数据
	ReportWithDPData(cid string, data []*dpmodel.WithDPValue) error

	// ReportWithoutDPData 上报不带dp点的数据
	ReportWithoutDPData(data *dpmodel.WithoutDPValue) error

	// ReportWithoutDPDataWithTopic 向指定的topic上报不带dp点的数据
	ReportWithoutDPDataWithTopic(topic string, data *dpmodel.WithoutDPValue) error

	// AllDevices 获取该驱动下的设备
	AllDevices() map[string]commons.DeviceInfo

	// GetActiveDevices 获取已激活的设备列表
	GetActiveDevices() map[string]commons.DeviceInfo

	// GetDeviceById 通过设备ID查询设备信息
	GetDeviceById(cid string) (commons.DeviceInfo, bool)

	// GetActiveDeviceById 通过设备ID查询已激活的设备信息
	// return true: 设备已激活，false: 设备未激活或者设备不存在
	GetActiveDeviceById(cid string) (commons.DeviceInfo, bool)

	// AddProduct 添加一个产品。
	AddProduct(product dpmodel.DPModelProductAddInfo) error

	// AllProducts 获取驱动下的产品信息
	AllProducts() map[string]dpmodel.DPModelProduct

	// GetProductById 根据产品ID查询产品信息
	GetProductById(pid string) (dpmodel.DPModelProduct, bool)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
type CommonBaseApi interface {
	//GetLogger 获取 TEdge logger 接口
	GetLogger() commons.TedgeLogger

	// GetCustomConfig 获取驱动自定义配置文件
	GetCustomConfig() map[string]interface{}

	// GetGatewayInfo 获取网关信息
	GetGatewayInfo() commons.GatewayInfo

	// SetDeviceExtendProperty 更新子设备附加属性，重复调用会覆盖之前的值
	SetDeviceExtendProperty(cid string, property commons.ExtendedProperty) error

	// SetDeviceBaseAttr 修改子设备名、IP地址、坐标 //比如：子设备名变化时，调用该接口修改，重复调用会覆盖之前的值
	SetDeviceBaseAttr(cid string, baseAttr commons.BaseProperty) error

	// DeleteDevice 删除/解绑一个子设备
	DeleteDevice(cid string) error

	// ReportDeviceStatus 上报子设备状态：在线、离线
	ReportDeviceStatus(data *commons.DeviceStatus) error

	// ReportAlert 驱动告警，告警内容会在tedge web前端页面展示
	ReportAlert(ctx context.Context, level commons.AlertLevel, content string) error

	////////////////////////////////////////////////////////////////////////////////////////////////
	// ReportThroughHttp 通过http上报设备数据
	ReportThroughHttp(api, version string, payload map[string]interface{}) (string, error)

	// CmdRespSuccess 上报执行结果成功
	CmdRespSuccess(sn int64) error

	// CmdRespFail 上报执行结果失败
	CmdRespFail(sn int64, message string) error

	// UploadFile 文件上传
	UploadFile(content []byte, fileName, subjectType string, timeout int32) (string, error)

	// UploadImage 新版文件上传
	// 返回值：id, {"bucket":"ty-cn-storage30","objectKey":"/88012e-34125598-194c9d2a7048680e/v","secretKey":"c86440d1912a413d8d90d13097391159","expireTime":1551083405}
	UploadFileV2(cid, fileName string, content []byte, timeout int32) (string, string, error)

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// 向App应用发数据
	SendToApp(appName string, req commons.AppDriverReq, cnnNum ...int) (commons.Response, error)

	// GetServiceId 获取驱动实例Id
	GetServiceId() string
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
type DriverMqtt interface {
	Publish(topic string, qos byte, retained bool, message []byte) error

	Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error

	UnSubscribe(topic string) error
}

type KVDatabase interface {
	////////////////////////////////////////////////////////////////////////////////////////////////
	// GetKV 根据key值获取驱动存储的自定义内容，不支持云端备份
	GetKV(keys []string) (map[string][]byte, error)

	// GetKVOne 根据key获取内容，不支持云端备份
	GetKVOne(key string) ([]byte, error)

	// PutKV 存储驱动的自定义内容，不支持云端备份
	PutKV(kvs map[string][]byte) error

	// PutKVOne 存储驱动的自定义内容，不支持云端备份
	PutKVOne(key string, value []byte) error

	// DeleteKV 根据key值删除驱动存储的自定义内容，不支持云端备份
	DeleteKV(keys []string) error

	// GetKVKeys 根据前缀筛选key，传空则返回所有key，不支持云端备份
	GetKVKeys(prefix string) ([]string, error)

	// QueryKV 根据前缀搜索KV存储，传空返回所有的结果，不支持云端备份
	QueryKV(prefix string) (map[string][]byte, error)

	// QueryBackupKV 根据前缀搜索KV存储，传空返回所有的结果，支持云端备份
	QueryBackupKV(prefix string) (map[string][]byte, error)

	// GetBackupKV 根据key获取KV存储，支持云端备份
	GetBackupKV(keys []string) (map[string][]byte, error)

	// GetBackupKV 根据key获取KV存储，支持云端备份
	GetBackupKVOne(key string) ([]byte, error)

	// GetBackupKVKeys 根据前缀获取keys，传空返回所有的结果，支持云端备份
	GetBackupKVKeys(prefix string) ([]string, error)

	// PutBackupKV 更新KV存储，支持云端备份
	PutBackupKV(kvs map[string][]byte) error

	// PutBackupKVOne 更新KV存储，支持云端备份
	PutBackupKVOne(key string, value []byte) error

	// DelBackupKV 删除KV存储，支持云端备份
	DelBackupKV(keys []string) error
}

type DPRetransmit interface {
	//////////////////////////////////////////////////////////////////////////////////////////////
	//DP点上报：支持断网续传
	ReportWithDPDataV2(cid string, data []*dpmodel.WithDPValue) error

	//打印因断网而存储在本地的DP点数据
	PrintUnPublishedDPS() error

	//Atop 失败自动重传
	ReportHttpWithRetrans(api, version string, payload map[string]interface{}) (string, error)

	//打印上传失败存储在本地的Atop数据
	PrintUnReportedAtops() error
}

type DriverProxy interface {
	//驱动Web服务注册
	RegistDriverProxy(proxyInfo commons.ProxyInfo) error
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
type TyBaseApi interface {
	Start(driver thingmodel.ThingModelDriver, opts ...driversvc.Option) error
	ActiveDevice(device commons.TMDeviceMeta) error
	AddDevice(device commons.TMDeviceMeta) error
	AllDevices() map[string]commons.TMDeviceInfo
	GetActiveDevices() map[string]commons.TMDeviceInfo
	GetDeviceById(cid string) (commons.TMDeviceInfo, bool)

	AddProduct(pr thingmodel.AddProductReq) error
	AllProducts() map[string]thingmodel.ThingModelProduct
	GetProductById(pid string) (thingmodel.ThingModelProduct, bool)

	PropertySetResponse(cid string, data thingmodel.CommonResponse) error
	PropertyGetResponse(cid string, data thingmodel.PropertyGetResponse) error
	ActionExecuteResponse(cid string, data thingmodel.ActionExecuteResponse) error
	PropertyReport(cid string, data thingmodel.PropertyReport) (thingmodel.CommonResponse, error)
	EventReport(cid string, data thingmodel.EventReport) (thingmodel.CommonResponse, error)
	BatchReport(cid string, data thingmodel.BatchReport) (thingmodel.CommonResponse, error)
	PropertyDesiredGet(cid string, data thingmodel.PropertyDesiredGet) (thingmodel.PropertyDesiredGetResponse, error)
	PropertyDesiredDelete(cid string, data thingmodel.PropertyDesiredDelete) (thingmodel.PropertyDesiredDeleteResponse, error)
	GetEventByPid(pid string) (map[string]thingmodel.Event, bool)

	ReportDevEvent(cid, event_type, deviceAddr, content string) error
	HttpRequestProxy(params thingmodel.HttpRequestParam, payload map[string]interface{}, timeout int) (string, error)
	GenHttpProxyParam(url, api, version string) thingmodel.HttpRequestParam
	GenRandomId() string
}



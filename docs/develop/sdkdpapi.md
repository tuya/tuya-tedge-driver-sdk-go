#常用 SDK API(DP模型)

####创建驱动服务
NewBaseService 和 NewDPServiceWithBase 结合使用，或者单独使用 NewDPService
```golang
func NewBaseService(l commons.TedgeLogger) *driversvc.BaseService

func NewDPServiceWithBase(bds *driversvc.BaseService) *DPDriverService`

func NewDPService(l commons.TedgeLogger) *DPDriverService
```

####启动驱动服务 
```golang
func (dds *DPDriverService) Start(driver dpmodel.DPModelDriver, opts ...driversvc.Option) error
```

####获取日志接口
GetLogger 获取 sdk 的 logger 接口
```golang
func (dds *DPDriverService) GetLogger() commons.TedgeLogger
```

####获取驱动自定义配置
GetCustomConfig 获取驱动自定义配置
```golang
func (dds *DPDriverService) GetCustomConfig() map[string]interface{}
```

####新增并激活一个子设备
ActiveDevice 新增并激活一个子设备
```golang
func (dds *DPDriverService) ActiveDevice(device commons.DeviceMeta) error
```

####更新子设备附加属性
SetDeviceExtendProperty 更新子设备附加属性
```golang
// SetDeviceExtendProperty 更新子设备附加属性，重复调用会覆盖之前的值
// property 示例:
// extenstion := make(map[string]interface{})
// extenstion["password"] = "12345678"
// property := devicemodel.ExtendedProperty{
// 	  InstallLocation: "xxxxxxx",
// 	  ExtendData: extenstion,
// }
func (dds *DPDriverService) SetDeviceExtendProperty(cid string, property commons.ExtendedProperty) error
```

####更新子设备名
SetDeviceBaseAttr 修改子设备名
```golang
func (dds *DPDriverService) SetDeviceBaseAttr(cid string, baseAttr commons.BaseProperty) error
```

####上报子设备状态
ReportDeviceStatus 上报子设备状态
```golang
func (dds *DPDriverService) ReportDeviceStatus(data *commons.DeviceStatus) error
```

####上报设备的dp点值
ReportWithDPData 上报设备的dp点值
```golang
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
func (dds *DPDriverService) ReportWithDPData(cid string, data []*dpmodel.WithDPValue) error
```

####通过http上报数据
ReportThroughHttp 通过http上报数据，特殊数据或业务需要通过http接口上报
```golang
func (dds *DPDriverService) ReportThroughHttp(api, version string, payload map[string]interface{}) (string, error)
```

####获取该驱动下的设备
AllDevices 获取该驱动下的子设备，一般驱动初始化时调用
```golang
func (dds *DPDriverService) AllDevices() map[string]commons.DeviceInfo
```

####获取已激活的子设备列表
GetActiveDevices 获取已激活的子设备列表，一般驱动初始化时调用
```golang
func (dds *DPDriverService) GetActiveDevices() map[string]commons.DeviceInfo
```

####添加一个产品
AddProduct 添加一个产品
```golang
func (dds *DPDriverService) AddProduct(product dpmodel.DPModelProductAddInfo) error
```

####驱动上报告警
ReportAlert 驱动告警，告警内容会在tedge web前端告警中心展示
```golang
func (dds *DPDriverService) ReportAlert(ctx context.Context, level commons.AlertLevel, content string) error
```

####获取网关信息
GetGatewayInfo 获取网关信息
```golang
func (dds *DPDriverService) GetGatewayInfo() commons.GatewayInfo
```

####上报执行结果成功
CmdRespSuccess 上报执行结果成功(根据具体业务逻辑需求调用)
```golang
func (dds *DPDriverService) CmdRespSuccess(sn int64) error
```

####上报执行结果成功
CmdRespFail 上报执行结果失败(根据具体业务逻辑需求调用)
```golang
func (dds *DPDriverService) CmdRespFail(sn int64, message string) error
```

####上传文件或图片
UploadFile 上传文件或图片
```golang
func (dds *DPDriverService) UploadFile(content []byte, fileName, subjectType string, timeout int32) (string, error)
```

####上传文件或图片
UploadFileV2 上传文件或图片
```golang
// 返回值：id, {"bucket":"ty-cn-storage30","objectKey":"/88012e-34125598-194c9d2a7048680e/v","secretKey":"c86440d1912a413d8d90d13097391159","expireTime":1551083405}
func (dds *DPDriverService) UploadFileV2(cid, fileName string, content []byte, timeout int32) (string, string, error)
```

####mqtt publish
```golang
func (dds *DPDriverService) Publish(topic string, qos byte, retained bool, message []byte) error
```

####mqtt Subscribe
```golang
func (dds *DPDriverService) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error
```

####mqtt UnSubscribe
```golang
func (dds *DPDriverService) UnSubscribe(topic string) error
```

####驱动Web服务注册
RegistDriverProxy 驱动Web服务注册
```golang
func (dds *DPDriverService) RegistDriverProxy(proxyInfo commons.ProxyInfo) error
```

上一章：[驱动本地调试指南](./driverdebug.md)

package commons

type (
	// ProtocolProperties 子设备的自定义配置，通过config.json配置
	ProtocolProperties map[string]interface{}

	// DeviceMeta 新增子设备
	DeviceMeta struct {
		Cid          string                        `json:"cid"`          //子设备cid，网关下唯一
		ProductId    string                        `json:"productKey"`   //产品id
		BaseAttr     BaseProperty                  `json:"baseAttr"`     //设备属性
		ExtendedAttr ExtendedProperty              `json:"extendedAttr"` //设备属性
		Protocols    map[string]ProtocolProperties `json:"protocols"`    //设备自定义协议配置
	}

	// BaseProperty 子设备基础属性：子设备名、IP地址、经度、纬度
	BaseProperty struct {
		Name string `json:"name"` //子设备名
		Ip   string `json:"ip"`   //子设备IP地址
		Lat  string `json:"lat"`  //纬度
		Lon  string `json:"lon"`  //经度
	}

	// ExtendedProperty 子设备扩展属性
	ExtendedProperty struct {
		//VendorCode      string                 `json:"vendorCode"`      //设备厂商
		InstallLocation string                 `json:"installLocation"` //安装地址
		ExtendData      map[string]interface{} `json:"extendData"`      //扩展字段 //map[string]interface{}
	}

	// DeviceInfo 回调接口, 包含激活信息
	DeviceInfo struct {
		DeviceMeta           // 原始信息
		ActiveStatus  string // 设备激活状态
		OnLineStatus  string // 设备在线/离线状态
		CloudDeviceId string // 设备云端对应的🆔(如果设备已激活)
	}

	// DeviceStatus device status 设备状态上报
	DeviceStatus struct {
		Online  []string // 设备上线列表
		Offline []string // 设备下线列表
	}
)

type (
	// TMDeviceMeta 物模型设备添加
	TMDeviceMeta struct {
		Cid          string                        `json:"cid"`          //子设备cid，网关下唯一
		ProductId    string                        `json:"productKey"`   //产品pid
		BaseAttr     BaseProperty                  `json:"baseAttr"`     //设备属性
		ExtendedAttr ExtendedProperty              `json:"extendedAttr"` //设备属性
		Protocols    map[string]ProtocolProperties `json:"protocols"`    //设备自定义协议配置
	}

	// TMDeviceInfo 物模型设备信息
	TMDeviceInfo struct {
		TMDeviceMeta `json:",inline"`
		DeviceId     string `json:"deviceId"`     // 设备在云端的IotId
		ActiveStatus string `json:"activeStatus"` // 设备激活状态
		OnLineStatus string `json:"onlineStatus"` // 设备在线/离线状态
	}
)

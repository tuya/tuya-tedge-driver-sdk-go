package commons

/////////////////////////////////////////////////////////////////////////////////////
type SubDataType string

const (
	ValueTypeUint8   SubDataType = "uint8"
	ValueTypeUint16  SubDataType = "uint16"
	ValueTypeUint32  SubDataType = "uint32"
	ValueTypeUint64  SubDataType = "uint64"
	ValueTypeInt8    SubDataType = "int8"
	ValueTypeInt16   SubDataType = "int16"
	ValueTypeInt32   SubDataType = "int32"
	ValueTypeInt64   SubDataType = "int64"
	ValueTypeFloat32 SubDataType = "float32"
	ValueTypeFloat64 SubDataType = "float64"
)

type DataType string

const (
	ValueType  DataType = "value" // int64
	BoolType   DataType = "bool"
	StringType DataType = "string"
	RawType    DataType = "raw"   // []byte
	EnumType   DataType = "enum"  // string
	FaultType  DataType = "fault" // []string
	FlotType   DataType = "float"
	DoubleType DataType = "double"
	DateType   DataType = "date"
	BitmapType DataType = "bitmap"
	StructType DataType = "struct"
	ArrayType  DataType = "array"
)

//////////////////////////////////////////////////////////////////////////////
type RWType string

const (
	ReadOnly  RWType = "ro" // 只上报
	WriteOnly RWType = "wo" // 只下发
	ReadWrite RWType = "rw" // 可上报可下发
)

//////////////////////////////////////////////////////////////////////////////
// device
const (
	// 设备激活状态 未激活、已激活、激活失败
	DeviceActiveStatusInactivated = "inactivated"
	DeviceActiveStatusActivated   = "activated"
	DeviceActiveStatusActiveFail  = "activeFail"
)

type DeviceOnlineStatus string

const (
	Online  DeviceOnlineStatus = "online"
	Offline DeviceOnlineStatus = "offline"
)

type DeviceNotifyType string

const (
	DeviceAddNotify    DeviceNotifyType = "add"
	DeviceUpdateNotify DeviceNotifyType = "update"
	DeviceDeleteNotify DeviceNotifyType = "delete"
	DeviceActiveNotify DeviceNotifyType = "active"
)

//////////////////////////////////////////////////////////////////////////////
type ProductNotifyType string

const (
	ProductAddNotify    ProductNotifyType = "add"
	ProductUpdateNotify ProductNotifyType = "update"
	ProductDeleteNotify ProductNotifyType = "delete"
)

//////////////////////////////////////////////////////////////////////////////
type AlertLevel int

const (
	ERROR AlertLevel = iota + 1
	WARN
	NOTIFY
)

const (
	DEVICE_REPORT_EVENT = "report"
	DEVICE_ALERT_EVENT  = "alert"
)

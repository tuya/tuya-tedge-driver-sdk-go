package dpmodel

import "github.com/tuya/tuya-tedge-driver-sdk-go/commons"

/////////////////////////////////////////////////////////////////////////////////////////////
//下发DP消息：涂鸦云-->TEdge-->驱动
// CommandRequest 云端下发的DP消息
type CommandRequest struct {
	Protocol int32 // mqtt协议号
	T        int64
	S        int64
	Data     map[string]interface{}
}

type (
	DPExtend struct {
		Property PropertyValue
		Attr     map[string]string
	}
	DPExtendInfo map[string]DPExtend
)

/////////////////////////////////////////////////////////////////////////////////////////////
//上报DP消息：驱动-->TEdge-->涂鸦云
type (
	// WithoutWithDPValue 设备不带dp点数据上报
	WithoutDPValue struct {
		Protocol int32
		S        int64
		T        int64
		Data     interface{}
	}

	// WithDPValue 设备带dp点数据上报
	WithDPValue struct {
		DPId   string
		DPType commons.DataType
		Value  interface{}
	}
)

func NewWithDPValue(dpId string, dpType commons.DataType, value interface{}) *WithDPValue {
	return &WithDPValue{
		DPId:   dpId,
		DPType: dpType,
		Value:  value,
	}
}

func NewWithoutDPValue(protocol int32, s, t int64, data interface{}) *WithoutDPValue {
	return &WithoutDPValue{
		Protocol: protocol,
		S:        s,
		T:        t,
		Data:     data,
	}
}

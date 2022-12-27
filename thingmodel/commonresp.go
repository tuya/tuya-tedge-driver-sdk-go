package thingmodel

const Version = "1.0"

type ACK struct {
	// 默认情况下设备属性上报后云端不会返回应答消息，
	// 但可以通过ack参数改变这一默认行为：
	// 0:不做应答(默认)
	// :处理之后返回应答消息
	Ack int8 `json:"ack"`
}

// CommonResponse 云端一般响应
type CommonResponse struct {
	Version string `json:"version"`
	MsgId   string `json:"msgId"`
	Time    int64  `json:"time"`
	Code    int64  `json:"code"` // 0代表成功，非0代表失败，默认0
}

type CommonRequest struct {
	Version string `json:"version"` // 默认仅有1.0
	MsgId   string `json:"msgId"`   // 长度不超过32位
	Time    int64  `json:"time"`    // 消息发送时的unix时间戳（10位秒级或13位毫秒级）
	Sys     ACK    `json:"sys"`     // 控制消息的系统行为
}

func NewCommonResponse(t, code int64, msgId string) CommonResponse {
	return CommonResponse{
		Version: Version,
		MsgId:   msgId,
		Time:    t,
		Code:    code,
	}
}

type HttpRequestParam struct {
	Url     string `json:"url"`
	Api     string `json:"api"`
	Version string `json:"version"`
}

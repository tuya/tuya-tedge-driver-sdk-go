package common

type AtopResp struct {
	T         int64  `json:"t"`
	Success   bool   `json:"success"`
	ErrorCode string `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

const (
	ATOP_ILLEGAL_ERR  = "ILLEGAL_ACCESS_API"
	ATOP_UNKNOW_ERR   = "REMOTE_API_RUN_UNKNOW_FAILED"
	ATOP_REPEATED_ERR = "REPEATED_REQUEST"
	ATOP_LIMIT_ERR    = "API_QPS_LIMIT_OR_DEGRADE"
	ATOP_INTERNAL_ERR = "INTERNAL_ERROR"
)

///////////////////////////////////////////////////////////////////////////////////////////
type IPCSkillSet struct {
	CloudGW    int32  `json:"cloudGW"` //支持推流网关
	WebRtc     int32  `json:"webrtc"`
	P2P        int32  `json:"p2p"`
	SdkVersion string `json:"sdk_version"`
}

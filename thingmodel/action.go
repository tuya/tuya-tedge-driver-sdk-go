package thingmodel

import "time"

//物模型：动作 Action
type (
	// 云端-->TEdge-->驱动：设备执行动作入参
	ActionDataIn struct {
		ActionCode  string                 `json:"actionCode"`
		InputParams map[string]interface{} `json:"inputParams"`
	}

	// ActionExecuteRequest 执行设备动作请求
	ActionExecuteRequest struct {
		CommonRequest `json:",inline"`
		Data          ActionDataIn `json:"data"`
		Spec          Action       `json:"-"`
	}

	// 驱动-->云端-->TEdge：设备执行动作结果出参
	ActionDataOut struct {
		ActionCode   string                 `json:"actionCode"`
		OutputParams map[string]interface{} `json:"outputParams"`
	}

	// ActionExecuteResponse 执行设备动作响应
	ActionExecuteResponse struct {
		CommonResponse `json:",inline"`
		Data           ActionDataOut `json:"data"`
	}
)

//statusCode: 响应状态码，0代表成功，非0代表失败，默认0
func NewActionExecuteResponse(statusCode int64, msgId string, data ActionDataOut) ActionExecuteResponse {
	return ActionExecuteResponse{
		CommonResponse: NewCommonResponse(time.Now().Unix(), statusCode, msgId),
		Data: data,
	}
}

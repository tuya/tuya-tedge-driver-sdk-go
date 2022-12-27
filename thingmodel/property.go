package thingmodel

import "time"

type (
	Property struct {
		Value interface{} `json:"value"` // 上报的属性值
		Time  int64       `json:"time"`  // 属性变更时间戳
	}
	// PropertyReport 属性上报 属性查询响应
	PropertyReport struct {
		CommonRequest `json:",inline"`
		Data          map[string]Property `json:"data"`
	}

	// 云端-->TEdge-->驱动：设备属性下发
	PropertySet struct {
		CommonRequest `json:",inline"`
		Data          map[string]interface{}  `json:"data"`
		Spec          map[string]PropertySpec `json:"-"`
	}

	// PropertyGet 属性查询,查询的属性code列表，如果为空代表查询所有属性
	PropertyGet struct {
		CommonRequest `json:",inline"`
		Data          []string                `json:"data"`
		Spec          map[string]PropertySpec `json:"-"`
	}

	// PropertyGetResponse 属性查询设备响应
	PropertyGetResponse struct {
		CommonResponse `json:",inline"`
		Data           map[string]Property `json:"data"`
	}
	// PropertyDesiredGet 设备拉取属性期望值
	PropertyDesiredGet struct {
		CommonRequest `json:",inline"`
		//Data          []string `json:"data"`
		PropertyData  `json:"data"`
	}

	PropertyData struct {
		Properties []string `json:"properties"`
	}

	// PropertyDesiredGetResponse 设备拉取属性期望值响应
	PropertyDesiredGetResponse struct {
		CommonResponse `json:",inline"`
		Data           map[string]PropertyDesiredGetValue `json:"data"`
	}

	// PropertyDesiredDelete 设备清除属性期望值
	PropertyDesiredDelete struct {
		CommonRequest `json:",inline"`
		Data          map[string]PropertyDesiredDeleteValue `json:"data"`
	}

	// PropertyDesiredDeleteResponse 设备清除属性期望值响应
	PropertyDesiredDeleteResponse struct {
		CommonResponse `json:",inline"`
		Data           map[string]PropertyDesiredGetValue `json:"data"`
	}

	PropertyDesiredGetValue struct {
		Value   interface{} `json:"value"`
		Version int64       `json:"version"`
	}

	PropertyDesiredDeleteValue struct {
		Version int64 `json:"version"`
	}
)

func NewProperty(value interface{}) Property {
	return Property{
		Value: value,
		Time:  time.Now().Unix(),
	}
}

func NewPropertyReport(needACK bool, data map[string]Property) PropertyReport {
	var ack int8
	if needACK {
		ack = 1
	}

	return PropertyReport{
		CommonRequest: CommonRequest{
			Version: Version,
			Time:    time.Now().Unix(),
			Sys: ACK{
				Ack: ack,
			},
		},
		Data: data,
	}
}

//statusCode: 响应状态码，0 代表成功，非 0 代表失败，默认 0
func NewPropertyGetResponse(statusCode int64, msgId string, data map[string]Property) PropertyGetResponse {
	return PropertyGetResponse{
		CommonResponse: CommonResponse{
			Version: Version,
			MsgId:   msgId,
			Time:    time.Now().Unix(),
			Code:    statusCode,
		},
		Data: data,
	}
}

func NewPropertyDesiredGet(data []string) PropertyDesiredGet {
	return PropertyDesiredGet{
		CommonRequest: CommonRequest{
			Version: Version,
			Time:    time.Now().Unix(),
			Sys: ACK{
				Ack: 1,
			},
		},
		//Data: data,
		PropertyData: PropertyData {
			Properties: data,
		},
	}
}

func NewPropertyDesiredDelete(data map[string]PropertyDesiredDeleteValue) PropertyDesiredDelete {
	return PropertyDesiredDelete{
		CommonRequest: CommonRequest{
			Version: Version,
			Time:    time.Now().Unix(),
			Sys: ACK{
				Ack: 1,
			},
		},
		Data: data,
	}
}

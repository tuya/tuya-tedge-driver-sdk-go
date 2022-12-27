package thingmodel

import "time"

type (
	// 驱动-->云端-->TEdge：设备向云端上报事件
	EventReport struct {
		CommonRequest `json:",inline"`
		Data          EventData `json:"data"`
	}
	EventData struct {
		EventCode    string                 `json:"eventCode"`
		EventTime    int64                  `json:"eventTime"`
		OutputParams map[string]interface{} `json:"outputParams"`
	}
)

func NewEventData(ec string, t int64, data map[string]interface{}) EventData {
	return EventData{
		EventCode:    ec,
		EventTime:    t,
		OutputParams: data,
	}
}

func NewEventReport(needACK bool, data EventData) EventReport {
	var ack int8
	if needACK {
		ack = 1
	}
	return EventReport{
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

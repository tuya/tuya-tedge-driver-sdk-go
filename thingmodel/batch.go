package thingmodel

import "time"

// BatchReport 设备批量上报
type (
	// BatchReport 设备批量上报
	BatchReport struct {
		CommonRequest `json:",inline"`
		Data          BatchData `json:"data"`
	}
	BatchProperty struct {
		Value interface{} `json:"value"` // 上报的属性值
		Time  int64       `json:"time"`  // 属性变更时间戳
	}
	BatchEvent struct {
		EventTime    int64                  `json:"eventTime"`
		OutputParams map[string]interface{} `json:"outputParams"`
	}
	BatchData struct {
		Properties map[string]BatchProperty `json:"properties"`
		Events     map[string]BatchEvent    `json:"events"`
	}
)

func NewBatchReport(needACK bool, data BatchData) BatchReport {
	var ack int8
	if needACK {
		ack = 1
	}
	return BatchReport{
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

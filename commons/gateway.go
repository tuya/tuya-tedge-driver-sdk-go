package commons

type GatewayInfo struct {
	GwId        string `json:"gwId"`
	LocalKey    string `json:"localKey"`
	Env         string `json:"env"`
	Region      string `json:"region"`
	Mode        RunningModel
	CloudState  bool   `json:"cloudState"`
	GatewayName string `json:"gatewayName"`
}

type RunningModel string

const (
	// DPModel dp模型
	DPModel RunningModel = "dp_model"
	// ThingModel 物模型
	ThingModel RunningModel = "thing_model"
)

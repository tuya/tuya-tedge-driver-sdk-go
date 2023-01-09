package commons

import (
	"github.com/eclipse/paho.mqtt.golang"
)

type Option interface {
	Apply(*Options)
}

type Options struct {
	MqttUsername   string
	MqttDriver     MqttDriver
	ConnHandler    mqtt.OnConnectHandler
	AppService     []string
	AppDataHandler AppCallBack
	EnableRTS      bool
}

func DefaultOptions() Options {
	return Options{
		MqttDriver:     nil,
		ConnHandler:    nil,
		AppService:     nil,
		MqttUsername:   "",
		AppDataHandler: nil,
		EnableRTS:      false,
	}
}

type AppCallBack func(appName string, req AppDriverReq) (Response, error)


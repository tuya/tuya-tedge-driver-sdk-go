package driversvc

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/clients"
)

type Option interface {
	Apply(*Options)
}

type Options struct {
	MqttUsername   string
	MqttDriver     commons.MqttDriver
	ConnHandler    mqtt.OnConnectHandler
	AppService     []string
	AppDataHandler clients.AppCallBack
	EnableRTS      bool
}

func defaultOptions() Options {
	return Options{
		MqttDriver:     nil,
		ConnHandler:    nil,
		AppService:     nil,
		MqttUsername:   "",
		AppDataHandler: nil,
		EnableRTS:      false,
	}
}

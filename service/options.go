package service

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/clients"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/driversvc"
)

type funcOption struct {
	f func(service *driversvc.Options)
}

func (fdo *funcOption) Apply(do *driversvc.Options) {
	fdo.f(do)
}

func newFuncOption(f func(*driversvc.Options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithMqtt(driver commons.MqttDriver, username string, connHandler mqtt.OnConnectHandler) driversvc.Option {
	return newFuncOption(func(o *driversvc.Options) {
		o.MqttDriver = driver
		o.ConnHandler = connHandler
		o.MqttUsername = username
	})
}

func WithAppService(appServiceNames []string) driversvc.Option {
	return newFuncOption(func(o *driversvc.Options) {
		o.AppService = appServiceNames
	})
}

func WithCallBack(callBack clients.AppCallBack) driversvc.Option {
	return newFuncOption(func(o *driversvc.Options) {
		o.AppDataHandler = callBack
	})
}

func WithRtsOption() driversvc.Option {
	return newFuncOption(func(o *driversvc.Options) {
		o.EnableRTS = true
	})
}

package service

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
)

type funcOption struct {
	f func(service *commons.Options)
}

func (fdo *funcOption) Apply(do *commons.Options) {
	fdo.f(do)
}

func newFuncOption(f func(*commons.Options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithMqtt(driver commons.MqttDriver, username string, connHandler mqtt.OnConnectHandler) commons.Option {
	return newFuncOption(func(o *commons.Options) {
		o.MqttDriver = driver
		o.ConnHandler = connHandler
		o.MqttUsername = username
	})
}

func WithAppService(appServiceNames []string) commons.Option {
	return newFuncOption(func(o *commons.Options) {
		o.AppService = appServiceNames
	})
}

func WithCallBack(callBack commons.AppCallBack) commons.Option {
	return newFuncOption(func(o *commons.Options) {
		o.AppDataHandler = callBack
	})
}

func WithRtsOption() commons.Option {
	return newFuncOption(func(o *commons.Options) {
		o.EnableRTS = true
	})
}

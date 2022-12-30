[English](README.md) | [中文版](README_CN.md)
# `tuya-edge-driver-sdk-go`: TEdge Southbound Driver SDK

## Terms
* TEdge: the Tuya IoT Edge Gateway, used to connect third-party devices to Tuya IoT Cloud.
* tedge-driver-sdk-go: the driver SDK, used to connect devices to TEdge through southbound interfaces.
* Driver: the southbound plug-in of TEdge, used to interface with third-party devices.

## Architecture of TEdge
![Image](./docs/images/Tedge架构图1.png)

## Get started

### Procedure
1. Follow the example and implement the driver interface `DPModelDriver`.
2. Package the driver into a Docker container.

- For more information about the sample, see [TEdge Driver Demo](https://github.com/tuya/tuya-tedge-driver-example).
- For more information about the guidelines on driver development, see [Documents on driver development](./docs/summary.md).

### DPModelDriver interface
```golang
type DPModelDriver interface {
	// The callback to invoke when a sub-device is added, activated, updated, or deleted in the TEdge console.
	DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, cid string, device commons.DeviceInfo) error

	// The callback to invoke when a product is added, updated, or deleted in the TEdge console.
	ProductNotify(ctx context.Context, t commons.ProductNotifyType, pid string, product DPModelProduct) error

	// The callback to invoke when a command is received and forwarded over MQTT in the following direction: Tuya IoT Cloud > TEdge > Sub-device
	HandleCommands(ctx context.Context, cid string, req CommandRequest, protocols map[string]commons.ProtocolProperties, dpExtend DPExtendInfo) error
    
        // The callback to invoke when a driver instance is stopped from running in the TEdge console. The driver program can be recycled.
	Stop(ctx context.Context) error
}
```

### Example
```golang
package main

import (
	"context"
	"fmt"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
	"github.com/tuya/tuya-tedge-driver-sdk-go/service"
)

// Example of TEdge driver development
func main() {
	sdkLog := commons.DefaultLogger(commons.DebugLevel, "driver-example")
	dpService := service.NewDPService(sdkLog)

	// DP model driver
	dpDriver := NewDemoDPDriver(dpService)
	go dpDriver.Run()

	// Note: `dpDriver` requires the implementation of the interface `type DPModelDriver interface`.
	//Start: blocked
	err := dpService.Start(dpDriver)
	if err != nil {
		sdkLog.Errorf("tedge-driver start err:%s", err)
		panic(fmt.Sprintf("tedge-driver start err:%s", err))
	}
}

// `DemoDpDriver` requires the implementation of the interface `type DPModelDriver interface`.
// The interface is defined in the SDK `tedge-driver-sdk-go/dpmodel/interface.go`.
type DemoDPDriver struct {
	dpService *service.DPDriverService
}

func NewDemoDPDriver(dpService *service.DPDriverService) *DemoDPDriver {
	return &DemoDPDriver{
		dpService: dpService,
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// Implementation of the interface `DPModelDriver`.
// 1. Receive messages sent by TEdge or Tuya IoT Cloud over MQTT.
// 2. Note: Do not perform blocking operations with the interface.
func (dd *DemoDPDriver) HandleCommands(ctx context.Context, cid string, msg dpmodel.CommandRequest, protocols map[string]commons.ProtocolProperties, dpExtend dpmodel.DPExtendInfo) error {
	//......
	//TODO: implement me

	return nil
}

// Implementation of the interface `DPModelDriver`.
// 1. The callback to invoke when a sub-device is added, activated, updated, or deleted in the TEdge console.
// 2. Note: Do not perform blocking operations with the interface.
// 3. Set the implementation of the interface to nil if you do not need to manually add sub-devices to the target gateway in the TEdge console.
func (dd *DemoDPDriver) DeviceNotify(ctx context.Context, action commons.DeviceNotifyType, cid string, device commons.DeviceInfo) error {
	//......
	//TODO: implement me
        //Send message to the real device

	return nil
}

// Implementation of the interface `DPModelDriver`.
// 1. ProductNotify: the callback to invoke when a product is added, updated, or deleted.
// 2. Note: Do not perform blocking operations with the interface.
func (dd *DemoDPDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, pid string, product dpmodel.DPModelProduct) error {
	return nil
}

// Implementation of the interface `DPModelDriver`.
// The callback to invoke when a driver instance is updated or stopped from running on TEdge. The driver program can be recycled.
func (dd *DemoDPDriver) Stop(ctx context.Context) error {
	return nil
}

// The callback to invoke when a driver instance is updated or stopped from running on TEdge. The driver program can be recycled.
func (dd *DemoDPDriver) Run() {
	//......
}
```

### Implementation driver with MQTT

```golang
type MqttDriver struct {
	logger commons.TedgeLogger
}

func NewMqttDriver(l commons.TedgeLogger) *MqttDriver {
	return &MqttDriver{
		logger: l,
	}
}

var _ commons.MqttDriver = (*MqttDriver)(nil)

// Auth: the MQTT authentication event. A value of `true` indicates successful authentication.
func (md *MqttDriver) Auth(clientId, username, password string) (bool, error) {
	md.logger.Debugf("auth: clientId: %s, username: %s, password: %s", clientId, username, password)
	return true, nil
}

// Sub: the MQTT subscription event. A value of `true` indicates a successful subscription.
func (md *MqttDriver) Sub(clientId, username, topic string, qos byte) (bool, error) {
	md.logger.Debugf("sub: clientId: %s, username: %s, topic: %s, qos: %d", clientId, username, topic, qos)
	return true, nil
}

// Pub: the MQTT authentication event. A value of `true` indicates successful publishing.
func (md *MqttDriver) Pub(clientId, username, topic string, qos byte, retained bool) (bool, error) {
	md.logger.Debugf("pub: clientId: %s, username: %s, topic: %s, qos: %d, retained: %t", clientId, username, topic, qos, retained)
	return true, nil
}

// UnSub: unsubscribe from an MQTT topic
func (md *MqttDriver) UnSub(clientId, username string, topics []string) {
	md.logger.Debugf("unsub: clientId: %s, username: %s, topics: %+v", clientId, username, topics)
}

// Connected
func (md *MqttDriver) Connected(clientId, username string) {
}

// Closed
func (md *MqttDriver) Closed(clientId, username string) {
	md.logger.Debugf("clised: clientId: %s, username: %s", clientId, username)
}

func (md *MqttDriver) OnConnectedHandler() mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		topic := "tuya/tedge/custom/test1"
		md.logger.Debugf("sub topic(%s) ...", topic)
		if token := client.Subscribe(topic, byte(1), md.OnMessageReceived()); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		md.logger.Debugf("sub topic(%s) success", topic)
	}
}

func (md *MqttDriver) OnMessageReceived() mqtt.MessageHandler {
	return func(client mqtt.Client, message mqtt.Message) {
		md.logger.Debugf("Received message on topic: %s, Message: %s", message.Topic(), message.Payload())
		client.Publish("tuya/tedge/custom/test2", byte(1), false, message.Payload())
	}
}
```

## SDK API
* DP model SDK API: `tedge-driver-sdk-go/service/dpmodelapi.go`
* TuyaLink model SDK API: `tedge-driver-sdk-go/service/tymodelapi.go`

## Technical support
Tuya IoT Developer Platform: https://developer.tuya.com/en/

Tuya Developer Help Center: https://support.tuya.com/en/help

Tuya Service Ticket System: https://service.console.tuya.com/


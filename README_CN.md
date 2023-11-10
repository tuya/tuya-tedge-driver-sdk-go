[English](README.md) | [中文版](README_CN.md)
# tuya-tedge-driver-sdk-go：涂鸦边缘网关 Tedge 南向驱动开发 SDK

## 名词简介
* Tedge: 涂鸦边缘网关，主要用于将第三方设备接入涂鸦云。
* tedge-driver-sdk-go: 驱动开发SDK，连接南向设备和Tedge。
* 驱动程序：Tedge南向插件，用来对接第三方设备。

## Tedge 架构
![Tedge架构图.png](./docs/images/Tedge架构图1.png)

## 快速开始

### 驱动开发步骤
1. 参考"驱动开发示例"，实现驱动接口 `DPModelDriver`
2. 将驱动打包成一个 docker 容器即可
3. 完整的示例请参考：[驱动程序Demo](https://github.com/tuya/tuya-tedge-driver-example)
4. 完整驱动开发指南请阅读：[驱动开发指南](./docs/summary.md)

### DPModelDriver interface
```golang
type DPModelDriver interface {
	// 在Tedge Web新增/激活/更新/删除一个子设备时，回调该接口
	DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, cid string, device commons.DeviceInfo) error

	// 在Tedge Web新增/更新/删除一个产品时，回调该接口
	ProductNotify(ctx context.Context, t commons.ProductNotifyType, pid string, product DPModelProduct) error

	// Tedge收到云端发往子设备的指令时(mqtt 消息)，回调该接口：tuya cloud-->Tedge-->device
	HandleCommands(ctx context.Context, cid string, req CommandRequest, protocols map[string]commons.ProtocolProperties, dpExtend DPExtendInfo) error
    
    // 在Tedge Web停止驱动实例运行时，回调该接口，驱动程序可以进行资源回收
	Stop(ctx context.Context) error
}
```

### 驱动开发示例
```golang
package main

import (
	"context"
	"fmt"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
	"github.com/tuya/tuya-tedge-driver-sdk-go/service"
)

//TEdge 驱动开发示例
func main() {
	sdkLog := commons.DefaultLogger(commons.DebugLevel, "driver-example")
	dpService := service.NewDPService(sdkLog)

	// DP 模型驱动
	dpDriver := NewDemoDPDriver(dpService)
	go dpDriver.Run()

	//注：dpDriver 必须实现接口 `type DPModelDriver interface`
	//Start: blocked
	err := dpService.Start(dpDriver)
	if err != nil {
		sdkLog.Errorf("tedge-driver start err:%s", err)
		panic(fmt.Sprintf("tedge-driver start err:%s", err))
	}
}

// DemoDpDriver 必须实现接口 `type DPModelDriver interface`
// 接口定义在sdk：`tedge-driver-sdk-go/dpmodel/interface.go`
type DemoDPDriver struct {
	dpService *service.DPDriverService
}

func NewDemoDPDriver(dpService *service.DPDriverService) *DemoDPDriver {
	return &DemoDPDriver{
		dpService: dpService,
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
//DPModelDriver 接口实现
//1.接收 "Tedge或涂鸦云端" 下发的MQTT消息
//2.注意：不要在该接口做阻塞性操作
func (dd *DemoDPDriver) HandleCommands(ctx context.Context, cid string, msg dpmodel.CommandRequest, protocols map[string]commons.ProtocolProperties, dpExtend dpmodel.DPExtendInfo) error {
	//......
	//TODO: implement me

	return nil
}

//DPModelDriver 接口实现
//1.在Tedge控制台页面，新增、激活、更新子设备属性、删除子设备时，回调该接口
//2.注意：不要在该接口做阻塞性操作
//3.如果接入的设备不需要在Tedge控制台页面手动新增子设备，则该接口实现为空即可
func (dd *DemoDPDriver) DeviceNotify(ctx context.Context, action commons.DeviceNotifyType, cid string, device commons.DeviceInfo) error {
	//......
	//TODO: implement me

	return nil
}

//DPModelDriver 接口实现
//1.ProductNotify 产品增加/删除/更新通知回调
//2.注意：不要在该接口做阻塞性操作
func (dd *DemoDPDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, pid string, product dpmodel.DPModelProduct) error {
	return nil
}

//DPModelDriver 接口实现
//在Tedge页面，更新驱动实例，或停止驱动实例时，回调该接口，驱动程序进行资源清理
func (dd *DemoDPDriver) Stop(ctx context.Context) error {
	return nil
}

//运行自定义逻辑
func (dd *DemoDPDriver) Run() {
	//......
}

```

### Implementation driver with mqtt

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

// Auth mqtt 鉴权事件，鉴权成功返回true
func (md *MqttDriver) Auth(clientId, username, password string) (bool, error) {
	md.logger.Debugf("auth: clientId: %s, username: %s, password: %s", clientId, username, password)
	return true, nil
}

// Sub mqtt subscribe 订阅事件，鉴权成功返回true
func (md *MqttDriver) Sub(clientId, username, topic string, qos byte) (bool, error) {
	md.logger.Debugf("sub: clientId: %s, username: %s, topic: %s, qos: %d", clientId, username, topic, qos)
	return true, nil
}

// Pub mqtt publish 事件，鉴权成功返回true
func (md *MqttDriver) Pub(clientId, username, topic string, qos byte, retained bool) (bool, error) {
	md.logger.Debugf("pub: clientId: %s, username: %s, topic: %s, qos: %d, retained: %t", clientId, username, topic, qos, retained)
	return true, nil
}

// UnSub mqtt unsubscribe 取消订阅事件
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
* DP模型SDK API: `tedge-driver-sdk-go/service/dpmodelapi.go`
* TuyaLink模型SDK API: `tedge-driver-sdk-go/service/tymodelapi.go`

## 技术支持
Tuya IoT Developer Platform: https://developer.tuya.com/en/

Tuya Developer Help Center: https://support.tuya.com/en/help

Tuya Work Order System: https://service.console.tuya.com/
package thingmodel

import (
	"context"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
)

type ThingModelDriver interface {
	// DeviceNotify 设备增删改通知,删除设备时device参数为空
	DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, cid string, device commons.TMDeviceInfo) error

	// ProductNotify 产品增删改通知,删除产品时product参数为空
	ProductNotify(ctx context.Context, t commons.ProductNotifyType, pid string, product ThingModelProduct) error


	// HandlePropertySet set device property, from cloud or web
	HandlePropertySet(ctx context.Context, cid string, data PropertySet, protocols map[string]commons.ProtocolProperties) error
	// HandlePropertyGet get device property, from cloud
	HandlePropertyGet(ctx context.Context, cid string, data PropertyGet, protocols map[string]commons.ProtocolProperties) error
	// HandleActionExecute device command, from cloud or web
	HandleActionExecute(ctx context.Context, cid string, data ActionExecuteRequest, protocols map[string]commons.ProtocolProperties) error

	Stop(ctx context.Context) error
}

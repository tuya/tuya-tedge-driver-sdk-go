package dpmodel

import (
	"context"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
)

type DPModelDriver interface {
	// DeviceNotify 设备新增/更新/删除回调,删除设备时device参数为空
	DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, cid string, device commons.DeviceInfo) error

	// ProductNotify 产品新增/更新/删除回调,删除产品时product参数为空
	ProductNotify(ctx context.Context, t commons.ProductNotifyType, pid string, product DPModelProduct) error

	// HandleCommands MQTT消息：云端-->Tedge-->设备
	HandleCommands(ctx context.Context, cid string, req CommandRequest, protocols map[string]commons.ProtocolProperties, dpExtend DPExtendInfo) error

	Stop(ctx context.Context) error
}

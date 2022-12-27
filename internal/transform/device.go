package transform

import (
	"encoding/json"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
)

// TYDeviceExt ty device extend data
type TYDeviceExt struct {
	LocalKey string `json:"localKey"`
	Mac      string `json:"mac"`
}

func ToDeviceModel(dev *proto.DeviceInfo) commons.DeviceInfo {
	var d commons.DeviceInfo
	d.Cid = dev.Id
	//d.Description = dev.GetDescription()
	d.Protocols = ToProtocolModels(dev.GetProtocols()) //TODO
	d.ProductId = dev.ProductId

	d.BaseAttr.Name = dev.Name
	d.BaseAttr.Ip = dev.Ip
	d.BaseAttr.Lat = dev.Lat
	d.BaseAttr.Lon = dev.Lon

	//d.ExtendedAttr.VendorCode = dev.GetVendorCode()
	d.ExtendedAttr.InstallLocation = dev.GetInstallLocation()

	d.OnLineStatus = dev.GetOnlineStatus()
	d.ActiveStatus = dev.GetActiveStatus()
	d.CloudDeviceId = dev.GetCloudDeviceId()

	// ty device
	if dev.GetIsScreenDevice() {
		d.ActiveStatus = commons.DeviceActiveStatusActivated
		d.ExtendedAttr.ExtendData = map[string]interface{}{
			"parentCloudId": dev.GetParentCloudId(),
			"screenType":    dev.GetScreenType(),
		}
		if len(dev.GetExtendData()) > 0 {
			var ext TYDeviceExt
			if err := json.Unmarshal([]byte(dev.GetExtendData()), &ext); err == nil {
				d.ExtendedAttr.ExtendData["localKey"] = ext.LocalKey
				d.ExtendedAttr.ExtendData["mac"] = ext.Mac
			}
		}
	} else {
		if len(dev.GetExtendData()) > 0 {
			var extendDataByte map[string]interface{}
			if err := json.Unmarshal([]byte(dev.GetExtendData()), &extendDataByte); err == nil {
				d.ExtendedAttr.ExtendData = extendDataByte
			}
		}
	}

	return d
}

func ToProtocolModels(protocol map[string]*proto.ProtocolProperties) map[string]commons.ProtocolProperties {
	protocolModels := make(map[string]commons.ProtocolProperties)
	for k, protocolProperties := range protocol {
		protocolModels[k] = ToProtocolPropertiesModel(protocolProperties)
	}
	return protocolModels
}

func ToProtocolPropertiesModel(p *proto.ProtocolProperties) map[string]interface{} {
	protocolProperties := make(commons.ProtocolProperties)
	for k, v := range p.Pp {
		protocolProperties[k] = v
	}
	return protocolProperties
}

func UpdateDeviceModelFieldsFromProto(dev *commons.DeviceInfo, patch *proto.DeviceUpdateInfo) {
	if patch.ProductId != nil {
		dev.ProductId = *patch.ProductId
	}

	if patch.Name != nil {
		dev.BaseAttr.Name = *patch.Name
	}

	if patch.Ip != nil {
		dev.BaseAttr.Ip = *patch.Ip
	}
	if patch.Lat != nil {
		dev.BaseAttr.Lat = *patch.Lat
	}
	if patch.Lon != nil {
		dev.BaseAttr.Lon = *patch.Lon
	}

	//if patch.VendorCode != nil {
	//	dev.ExtendedAttr.VendorCode = *patch.VendorCode
	//}

	if patch.InstallLocation != nil {
		dev.ExtendedAttr.InstallLocation = *patch.InstallLocation
	}

	dev.Protocols = ToProtocolModels(patch.Protocols)
	dev.ActiveStatus = *patch.ActiveStatus
	dev.OnLineStatus = *patch.OnlineStatus
	dev.CloudDeviceId = *patch.CloudDeviceId

	if patch.ExtendData != nil && len(*patch.ExtendData) > 0 {
		if *patch.IsScreenDevice {
			dev.ActiveStatus = commons.DeviceActiveStatusActivated
			dev.ExtendedAttr.ExtendData = map[string]interface{}{
				"parentCloudId": patch.GetParentCloudId(),
				"screenType":    patch.GetScreenType(),
			}
			var ext TYDeviceExt
			if err := json.Unmarshal([]byte(patch.GetExtendData()), &ext); err == nil {
				dev.ExtendedAttr.ExtendData["localKey"] = ext.LocalKey
				dev.ExtendedAttr.ExtendData["mac"] = ext.Mac
			}
		} else {
			var extendDataMap map[string]interface{}
			err := json.Unmarshal([]byte(*patch.ExtendData), &extendDataMap)
			if err == nil {
				dev.ExtendedAttr.ExtendData = extendDataMap
			}
		}
	}
}

func NewAddDeviceRequest(dev *proto.DeviceAddInfo) *proto.AddDeviceRequest {
	return &proto.AddDeviceRequest{
		Device: dev,
	}
}

func FromDeviceActiveInfoModelToProto(d commons.DeviceMeta, serviceId string, isIPC bool) (*proto.DeviceAddInfo, error) {
	var device proto.DeviceAddInfo

	device.Id = d.Cid
	device.ProductId = d.ProductId

	//基础属性
	device.Name = d.BaseAttr.Name
	device.Ip = d.BaseAttr.Ip
	device.Lat = d.BaseAttr.Lat
	device.Lon = d.BaseAttr.Lon

	//附加属性
	//device.VendorCode = d.ExtendedAttr.VendorCode

	device.InstallLocation = d.ExtendedAttr.InstallLocation
	if d.ExtendedAttr.ExtendData != nil {
		extendDataByte, err := json.Marshal(d.ExtendedAttr.ExtendData)
		if err == nil {
			extendDataStr := string(extendDataByte)
			device.ExtendData = extendDataStr
		}
	}

	device.IsIpcDev = isIPC //是否是IPC类子设备
	device.ServiceId = serviceId

	device.Source = common.DeviceFromDriver // 标识驱动添加

	return &device, nil
}

func GenIpcSkillSet() (string, error) {
	ipcSkillSet := common.IPCSkillSet{
		WebRtc:     0x7f,     //固定值，能力集有点乱，全部设置
		P2P:        3,        //固定值
		SdkVersion: "e1.0.1", //Tedge专用版本号，P2P定位归类使用
		CloudGW:    1,        //支持推流网关
	}

	ipcSkillSetStr, err := json.Marshal(ipcSkillSet)
	if err != nil {
		return "", err
	}

	return string(ipcSkillSetStr), nil
}

func DeviceAddToModel(dev *proto.DeviceAddInfo) commons.DeviceInfo {
	var d commons.DeviceInfo
	d.Cid = dev.GetId()
	d.Protocols = ToProtocolModels(dev.GetProtocols()) //TODO
	d.ProductId = dev.GetProductId()

	d.BaseAttr.Name = dev.GetName()
	d.BaseAttr.Ip = dev.GetIp()
	d.BaseAttr.Lat = dev.GetLat()
	d.BaseAttr.Lon = dev.GetLon()

	//d.ExtendedAttr.VendorCode = dev.GetVendorCode()
	d.ExtendedAttr.InstallLocation = dev.GetInstallLocation()

	d.OnLineStatus = dev.GetOnlineStatus()
	d.ActiveStatus = dev.GetActiveStatus()
	d.CloudDeviceId = dev.GetCloudDeviceId()

	// ty device
	if dev.GetIsScreenDevice() {
		d.ActiveStatus = commons.DeviceActiveStatusActivated
		d.ExtendedAttr.ExtendData = map[string]interface{}{
			"parentCloudId": dev.GetParentCloudId(),
			"screenType":    dev.GetScreenType(),
		}
		if len(dev.GetExtendData()) > 0 {
			var ext TYDeviceExt
			if err := json.Unmarshal([]byte(dev.GetExtendData()), &ext); err == nil {
				d.ExtendedAttr.ExtendData["localKey"] = ext.LocalKey
				d.ExtendedAttr.ExtendData["mac"] = ext.Mac
			}
		}
	} else {
		if len(dev.GetExtendData()) > 0 {
			var extendDataByte map[string]interface{}
			err := json.Unmarshal([]byte(dev.GetExtendData()), &extendDataByte)
			if err == nil {
				d.ExtendedAttr.ExtendData = extendDataByte
			}
		}
	}

	return d
}

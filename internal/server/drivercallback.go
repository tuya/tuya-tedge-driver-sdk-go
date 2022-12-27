package server

import (
	"context"
	"fmt"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/interfaces"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/transform"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func DeviceInfoToTM(dev commons.DeviceInfo) commons.TMDeviceInfo {
	return commons.TMDeviceInfo{
		TMDeviceMeta: transform.FromDeviceToTMDevice(dev.DeviceMeta),
		DeviceId:     dev.CloudDeviceId,
		ActiveStatus: dev.ActiveStatus,
		OnLineStatus: dev.OnLineStatus,
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func AddDeviceCallback(
	ctx context.Context,
	in *proto.DeviceAddInfo,
	serviceItf interface{},
	tEdgeModel commons.RunningModel) (*emptypb.Empty, error) {

	cid := in.GetId()
	if len(cid) == 0 {
		return common.EmptyPb, fmt.Errorf("cid is null")
	}

	commonItf := serviceItf.(interfaces.DriverCommonItf)
	lc := commonItf.GetLogger()
	dev := transform.DeviceAddToModel(in)
	lc.Debugf("AddDeviceCallback receive new device:%+v", dev)

	devCache := commonItf.GetDevCache()
	devCache.Add(dev)
	lc.Infof("AddDeviceCallback device %s added", dev.Cid)

	// 设备添加请求来源不是驱动才通知驱动，否则不通知
	source := int32(in.GetSource())
	if source != common.DeviceFromDriver {
		var err error
		if tEdgeModel == commons.DPModel {
			driverItf := serviceItf.(interfaces.DpDriverService)
			driver := driverItf.GetDriver()
			err = driver.DeviceNotify(ctx, commons.DeviceAddNotify, dev.Cid, dev)
		} else {
			driverItf := serviceItf.(interfaces.TyDriverService)
			driver := driverItf.GetDriver()
			err = driver.DeviceNotify(ctx, commons.DeviceAddNotify, dev.Cid, DeviceInfoToTM(dev))
		}

		if err != nil {
			lc.Errorf("AddDeviceCallback callback failed for %s, err:%v", dev.Cid, err)
			return common.EmptyPb, status.Errorf(codes.Internal, "callback failed for %s", dev.Cid)
		}
		lc.Debugf("AddDeviceCallback callback for %s", dev.Cid)
	}

	return common.EmptyPb, nil
}

func UpdateDeviceCallback(
	ctx context.Context,
	in *proto.DeviceUpdateInfo,
	serviceItf interface{},
	tEdgeModel commons.RunningModel) (*emptypb.Empty, error) {

	cid := in.GetId()
	if len(cid) == 0 {
		return common.EmptyPb, fmt.Errorf("cid is null")
	}

	commonItf := serviceItf.(interfaces.DpDriverService)
	lc := commonItf.GetLogger()

	notifyType := commons.DeviceUpdateNotify
	devCache := commonItf.GetDevCache()
	dev, ok := devCache.ById(cid)
	if !ok {
		notifyType = commons.DeviceAddNotify
		dev.Cid = cid
		transform.UpdateDeviceModelFieldsFromProto(&dev, in)
		devCache.Add(dev)
		lc.Warnf("UpdateDeviceCallback failed to find device cid:%s, add it", cid)
	} else {
		transform.UpdateDeviceModelFieldsFromProto(&dev, in)
		devCache.Update(dev)
		lc.Infof("UpdateDeviceCallback device updated: %+v", dev)
	}

	source := int32(in.GetSource())
	if source != common.DeviceFromDriver {
		var err error
		if tEdgeModel == commons.DPModel {
			driverItf := serviceItf.(interfaces.DpDriverService)
			driver := driverItf.GetDriver()
			err = driver.DeviceNotify(ctx, notifyType, cid, dev)
		} else {
			driverItf := serviceItf.(interfaces.TyDriverService)
			driver := driverItf.GetDriver()
			err = driver.DeviceNotify(ctx, notifyType, cid, DeviceInfoToTM(dev))
		}

		if err != nil {
			lc.Errorf("UpdateDeviceCallback callback failed cid:%s, err:%v", cid, err)
			return common.EmptyPb, status.Errorf(codes.Internal, "internal.UpdateDevice callback failed for %s", cid)
		}
		lc.Debugf("UpdateDeviceCallback callback cid:%s", cid)
	}

	return common.EmptyPb, nil
}

func DeleteDeviceCallback(
	ctx context.Context,
	in *proto.DeleteDeviceByIdRequest,
	serviceItf interface{},
	tEdgeModel commons.RunningModel) (*emptypb.Empty, error) {

	commonItf := serviceItf.(interfaces.DpDriverService)
	devCache := commonItf.GetDevCache()

	id := in.Id
	dev, ok := devCache.ById(id)
	if !ok {
		return common.EmptyPb, nil
	}

	devCache.RemoveById(id)
	lc := commonItf.GetLogger()
	lc.Infof("DeleteDeviceCallback device cid:%s", id)

	var err error
	if tEdgeModel == commons.DPModel {
		driverItf := serviceItf.(interfaces.DpDriverService)
		driver := driverItf.GetDriver()
		err = driver.DeviceNotify(ctx, commons.DeviceDeleteNotify, id, dev)
	} else {
		driverItf := serviceItf.(interfaces.TyDriverService)
		driver := driverItf.GetDriver()
		err = driver.DeviceNotify(ctx, commons.DeviceDeleteNotify, id, DeviceInfoToTM(dev))
	}

	if err != nil {
		lc.Errorf("DeleteDeviceCallback callback failed  cid:%s, err: %v", id, err)
		return common.EmptyPb, status.Errorf(codes.Internal, "internal.RemoveDevice callback failed cid:%s", dev.Cid)
	}

	return common.EmptyPb, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ChangeLogLevel(
	req *proto.LogLevelRequest,
	serviceItf interface{}) (*emptypb.Empty, error) {

	commonItf := serviceItf.(interfaces.DriverCommonItf)
	lc := commonItf.GetLogger()
	level := transform.ToLogLevel(req)
	if len(level) <= 0 {
		lc.Errorf("ChangeLogLevel level err, got:%d", req.GetLogLevel())
		return common.EmptyPb, status.Errorf(codes.InvalidArgument, "recv invalid log level: %d", req.GetLogLevel())
	}

	err := commonItf.ChangeLogLevel(level)
	return common.EmptyPb, err
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func AppServiceAddress(
	ctx context.Context,
	req *proto.AppBaseAddress,
	serviceItf interface{}) (*emptypb.Empty, error) {

	commonItf := serviceItf.(interfaces.DriverCommonItf)
	err := commonItf.UpdateAppAddress(ctx, req)
	return common.EmptyPb, err
}

//////////////////////////////////////////////////////////////////////////////////////////////
func GatewayStateCallback(
	req *proto.GatewayState,
	serviceItf interface{}) (*emptypb.Empty, error) {

	commonItf := serviceItf.(interfaces.DriverCommonItf)
	lc := commonItf.GetLogger()
	lc.Infof("GatewayStateCallback req:%+v", req)

	//set gateway status
	commonItf.SetCloudStatus(req.ConnStatus)
	return common.EmptyPb, nil
}

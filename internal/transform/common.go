package transform

import (
	"github.com/google/uuid"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
)

func ToAlertReportProto(serviceName string, t int64, level commons.AlertLevel, content string) *proto.AlertReportReq {
	return &proto.AlertReportReq{
		Id:          uuid.NewString(),
		Version:     "v1",
		ServiceName: serviceName,
		AlertType:   proto.AlertType_DRIVER_ALERT,
		AlertLevel:  *proto.AlertLevel(int32(level)).Enum(),
		T:           t,
		Content:     content,
	}
}

func ToGatewayModel(gwInfo *proto.GateWayInfoResponse) commons.GatewayInfo {
	if gwInfo == nil {
		return commons.GatewayInfo{}
	}
	var mode commons.RunningModel
	if gwInfo.GetIsNewModel() {
		mode = commons.ThingModel
	} else {
		mode = commons.DPModel
	}
	return commons.GatewayInfo{
		GwId:        gwInfo.GetGwId(),
		LocalKey:    gwInfo.GetLocalKey(),
		Env:         gwInfo.GetEnv(),
		Region:      gwInfo.GetRegion(),
		Mode:        mode,
		GatewayName: gwInfo.GatewayName,
		CloudState:  gwInfo.CloudState,
	}
}

func ToLogLevel(req *proto.LogLevelRequest) string {
	return logLevelMap[req.GetLogLevel()]
}

var logLevelMap = map[proto.EnumLogLevel]string{
	proto.EnumLogLevel_ENUM_LOG_LEVEL_DEBUG:   "DEBUG",
	proto.EnumLogLevel_ENUM_LOG_LEVEL_INFO:    "INFO",
	proto.EnumLogLevel_ENUM_LOG_LEVEL_WARNING: "WARNING",
	proto.EnumLogLevel_ENUM_LOG_LEVEL_ERROR:   "ERROR",
}

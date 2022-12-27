package server

import (
	"context"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/interfaces"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/rpcserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TyLink模型grpc.Server注册
// TEdge-->SDK
func TyServerRegister(rpcSrv *grpc.Server, tyService interfaces.TyDriverService) {
	NewCommonRPCServer().RegisterServer(rpcSrv)
	NewAppToDriverServer(tyService).RegisterServer(rpcSrv)
	NewTyCommandServer(tyService).RegisterServer(rpcSrv)
	NewTyDriverCommonServer(tyService).RegisterServer(rpcSrv)
}

/////////////////////////////////////////////////////////////////////////////////////////////
type TyDriverRpcServer struct {
	impServer interfaces.TyDriverService
}

func NewTyDriverRpcServer(implService interfaces.TyDriverService) *TyDriverRpcServer {
	return &TyDriverRpcServer{
		impServer: implService,
	}
}

func (tys *TyDriverRpcServer) Serve() error {
	lc := tys.impServer.GetLogger()
	regFunc := func(rserver *grpc.Server) {
		TyServerRegister(rserver, tys.impServer)
		reflection.Register(rserver)
	}

	rpcConfig := tys.impServer.GetRPCConfig()
	ctx, cancel := tys.impServer.GetContext()
	rpcServer, err := rpcserver.NewRPCServer(ctx, cancel, regFunc, rpcConfig, lc)
	if err != nil {
		lc.Errorf("TyDriverCommonServer NewRPCServer err:%s", err)
		return err
	}

	return rpcServer.Serve()
}

/////////////////////////////////////////////////////////////////////////////////////////////
var _ proto.DriverCommonServiceServer = (*TyDriverCommonServer)(nil)

type TyDriverCommonServer struct {
	impServer interfaces.TyDriverService
}

func NewTyDriverCommonServer(implService interfaces.TyDriverService) *TyDriverCommonServer {
	return &TyDriverCommonServer{
		impServer: implService,
	}
}

func (dps *TyDriverCommonServer) RegisterServer(s *grpc.Server) {
	proto.RegisterDriverCommonServiceServer(s, dps)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (tycs *TyDriverCommonServer) AddDeviceCallback(
	ctx context.Context,
	in *proto.DeviceAddInfo) (*emptypb.Empty, error) {

	return AddDeviceCallback(ctx, in, tycs.impServer, commons.ThingModel)
}

func (tycs *TyDriverCommonServer) UpdateDeviceCallback(
	ctx context.Context,
	in *proto.DeviceUpdateInfo) (*emptypb.Empty, error) {

	return UpdateDeviceCallback(ctx, in, tycs.impServer, commons.ThingModel)
}

func (tycs *TyDriverCommonServer) DeleteDeviceCallback(
	ctx context.Context,
	in *proto.DeleteDeviceByIdRequest) (*emptypb.Empty, error) {

	return DeleteDeviceCallback(ctx, in, tycs.impServer, commons.ThingModel)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (tycs *TyDriverCommonServer) ChangeLogLevel(
	ctx context.Context,
	req *proto.LogLevelRequest) (*emptypb.Empty, error) {

	return ChangeLogLevel(req, tycs.impServer)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (tycs *TyDriverCommonServer) AppServiceAddress(
	ctx context.Context,
	req *proto.AppBaseAddress) (*emptypb.Empty, error) {

	return AppServiceAddress(ctx, req, tycs.impServer)
}

//////////////////////////////////////////////////////////////////////////////////////////////
func (tycs *TyDriverCommonServer) GatewayStateCallback(
	ctx context.Context,
	req *proto.GatewayState) (*emptypb.Empty, error) {

	return GatewayStateCallback(req, tycs.impServer)
}

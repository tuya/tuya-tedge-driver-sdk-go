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

// DP模型grpc.Server注册
// TEdge-->SDK
func DPServerRegister(rpcSrv *grpc.Server, dpService interfaces.DpDriverService) {
	NewCommonRPCServer().RegisterServer(rpcSrv)
	NewAppToDriverServer(dpService).RegisterServer(rpcSrv)
	NewDPCommandServer(dpService).RegisterServer(rpcSrv)
	NewDPDriverCommonServer(dpService).RegisterServer(rpcSrv)
}

/////////////////////////////////////////////////////////////////////////////////////////////
type DPDriverRpcServer struct {
	impServer interfaces.DpDriverService
}

func NewDPDriverRpcServer(implService interfaces.DpDriverService) *DPDriverRpcServer {
	return &DPDriverRpcServer{
		impServer: implService,
	}
}

func (dpc *DPDriverRpcServer) Serve() error {
	lc := dpc.impServer.GetLogger()
	regFunc := func(rserver *grpc.Server) {
		DPServerRegister(rserver, dpc.impServer)
		reflection.Register(rserver)
	}

	rpcConfig := dpc.impServer.GetRPCConfig()
	ctx, cancel := dpc.impServer.GetContext()
	rpcServer, err := rpcserver.NewRPCServer(ctx, cancel, regFunc, rpcConfig, lc)
	if err != nil {
		lc.Errorf("DPDriverRpcServer NewRPCServer err:%s", err)
		return err
	}

	return rpcServer.Serve()
}

/////////////////////////////////////////////////////////////////////////////////////////////
var _ proto.DriverCommonServiceServer = (*DPDriverCommonServer)(nil)

type DPDriverCommonServer struct {
	impServer interfaces.DpDriverService
}

func NewDPDriverCommonServer(implService interfaces.DpDriverService) *DPDriverCommonServer {
	return &DPDriverCommonServer{
		impServer: implService,
	}
}

func (dpcs *DPDriverCommonServer) RegisterServer(s *grpc.Server) {
	proto.RegisterDriverCommonServiceServer(s, dpcs)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (dpcs *DPDriverCommonServer) AddDeviceCallback(
	ctx context.Context,
	in *proto.DeviceAddInfo) (*emptypb.Empty, error) {

	return AddDeviceCallback(ctx, in, dpcs.impServer, commons.DPModel)
}

func (dpcs *DPDriverCommonServer) UpdateDeviceCallback(
	ctx context.Context,
	in *proto.DeviceUpdateInfo) (*emptypb.Empty, error) {

	return UpdateDeviceCallback(ctx, in, dpcs.impServer, commons.DPModel)
}

func (dpcs *DPDriverCommonServer) DeleteDeviceCallback(
	ctx context.Context,
	in *proto.DeleteDeviceByIdRequest) (*emptypb.Empty, error) {

	return DeleteDeviceCallback(ctx, in, dpcs.impServer, commons.DPModel)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (dpcs *DPDriverCommonServer) ChangeLogLevel(
	ctx context.Context,
	req *proto.LogLevelRequest) (*emptypb.Empty, error) {

	return ChangeLogLevel(req, dpcs.impServer)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (dpcs *DPDriverCommonServer) AppServiceAddress(
	ctx context.Context,
	req *proto.AppBaseAddress) (*emptypb.Empty, error) {

	return AppServiceAddress(ctx, req, dpcs.impServer)
}

//////////////////////////////////////////////////////////////////////////////////////////////
func (dpcs *DPDriverCommonServer) GatewayStateCallback(
	ctx context.Context,
	req *proto.GatewayState) (*emptypb.Empty, error) {

	return GatewayStateCallback(req, dpcs.impServer)
}

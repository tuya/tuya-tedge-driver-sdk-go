package server

import (
	"context"
	"time"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ proto.DriverCommonServiceServer = (*DPDriverCommonServer)(nil)

type CommonRPCServer struct{}

func NewCommonRPCServer() *CommonRPCServer {
	return &CommonRPCServer{
	}
}

func (crs *CommonRPCServer) RegisterServer(s *grpc.Server) {
	proto.RegisterCommonServer(s, crs)
}

// Ping tests whether the service is working
func (crs *CommonRPCServer) Ping(context.Context, *emptypb.Empty) (*proto.Pong, error) {
	return &proto.Pong{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// Version obtains version information from the target service.
func (crs *CommonRPCServer) Version(context.Context, *emptypb.Empty) (*proto.VersionResponse, error) {
	return &proto.VersionResponse{
		Version: commons.SDKVersion,
	}, nil
}

package rpcserver

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var Kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second,
	PermitWithoutStream: true,
}

// resource --> driver是短连接,以下参数可有可无
var Kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     30 * time.Second,
	MaxConnectionAge:      30 * time.Second,
	MaxConnectionAgeGrace: 5 * time.Second,
	Time:                  5 * time.Second,
	Timeout:               3 * time.Second,
}

type RegisterFunc func(serve *grpc.Server)

type RPCServer struct {
	ctx     context.Context
	cancelF context.CancelFunc

	grpcSrv  *grpc.Server
	listener net.Listener
	logger   commons.TedgeLogger
}

func NewRPCServer(
	ctx context.Context,
	cancel context.CancelFunc,
	registerFunc RegisterFunc,
	rpcConfig config.RPCConfig,
	lc commons.TedgeLogger) (*RPCServer, error) {

	address := rpcConfig.Address
	if address == "" {
		lc.Errorf("NewRPCServer required rpc address")
		return nil, errors.New("required rpc address")
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		lc.Errorf("NewRPCServer listen address:%s failed, err:%s", address, err)
		return nil, err
	}
	lc.Infof("NewRPCServer listen address:%s", address)

	var opts []grpc.ServerOption
	if rpcConfig.UseTLS {
		cred, err := credentials.NewServerTLSFromFile(common.Path(rpcConfig.CertFile), common.Path(rpcConfig.KeyFile))
		if err != nil {
			lc.Errorf("NewRPCServer failed to create credentials: %v", err)
			return nil, err
		}
		opts = append(opts, grpc.Creds(cred))
	}

	// grpc 长连接策略
	opts = append(opts, grpc.KeepaliveEnforcementPolicy(Kaep))
	opts = append(opts, grpc.KeepaliveParams(Kasp))

	rpcSrv := RPCServer{
		listener: lis,
		logger:   lc,
		grpcSrv:  grpc.NewServer(opts...),

		ctx:     ctx,
		cancelF: cancel,
	}

	// registry Server
	registerFunc(rpcSrv.grpcSrv)

	return &rpcSrv, nil
}

func (rpcSrv *RPCServer) Serve() error {
	lc := rpcSrv.logger
	go func() {
		<-rpcSrv.ctx.Done()
		rpcSrv.grpcSrv.Stop()
		lc.Infof("NewRPCServer server shut down done")
	}()

	err := rpcSrv.grpcSrv.Serve(rpcSrv.listener)
	if err != nil {
		lc.Errorf("NewRPCServer server failed: %v", err)
		rpcSrv.cancelF()
		return err
	}

	lc.Infof("NewRPCServer server stopped")
	return nil
}

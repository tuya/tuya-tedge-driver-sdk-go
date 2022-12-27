package clients

import (
	"context"
	"errors"
	"time"

	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/config"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

//grpc client params
var Kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // grpc-go对客户端的ping周期最小限制为10s
	Timeout:             3 * time.Second,  // wait 3 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

var ConnParams = grpc.ConnectParams{
	Backoff: backoff.Config{
		BaseDelay:  time.Second * 1.0,
		Multiplier: 1.0,
		Jitter:     0,
		MaxDelay:   10 * time.Second,
	},
	MinConnectTimeout: time.Second * 3,
}

type NewCliFunc func() *DPReportClient

func dial(cfg config.ClientInfo, serverName string) (*grpc.ClientConn, error) {
	var (
		err         error
		creds       credentials.TransportCredentials
		conn        *grpc.ClientConn
		ctx, cancel = context.WithTimeout(context.Background(), common.GRPCTimeout)
	)
	defer cancel()
	if cfg.UseTLS {
		if creds, err = credentials.NewClientTLSFromFile(common.Path(cfg.CertFilePath), serverName); err != nil {
			return nil, err
		}
		if conn, err = grpc.DialContext(ctx, cfg.Address, grpc.WithTransportCredentials(creds), grpc.WithBlock(), grpc.WithKeepaliveParams(Kacp), grpc.WithConnectParams(ConnParams)); err != nil {
			return nil, err
		}
	} else {
		if conn, err = grpc.DialContext(ctx, cfg.Address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithKeepaliveParams(Kacp), grpc.WithConnectParams(ConnParams)); err != nil {
			return nil, err
		}
	}
	return conn, nil
}

type ResourceClient struct {
	address string
	Conn    *grpc.ClientConn
	proto.CommonClient
	proto.RpcDeviceClient
	proto.RpcDeviceServiceClient
	proto.RpcProductClient
	proto.RpcEventClient
	proto.AlertReportServiceClient  // alert
	proto.RpcGatewayForDeviceClient // gwinfo
	proto.RpcThingModelClient
	proto.DriverStorageClient // leveldb storage
	proto.ThingModelUpServiceClient
	proto.EventReportServiceClient
}

func NewResourceClient(cfg config.ClientInfo) (*ResourceClient, error) {
	var (
		err  error
		conn *grpc.ClientConn
		rc   *ResourceClient
	)
	if cfg.Address == "" {
		return nil, errors.New("required address")
	}
	if conn, err = dial(cfg, common.Resource); err != nil {
		return nil, err
	}
	rc = &ResourceClient{
		address:                   cfg.Address,
		Conn:                      conn,
		CommonClient:              proto.NewCommonClient(conn),
		RpcDeviceClient:           proto.NewRpcDeviceClient(conn),
		RpcDeviceServiceClient:    proto.NewRpcDeviceServiceClient(conn),
		RpcProductClient:          proto.NewRpcProductClient(conn),
		RpcEventClient:            proto.NewRpcEventClient(conn),
		AlertReportServiceClient:  proto.NewAlertReportServiceClient(conn),
		RpcGatewayForDeviceClient: proto.NewRpcGatewayForDeviceClient(conn),
		RpcThingModelClient:       proto.NewRpcThingModelClient(conn),
		DriverStorageClient:       proto.NewDriverStorageClient(conn),
		ThingModelUpServiceClient: proto.NewThingModelUpServiceClient(conn),
		EventReportServiceClient:  proto.NewEventReportServiceClient(conn),
	}
	return rc, nil
}

func (c *ResourceClient) Close() error {
	return c.Conn.Close()
}

type DPReportClient struct {
	Conn *grpc.ClientConn
	proto.CommonClient
	proto.RpcEventClient
}

func NewDPReportClient(cfg config.ClientInfo) (*DPReportClient, error) {
	var (
		err  error
		conn *grpc.ClientConn
	)
	if cfg.Address == "" {
		return nil, errors.New("required address")
	}
	if conn, err = dial(cfg, common.Resource); err != nil {
		return nil, err
	}
	return &DPReportClient{
		Conn:           conn,
		CommonClient:   proto.NewCommonClient(conn),
		RpcEventClient: proto.NewRpcEventClient(conn),
	}, nil
}

func (rdc *DPReportClient) Close() error {
	return rdc.Conn.Close()
}

type TyModelReportClient struct {
	Conn *grpc.ClientConn
	proto.CommonClient
	proto.ThingModelUpServiceClient
}

func NewTyModelReportClient(cfg config.ClientInfo) (*TyModelReportClient, error) {
	var (
		err  error
		conn *grpc.ClientConn
	)
	if cfg.Address == "" {
		return nil, errors.New("required address")
	}
	if conn, err = dial(cfg, common.Resource); err != nil {
		return nil, err
	}
	return &TyModelReportClient{
		Conn:                      conn,
		CommonClient:              proto.NewCommonClient(conn),
		ThingModelUpServiceClient: proto.NewThingModelUpServiceClient(conn),
	}, nil
}

func (tmrc *TyModelReportClient) Close() error {
	return tmrc.Conn.Close()
}

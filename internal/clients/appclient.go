package clients

import (
	"fmt"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

const (
	rpcScheme = "driver"
)

type AppService struct {
	Name      string
	RpcClient AppClient
}

type AppCallBack func(appName string, req commons.AppDriverReq) (commons.Response, error)

func dail(builder *AppResolverBuilder) (*grpc.ClientConn, error) {

	resolver.Register(builder)

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", rpcScheme, builder.AppSerivceName), // Dial to "example:///resolver.example.grpc.io"
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(Kacp),
		grpc.WithConnectParams(ConnParams),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type AppResolverBuilder struct {
	AppSerivceName string
	Addr           string
	AddrChan       chan string
}

func (e *AppResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &appResolver{
		addr:   e.Addr,
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			e.AppSerivceName: {e.Addr},
		},
		addrChan:    e.AddrChan,
		serviceName: e.AppSerivceName,
	}
	r.start()
	return r, nil
}
func (*AppResolverBuilder) Scheme() string { return rpcScheme }

type appResolver struct {
	serviceName string
	addr        string
	target      resolver.Target
	cc          resolver.ClientConn
	addrsStore  map[string][]string
	addrChan    chan string
}

func (r *appResolver) start() {
	//动态更新驱动endpoint
	addr := make([]resolver.Address, 0)
	addr = append(addr, resolver.Address{Addr: r.addr})
	r.cc.UpdateState(resolver.State{Addresses: addr})
	go func() {
		for {
			a := <-r.addrChan
			addr = addr[:0]
			addr = append(addr, resolver.Address{Addr: a})
			r.cc.UpdateState(resolver.State{Addresses: addr})
		}
	}()
}
func (*appResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*appResolver) Close()                                  {}

type AppClientCnnMap map[int]*AppClient //rpc连接数和AppClient对应关系

type AppClient struct {
	Builder *AppResolverBuilder
	Conn    *grpc.ClientConn
	proto.DriverToAppServiceClient
}

func NewAppRpcClient(builder *AppResolverBuilder) (*AppClient, error) {
	conn, err := dail(builder)
	if err != nil {
		return nil, err
	}
	return &AppClient{
		Conn:                     conn,
		DriverToAppServiceClient: proto.NewDriverToAppServiceClient(conn),
	}, nil
}

func (c *AppClient) Close() error {
	return c.Conn.Close()
}

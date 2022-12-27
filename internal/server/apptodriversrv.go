package server

import (
	"context"
	"fmt"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/interfaces"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"google.golang.org/grpc"
)

var _ proto.AppToDriverServiceServer = (*AppToDriverServer)(nil)

type AppToDriverServer struct {
	impServer interfaces.DriverCommonItf
}

func NewAppToDriverServer(impServer interfaces.DriverCommonItf) *AppToDriverServer {
	return &AppToDriverServer{
		impServer: impServer,
	}
}

func (app *AppToDriverServer) RegisterServer(s *grpc.Server) {
	proto.RegisterAppToDriverServiceServer(s, app)
}

func (ds *AppToDriverServer) SendToDriver(ctx context.Context, req *proto.Data) (*proto.SendResponse, error) {
	lc := ds.impServer.GetLogger()
	lc.Debugf("SendToDriver receive data: %v", req)

	appHandler := ds.impServer.GetAppHandler()
	if appHandler == nil {
		return &proto.SendResponse{
			Message: fmt.Sprintf("SendToDriver appService %s not found", req.Name),
		}, nil
	}

	resp, err := appHandler(req.Name, commons.AppDriverReq{
		Header: &commons.Header{
			Tag:    req.Header.Tag,
			From:   req.Header.From,
			Option: req.Header.Option,
		},
		Payload: req.Payload,
	})

	if err != nil {
		return &proto.SendResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &proto.SendResponse{
		Success: resp.Success,
		Message: resp.Message,
		Payload: resp.Payload,
	}, nil
}

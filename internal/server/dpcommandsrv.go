package server

import (
	"context"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/interfaces"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/transform"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DPCommandServer struct {
	impServer interfaces.DpDriverService
}

func NewDPCommandServer(impServer interfaces.DpDriverService) *DPCommandServer {
	return &DPCommandServer{
		impServer: impServer,
	}
}

func (dpc *DPCommandServer) RegisterServer(s *grpc.Server) {
	proto.RegisterDPModelDriverServiceServer(s, dpc)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (dpc *DPCommandServer) UpdateProductCallback(ctx context.Context, in *proto.Product) (*emptypb.Empty, error) {
	lc := dpc.impServer.GetLogger()

	pid := in.GetId()
	driver := dpc.impServer.GetDriver()
	product := transform.ToDPProductModel(in)
	pdCache := dpc.impServer.GetPdCache()
	_, ok := pdCache.ById(pid)
	if !ok {
		lc.Debugf("product id:%s not found, add it", pid)
		pdCache.Add(product)
		if err := driver.ProductNotify(ctx, commons.ProductAddNotify, pid, product); err != nil {
			lc.Errorf("UpdateProductCallback product:%s notify error: %s", pid, err)
			return common.EmptyPb, status.Errorf(codes.Internal, err.Error())
		}
	} else {
		pdCache.Update(product)
		if err := driver.ProductNotify(ctx, commons.ProductUpdateNotify, pid, product); err != nil {
			lc.Errorf("UpdateProductCallback product:%s notify error: %s", pid, err)
			return common.EmptyPb, status.Errorf(codes.Internal, err.Error())
		}
	}

	lc.Infof("UpdateProductCallback product:%s updated", pid)
	return common.EmptyPb, nil
}

func (dpc *DPCommandServer) IssueCommand(ctx context.Context, in *proto.CmdRequest) (*emptypb.Empty, error) {
	lc := dpc.impServer.GetLogger()

	var (
		err       error
		exist     bool
		cid       string
		dev       commons.DeviceInfo
		dp        dpmodel.DP
		data      map[string]interface{}
		req       dpmodel.CommandRequest
		protocols = make(map[string]commons.ProtocolProperties)
		dpExtend  = make(dpmodel.DPExtendInfo)
	)

	req.Protocol = int32(in.GetProtocol())
	req.T = in.GetT()
	req.S = in.GetS()

	if err = utils.JsonDecoder(in.GetData(), &data); err != nil {
		lc.Errorf("DP IssueCommand decode error: %s", err)
		return common.EmptyPb, status.Errorf(codes.Internal, err.Error())
	}
	req.Data = data
	lc.Infof("DP IssueCommand get command:%+v", req)

	cid = in.GetCid()
	devCache := dpc.impServer.GetDevCache()
	dev, exist = devCache.ById(cid)
	if !exist {
		lc.Errorf("DP IssueCommand cid:%s not found in local cache", in.Cid)
		return common.EmptyPb, status.Error(codes.NotFound, "device not found in local cache")
	}
	protocols = dev.Protocols

	switch in.GetProtocol() {
	case proto.EnumTUYAMQTTProtocol_ENUM_TUYAMQTT_PROTOCOL_5:
		dps, ok := data["dps"]
		if !ok {
			lc.Warnf("DP IssueCommand dps not found, data:%v", data)
		} else {
			if _, ok = dps.(map[string]interface{}); !ok {
				lc.Errorf("DP IssueCommand dps type error")
				return common.EmptyPb, status.Error(codes.InvalidArgument, "dps type error")
			}

			pdCache := dpc.impServer.GetPdCache()
			for k := range dps.(map[string]interface{}) {
				if dp, exist = pdCache.Dp(dev.ProductId, k); !exist {
					lc.Warnf("DP IssueCommand dp(%s) not found in product(%s)", k, dev.ProductId)
					continue
				} else {
					dpExtend[k] = dpmodel.DPExtend{
						Property: dp.Properties,
						Attr:     dp.Attributes,
					}
				}
			}
		}
	case proto.EnumTUYAMQTTProtocol_ENUM_TUYAMQTT_PROTOCOL_64,
		proto.EnumTUYAMQTTProtocol_ENUM_TUYAMQTT_PROTOCOL_302,
		proto.EnumTUYAMQTTProtocol_ENUM_TUYAMQTT_PROTOCOL_312:
	default:
		lc.Errorf("DP IssueCommand unsupported protocol: %+v", in.GetProtocol())
		return common.EmptyPb, status.Error(codes.Unimplemented, "unsupported protocol")
	}

	driver := dpc.impServer.GetDriver()
	if err = driver.HandleCommands(ctx, cid, req, protocols, dpExtend); err != nil {
		lc.Errorf("DP IssueCommand handle command error: %s", err)
		return common.EmptyPb, status.Error(codes.Internal, err.Error())
	}

	return common.EmptyPb, nil
}

package server

import (
	"context"
	"fmt"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/common"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/interfaces"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/transform"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/utils"
	"github.com/tuya/tuya-tedge-driver-sdk-go/thingmodel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TyCommandServer struct {
	impServer interfaces.TyDriverService
}

func NewTyCommandServer(impServer interfaces.TyDriverService) *TyCommandServer {
	return &TyCommandServer{
		impServer: impServer,
	}
}

func (tcs *TyCommandServer) RegisterServer(s *grpc.Server) {
	proto.RegisterThingModelDownServiceServer(s, tcs)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (tcs *TyCommandServer) TMUpdateProductCallback(ctx context.Context, in *proto.ThingModelProduct) (*emptypb.Empty, error) {
	pdCache := tcs.impServer.GetPdCache()
	product := transform.ToTMProductModel(in)
	pdCache.Update(product)

	lc := tcs.impServer.GetLogger()
	pid := in.GetId()
	driver := tcs.impServer.GetDriver()
	if err := driver.ProductNotify(ctx, commons.ProductUpdateNotify, pid, product); err != nil {
		lc.Errorf("TMUpdateProductCallback product:%s ProductNotify err:%s", in.Id, err)
		return common.EmptyPb, status.Errorf(codes.Internal, err.Error())
	}

	lc.Debugf("TMUpdateProductCallback product:%s updated", in.Id)
	return common.EmptyPb, nil
}

func (tcs *TyCommandServer) ThingModelMsgIssue(ctx context.Context, in *proto.ThingModelMsg) (*emptypb.Empty, error) {
	lc := tcs.impServer.GetLogger()
	lc.Debugf("ThingModelMsgIssue receive msg: cid: %s, op: %d, data: %s", in.GetCid(), in.GetOpType(), in.GetData())

	devCache := tcs.impServer.GetDevCache()
	cid := in.GetCid()
	dev, ok := devCache.ById(cid)
	if !ok {
		lc.Errorf("ThingModelMsgIssue can't find cid: %s in local cache", cid)
		return common.EmptyPb, status.Errorf(codes.NotFound, "can't find cid: %s in local cache", cid)
	}

	driver := tcs.impServer.GetDriver()
	pdCache := tcs.impServer.GetPdCache()
	switch in.GetOpType() {
	case proto.ThingModelOperationType_PROPERTY_SET: // 设备属性下发
		var req thingmodel.PropertySet
		if err := utils.JsonDecoder([]byte(in.GetData()), &req); err != nil {
			lc.Errorf("ThingModelMsgIssue decode data error: %s", err)
			return common.EmptyPb, status.Errorf(codes.Internal, "decode data error: %s", err)
		}
		lc.Debugf("ThingModelMsgIssue req: %+v", req)

		req.Spec = make(map[string]thingmodel.PropertySpec, len(req.Data))
		for k := range req.Data {
			if ps, ok := pdCache.GetPropertySpecByCode(dev.ProductId, k); !ok {
				lc.Warnf("ThingModelMsgIssue can't find property(%s) spec in product(%s)", k, dev.ProductId)
				continue
			} else {
				req.Spec[k] = ps
			}
		}

		err := driver.HandlePropertySet(ctx, cid, req, dev.Protocols)
		if err != nil {
			lc.Errorf("ThingModelMsgIssue handlePropertySet error: %s", err)
			return common.EmptyPb, status.Errorf(codes.Unknown, err.Error())
		}
	case proto.ThingModelOperationType_PROPERTY_GET:
		var req thingmodel.PropertyGet
		if err := utils.JsonDecoder([]byte(in.GetData()), &req); err != nil {
			lc.Errorf("decode data error: %s", err)
			return common.EmptyPb, status.Errorf(codes.Internal, "decode data error: %s", err)
		}

		req.Spec = make(map[string]thingmodel.PropertySpec, len(req.Data))
		for _, k := range req.Data {
			if ps, ok := pdCache.GetPropertySpecByCode(dev.ProductId, k); !ok {
				lc.Warnf("can't find property(%s) spec in product(%s)", k, dev.ProductId)
				continue
			} else {
				req.Spec[k] = ps
			}
		}

		err := driver.HandlePropertyGet(ctx, cid, req, dev.Protocols)
		if err != nil {
			lc.Errorf("handlePropertyGet error: %s", err)
			return common.EmptyPb, status.Errorf(codes.Unknown, err.Error())
		}
	case proto.ThingModelOperationType_ACTION_EXECUTE:
		var req thingmodel.ActionExecuteRequest
		if err := utils.JsonDecoder([]byte(in.GetData()), &req); err != nil {
			lc.Errorf("decode data error: %s", err)
			return common.EmptyPb, status.Errorf(codes.Internal, "decode data error: %s", err)
		}

		if action, ok := pdCache.GetActionSpecByCode(dev.ProductId, req.Data.ActionCode); !ok {
			lc.Warnf("can't find action(%s) spec in product(%s)", req.Data.ActionCode, dev.ProductId)
		} else {
			req.Spec = action
		}

		err := driver.HandleActionExecute(ctx, cid, req, dev.Protocols)
		if err != nil {
			lc.Errorf("handleActionExecute error: %s", err)
		}
	case proto.ThingModelOperationType_PROPERTY_REPORT_RESPONSE,
		proto.ThingModelOperationType_EVENT_TRIGGER_RESPONSE,
		proto.ThingModelOperationType_DATA_BATCH_REPORT_RESPONSE:
		//proto.ThingModelOperationType_PROPERTY_SET_RESPONSE:
		var resp thingmodel.CommonResponse
		if err := utils.JsonDecoder([]byte(in.GetData()), &resp); err != nil {
			lc.Errorf("decode data error: %s", err)
			return common.EmptyPb, status.Errorf(codes.Internal, "decode data error: %s", err)
		}
		if err := tcs.handleAck(resp.MsgId, resp); err != nil {
			return common.EmptyPb, status.Errorf(codes.NotFound, err.Error())
		}
	case proto.ThingModelOperationType_PROPERTY_DESIRED_GET_RESPONSE:
		var resp thingmodel.PropertyDesiredGetResponse
		if err := utils.JsonDecoder([]byte(in.GetData()), &resp); err != nil {
			lc.Errorf("decode data error: %s", err)
			return common.EmptyPb, status.Errorf(codes.Internal, "decode data error: %s", err)
		}
		if err := tcs.handleAck(resp.MsgId, resp); err != nil {
			return common.EmptyPb, status.Errorf(codes.NotFound, err.Error())
		}
	case proto.ThingModelOperationType_PROPERTY_DESIRED_DELETE_RESPONSE:
		var resp thingmodel.PropertyDesiredDeleteResponse
		if err := utils.JsonDecoder([]byte(in.GetData()), &resp); err != nil {
			lc.Errorf("decode data error: %s", err)
			return common.EmptyPb, status.Errorf(codes.Internal, "decode data error: %s", err)
		}
		if err := tcs.handleAck(resp.MsgId, resp); err != nil {
			return common.EmptyPb, status.Errorf(codes.NotFound, err.Error())
		}
	default:
		return common.EmptyPb, status.Errorf(codes.InvalidArgument, "unsupported operation type")
	}
	return common.EmptyPb, nil
}

func (tycs *TyCommandServer) handleAck(id string, resp interface{}) error {
	lc := tycs.impServer.GetLogger()
	workerPool := tycs.impServer.GetTMWorkerPool()
	ack, ok := workerPool.LoadMsgChan(id)
	if !ok {
		lc.Errorf("handleAck can't find ackChan with id: %s", id)
		return fmt.Errorf("can't find msgId: %s", id)
	}
	if !ack.(*common.MsgAckChan).TrySendDataAndCloseChan(resp) {
		lc.Errorf("handleAck send data to chan error, because chan is closed")
	}

	workerPool.DeleteMsgId(id)
	return nil
}

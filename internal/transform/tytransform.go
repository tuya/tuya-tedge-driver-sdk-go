package transform

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"github.com/tuya/tuya-tedge-driver-sdk-go/thingmodel"
)

type Optype int8

const (
	PropertyReport Optype = iota + 1
	PropertySetResponse
	PropertyGetResponse
	ActionExecuteResponse
	EventReport
	BatchReport
	PropertyDesiredGet
	PropertyDesiredDelete
)

func TMDataToProto(cid string, t Optype, data interface{}) (*proto.ThingModelMsg, error) {
	var opt proto.ThingModelOperationType
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	switch t {
	case PropertyReport:
		opt = proto.ThingModelOperationType_PROPERTY_REPORT
	case PropertySetResponse:
		opt = proto.ThingModelOperationType_PROPERTY_SET_RESPONSE
	case PropertyGetResponse:
		opt = proto.ThingModelOperationType_PROPERTY_GET_RESPONSE
	case ActionExecuteResponse:
		opt = proto.ThingModelOperationType_ACTION_EXECUTE_RESPONSE
	case EventReport:
		opt = proto.ThingModelOperationType_EVENT_TRIGGER
	case BatchReport:
		opt = proto.ThingModelOperationType_DATA_BATCH_REPORT
	case PropertyDesiredGet:
		opt = proto.ThingModelOperationType_PROPERTY_DESIRED_GET
	case PropertyDesiredDelete:
		opt = proto.ThingModelOperationType_PROPERTY_DESIRED_DELETE
	default:
		return nil, errors.New("unsupported optype")
	}
	return &proto.ThingModelMsg{
		Cid:    cid,
		OpType: opt,
		Data:   string(payload),
	}, nil
}

func ProtoCommonRespToTMData(data string) (thingmodel.CommonResponse, error) {
	var resp thingmodel.CommonResponse
	d := json.NewDecoder(strings.NewReader(data))
	d.UseNumber()
	err := d.Decode(&resp)
	return resp, err
}

func FromAddTyProductModel2Proto(p *thingmodel.AddProductReq, driverLibraryId string) (*proto.Product, error) {
	if p == nil {
		return nil, errors.New("product point is nil")
	}

	return &proto.Product{
		Id:              p.Id,
		Name:            p.Name,
		Description:     p.Description,
		Model:           p.Model,
		DeviceLibraryId: driverLibraryId,
	}, nil
}

func HttpWithCustomParamsToProto(params thingmodel.HttpRequestParam, payload map[string]interface{}) (*proto.HttpRequestParam, error) {
	var (
		err  error
		buff []byte
	)
	if buff, err = json.Marshal(payload); err != nil {
		return nil, err
	}
	return &proto.HttpRequestParam{
		HttpUrl: params.Url,
		HttpApi: params.Api,
		Version: params.Version,
		Payload: buff,
	}, nil
}

func FromTMDeviceToDevice(device commons.TMDeviceMeta) commons.DeviceMeta {
	dev := commons.DeviceMeta{
		Cid:          device.Cid,
		ProductId:    device.ProductId,
		BaseAttr:     device.BaseAttr,
		ExtendedAttr: device.ExtendedAttr,
		Protocols:    device.Protocols,
	}

	return dev
}

func FromDeviceToTMDevice(device commons.DeviceMeta) commons.TMDeviceMeta {
	dev := commons.TMDeviceMeta{
		Cid:          device.Cid,
		ProductId:    device.ProductId,
		BaseAttr:     device.BaseAttr,
		ExtendedAttr: device.ExtendedAttr,
		Protocols:    device.Protocols,
	}

	return dev
}

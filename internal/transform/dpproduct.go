package transform

import (
	"errors"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
)

func ToDPProductModel(p *proto.Product) dpmodel.DPModelProduct {
	return dpmodel.DPModelProduct{
		Id:              p.GetId(),
		Name:            p.GetName(),
		//Description:     p.GetDescription(),
		Model:           p.GetModel(),
		Dps:             ToDPModels(p.GetDps()),
		DeviceLibraryId: p.GetDeviceLibraryId(),
	}
}

func ToDPModels(dp []*proto.DP) []dpmodel.DP {
	deviceResourceModels := make([]dpmodel.DP, len(dp))
	for i, d := range dp {
		deviceResourceModels[i] = ToDPModel(d)
	}
	return deviceResourceModels
}

func ToDPModel(dp *proto.DP) dpmodel.DP {
	return dpmodel.DP{
		//Description: dp.GetDescription(),
		Id:          dp.GetId(),
		Properties:  ToPropertyValueModel(dp.Properties),
		Attributes:  dp.Attributes,
	}
}

func ToPropertyValueModel(p *proto.PropertyValue) dpmodel.PropertyValue {
	var (
		enum2  []string
		fault2 []string
	)
	enum1 := p.GetEnum()
	if len(enum1) > 32 {
		enum2 = make([]string, 32)
		copy(enum2, enum1)
	} else {
		enum2 = enum1
	}

	fault1 := p.GetFault()
	if len(fault1) > 30 {
		fault2 = make([]string, 30)
		copy(fault2, fault1)
	} else {
		fault2 = fault1
	}

	if enum2 == nil {
		enum2 = make([]string, 0)
	}
	if fault2 == nil {
		fault2 = make([]string, 0)
	}

	return dpmodel.PropertyValue{
		Type:         commons.DataType(p.Type),
		ReadWrite:    p.ReadWrite,
		Units:        p.Units,
		Minimum:      p.Minimum,
		Maximum:      p.Maximum,
		DefaultValue: p.DefaultValue,
		Shift:        p.Shift,
		Scale:        p.Scale,
		Enum:         enum2,
		Fault:        fault2,
	}
}

func FromAddProductModel2Proto(p *dpmodel.DPModelProductAddInfo, driverLibraryId string) (*proto.Product, error) {
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

package transform

import (
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
	"github.com/tuya/tuya-tedge-driver-sdk-go/thingmodel"
)

func ToTMProductModel(p *proto.ThingModelProduct) thingmodel.ThingModelProduct {
	return thingmodel.ThingModelProduct{
		Id:              p.GetId(),
		Name:            p.GetName(),
		//Description:     p.GetDescription(),
		Model:           p.GetModel(),
		Action:          ToTMActionModels(p.GetAction()),
		Event:           ToTMEventModels(p.GetEvent()),
		Property:        ToTMPropertyModels(p.GetProperty()),
		DeviceLibraryId: p.GetDeviceLibraryId(),
	}
}

func ToTMActionModels(as []*proto.Action) []thingmodel.Action {
	rets := make([]thingmodel.Action, 0, len(as))
	for i := range as {
		rets = append(rets, thingmodel.Action{
			AbilityId:    as[i].GetAbilityId(),
			Code:         as[i].GetCode(),
			Name:         as[i].GetName(),
			//Description:  as[i].GetDescription(),
			InputParams:  ToTMParamsModels(as[i].GetInputParams()),
			OutputParams: ToTMParamsModels(as[i].GetOutputParams()),
			Attributes:   as[i].GetAttributes(),
		})
	}
	return rets
}

func ToTMParamsModels(params []*proto.InputOutput) []thingmodel.InputOutput {
	rets := make([]thingmodel.InputOutput, 0, len(params))
	for i := range params {
		rets = append(rets, thingmodel.InputOutput{
			Code:     params[i].GetCode(),
			Name:     params[i].GetName(),
			TypeSpec: ToTMTypeSpecModels(params[i].GetTypeSpec()),
		})
	}
	return rets
}

func ToTMTypeSpecModels(spec *proto.TypeSpec) thingmodel.TypeSpec {
	return thingmodel.TypeSpec{
		Type:         commons.DataType(spec.GetType()),
		Min:          spec.GetMin(),
		Max:          spec.GetMax(),
		Step:         spec.GetStep(),
		Unit:         spec.GetUnit(),
		Scale:        spec.GetScale(),
		MaxLen:       spec.GetMaxLen(),
		Range:        spec.GetRange(),
		Label:        spec.GetLabel(),
		DefaultValue: spec.GetDefaultValue(),
	}
}

func ToTMEventModels(e []*proto.TMEvent) []thingmodel.Event {
	rets := make([]thingmodel.Event, 0, len(e))
	for i := range e {
		rets = append(rets, thingmodel.Event{
			AbilityId:    e[i].GetAbilityId(),
			Code:         e[i].GetCode(),
			Name:         e[i].GetName(),
			//Description:  e[i].GetDescription(),
			OutputParams: ToTMParamsModels(e[i].GetOutputParams()),
			Attributes:   e[i].GetAttributes(),
		})
	}
	return rets
}

func ToTMPropertyModels(p []*proto.Property) []thingmodel.PropertySpec {
	rets := make([]thingmodel.PropertySpec, 0, len(p))
	for i := range p {
		rets = append(rets, thingmodel.PropertySpec{
			AbilityId:  p[i].GetAbilityId(),
			Code:       p[i].GetCode(),
			Name:       p[i].GetName(),
			AccessMode: p[i].GetAccessMode(),
			TypeSpec:   ToTMTypeSpecModels(p[i].GetTypeSpec()),
			Attributes: p[i].GetAttributes(),
		})
	}
	return rets
}

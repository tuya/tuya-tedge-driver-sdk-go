package thingmodel

import "github.com/tuya/tuya-tedge-driver-sdk-go/commons"

type (
	// ThingModelProduct 物模型产品定义
	ThingModelProduct struct {
		Id              string
		Name            string
		Description     string
		Model           string
		Action          []Action
		Event           []Event
		Property        []PropertySpec
		DeviceLibraryId string
	}

	// Action 物模型动作定义
	Action struct {
		AbilityId    string
		ProductId    string
		Code         string
		Name         string
		Description  string
		InputParams  []InputOutput
		OutputParams []InputOutput
		Attributes   map[string]string // 自定义
	}

	InputOutput struct {
		Code     string
		Name     string
		TypeSpec TypeSpec
	}

	TypeSpec struct {
		Properties      map[string]PropertySpec
		ElementTypeSpec *TypeSpec
		Type            commons.DataType
		Min             int64
		Max             int64
		Step            int64
		Unit            string
		Scale           int64
		MaxLen          int64
		Range           []string
		Label           []string
		DefaultValue    string
	}

	// Event 物模型事件定义
	Event struct {
		AbilityId    string
		ProductId    string
		Code         string
		Name         string
		Description  string
		OutputParams []InputOutput
		Attributes   map[string]string // 自定义
	}

	// PropertySpec 物模型属性定义
	PropertySpec struct {
		AbilityId  string
		ProductId  string
		Code       string
		Name       string
		AccessMode string
		TypeSpec   TypeSpec
		Attributes map[string]string // 自定义
	}

	// AddProductReq 添加产品
	AddProductReq struct {
		Id          string
		Name        string
		Model       string
		Description string
	}
)

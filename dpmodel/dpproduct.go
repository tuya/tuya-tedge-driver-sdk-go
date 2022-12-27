package dpmodel

import "github.com/tuya/tuya-tedge-driver-sdk-go/commons"

type DPModelProductAddInfo struct {
	Id          string //
	Name        string //产品名称
	Model       string //产品型号
	Description string //产品描述
}

type DPModelProduct struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"-"`
	Manufacturer    string `json:"manufacturer"`
	Model           string `json:"model"`
	Dps             []DP
	DeviceLibraryId string `json:"deviceLibraryId"`
}

type DP struct {
	Description string
	Id          string
	Properties  PropertyValue
	Attributes  map[string]string // must have valueType if Properties.Type=value, one of (int8,int16,int32,int64,uint8,uint16,uint32,uint64,float32,float64)
}

type PropertyValue struct {
	Type         commons.DataType // one of value, string, enum, bool, fault, raw
	ReadWrite    string
	Units        string
	Minimum      int64
	Maximum      int64
	DefaultValue string
	Shift        string
	Scale        string
	Offset       string
	Enum         []string // limit length 32
	Fault        []string // limit length 30
}

package transform

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/cache"
	"github.com/tuya/tuya-tedge-driver-sdk-go/internal/proto"
)

const (
	DefaultOffset  = 0
	DefaultLimit   = 20
	MaxBinaryBytes = 0xFF
	FaultMask      = 0x3FFFFFFF
)

type (
	// DPValue dp点对应的值
	DPValue struct {
		// dpId is the name of Device Resource for this command
		DPId string
		// Origin is an int64 value which indicates the time the reading
		// contained in the DPValue was read by the ProtocolDriver
		// instance.
		Origin int64
		// Type is a ValueType value which indicates what type of
		// value was returned from the ProtocolDriver instance in
		// response to HandleCommand being called to handle a single
		// ResourceOperation.
		Type commons.DataType
		// NumericValue is a byte slice with a maximum capacity of
		// 64 bytes, used to hold a numeric value returned by a
		// ProtocolDriver instance. The value can be converted to
		// its native type by referring to the the value of ResType.
		NumericValue []byte
		// stringValue is a string value returned as a value by a ProtocolDriver instance.
		stringValue string
		// BinValue is a binary value with a maximum capacity of 16 MB,
		// used to hold binary values returned by a ProtocolDriver instance.
		BinValue []byte
	}

	WithDPData struct {
		Id        string
		DeviceId  string
		ProductId string
		Created   int64
		Origin    int64
		Readings  []Reading
	}
	WithDPReport struct {
		Cid  string
		Data []*DPValue
	}

	WithoutDPReport struct {
		Topic string
		Data  *dpmodel.WithoutDPValue
	}
	BaseReading struct {
		Id          string
		Created     int64
		Origin      int64
		DeviceId    string
		DpId        string
		ProductId   string
		ValueType   commons.DataType
		Value       string
		BinaryValue []byte
		//MediaType   string
	}

	BinaryReading struct {
		BaseReading `json:",inline"`
		BinaryValue []byte
		//MediaType   string
	}

	SimpleReading struct {
		BaseReading `json:",inline"`
		Value       string
	}
)

type Reading interface {
	GetBaseReading() BaseReading
}

func (b BinaryReading) GetBaseReading() BaseReading { return b.BaseReading }
func (s SimpleReading) GetBaseReading() BaseReading { return s.BaseReading }

func WithDPValueToReading(dpv *DPValue, cid, productId string) Reading {
	baseReading := BaseReading{
		Id:        uuid.New().String(),
		DeviceId:  cid,
		DpId:      dpv.DPId,
		ProductId: productId,
		Created:   time.Now().Unix(),
		ValueType: dpv.Type,
	}
	if dpv.Origin > 0 {
		baseReading.Origin = dpv.Origin
	} else {
		baseReading.Origin = time.Now().Unix()
	}
	if dpv.Type == commons.RawType {
		return BinaryReading{
			BaseReading: baseReading,
			BinaryValue: dpv.BinValue,
		}
	} else {
		value, err := dpv.ValueToString()
		if err != nil {
			return nil
		}
		return SimpleReading{
			BaseReading: baseReading,
			Value:       value,
		}
	}
}

func NewWithDPData(cid string, data []*dpmodel.WithDPValue, dev commons.DeviceInfo, pc *cache.DPModelProductCache) (*WithDPData, error) {
	// check dp value
	var (
		err   error
		ok    bool
		dp    dpmodel.DP
		value *DPValue
		v     interface{}
	)

	readings := make([]Reading, 0, len(data))
	for _, dpv := range data {
		dpId := dpv.DPId
		if dpv == nil {
			return nil, fmt.Errorf("cid:%s dpid:%s data slice has nil value", cid, dpId)
		}
		v = dpv.Value

		// check
		if dp, ok = pc.Dp(dev.ProductId, dpId); !ok {
			return nil, fmt.Errorf("dpid(%s) not found in product(%s)", dpId, dev.ProductId)
		}
		if strings.Compare(string(dp.Properties.Type), string(dpv.DPType)) != 0 {
			return nil, fmt.Errorf("dp(%s) value type error, want: %s, get: %s", dpId, dp.Properties.Type, dpv.DPType)
		}

		// check value
		if dp.Properties.Type == commons.ValueType {
			vv, ok := v.(int64)
			if !ok {
				return nil, fmt.Errorf("dpId:%s value type error, want int64", dpId)
			}
			if vv < dp.Properties.Minimum || vv > dp.Properties.Maximum {
				return nil, fmt.Errorf("dpId:%s value out of range, want: [%d, %d], got: %d", dpId, dp.Properties.Minimum, dp.Properties.Maximum, vv)
			}
		} else if dp.Properties.Type == commons.FaultType || dp.Properties.Type == commons.BitmapType {
			// transform
			switch dpv.Value.(type) {
			case int32:
				break
			case []string:
				if v = faultTransform(dpv.Value.([]string), dp.Properties.Fault); v.(int32) <= 0 {
					continue
				}
			default:
				return nil, fmt.Errorf("dpId:%s unknown value:%v type: %v", dpId, dpv.Value, reflect.TypeOf(dpv.Value))
			}
		}

		if value, err = NewDPValue(dpv.DPId, dpv.DPType, v); err != nil {
			return nil, err
		}
		reading := WithDPValueToReading(value, dev.Cid, dev.ProductId)
		if reading == nil {
			return nil, fmt.Errorf("dpId:%s value to reading return nil", dpId)
		}
		readings = append(readings, reading)
	}

	if len(readings) <= 0 {
		return nil, fmt.Errorf("cid:%s, reading is empty", cid)
	}
	return &WithDPData{
		Id:        uuid.New().String(),
		DeviceId:  dev.Cid,
		ProductId: dev.ProductId,
		Created:   time.Now().Unix(),
		Origin:    time.Now().Unix(),
		Readings:  readings,
	}, nil
}

func faultTransform(value, fault []string) int32 {
	var result int32
	for i := range value {
		if len(value[i]) == 0 {
			continue
		}
		for j := range fault {
			if value[i] == fault[j] {
				result |= 1 << j
				break
			}
		}
	}
	return result & FaultMask
}

func (d *WithDPData) ToProto() *proto.Event {
	var readings []*proto.BaseReading
	for _, reading := range d.Readings {
		readings = append(readings, toProtoReading(reading))
	}
	return &proto.Event{
		Id:        d.Id,
		DeviceId:  d.DeviceId,
		ProductId: d.ProductId,
		Created:   d.Created,
		Origin:    d.Origin,
		Readings:  readings,
	}
}

func toProtoReading(reading Reading) *proto.BaseReading {
	var baseReading *proto.BaseReading
	switch reading.(type) {
	case BinaryReading:
		r := reading.(BinaryReading)
		baseReading = &proto.BaseReading{
			Id:        r.Id,
			Created:   r.Created,
			Origin:    r.Origin,
			DeviceId:  r.DeviceId,
			DpId:      r.DpId,
			ProductId: r.ProductId,
			ValueType: string(r.ValueType),
			Reading: &proto.BaseReading_BinaryReading{
				BinaryReading: &proto.BinaryReading{
					BinaryValue: r.BinaryValue,
				},
			},
		}
	case SimpleReading:
		r := reading.(SimpleReading)
		baseReading = &proto.BaseReading{
			Id:        r.Id,
			Created:   r.Created,
			Origin:    r.Origin,
			DeviceId:  r.DeviceId,
			DpId:      r.DpId,
			ProductId: r.ProductId,
			ValueType: string(r.ValueType),
			Reading: &proto.BaseReading_SimpleReading{
				SimpleReading: &proto.SimpleReading{
					Value: r.Value,
				},
			},
		}
	}
	return baseReading
}

func WithoutDPDataToProto(data *WithoutDPReport) (*proto.WithoutDpReport, error) {
	var (
		err  error
		buff []byte
	)
	if data == nil {
		return nil, errors.New("nil")
	}
	if buff, err = json.Marshal(data.Data.Data); err != nil {
		return nil, err
	}
	return &proto.WithoutDpReport{
		Topic:    data.Topic,
		Protocol: data.Data.Protocol,
		S:        data.Data.S,
		T:        data.Data.T,
		Data:     buff,
	}, nil
}

func HTTPDataReportToProto(api, version string, payload map[string]interface{}) (*proto.HttpReport, error) {
	var (
		err  error
		buff []byte
	)
	if buff, err = json.Marshal(payload); err != nil {
		return nil, err
	}
	return &proto.HttpReport{
		HttpApi: api,
		Version: version,
		Payload: buff,
	}, nil
}

func DeviceStatusToProto(lists *commons.DeviceStatus) (*proto.DeviceOnlineAndOfflineList, error) {
	if len(lists.Online) == 0 && len(lists.Offline) == 0 {
		return nil, errors.New("list is nil")
	}
	return &proto.DeviceOnlineAndOfflineList{
		Online:  lists.Online,
		Offline: lists.Offline,
	}, nil
}

func NewDPValue(dpId string, dpType commons.DataType, value interface{}) (*DPValue, error) {
	var (
		err    = errors.New("type error")
		origin = time.Now().Unix()
	)
	switch dpType {
	case commons.BoolType:
		v, ok := value.(bool)
		if !ok {
			return nil, err
		}
		return NewBoolValue(dpId, origin, v)
	case commons.StringType:
		v, ok := value.(string)
		if !ok {
			return nil, err
		}
		return NewStringValue(dpId, origin, v)
	case commons.RawType:
		v, ok := value.([]byte)
		if !ok {
			return nil, err
		}
		return NewBinaryValue(dpId, origin, v)
	case commons.ValueType:
		v, ok := value.(int64)
		if !ok {
			return nil, err
		}
		return NewInt64Value(dpId, origin, v)
	case commons.EnumType:
		v, ok := value.(string)
		if !ok {
			return nil, err
		}
		return NewEnumValue(dpId, origin, v)
	case commons.FaultType, commons.BitmapType:
		v, ok := value.(int32)
		if !ok {
			return nil, err
		}
		return NewFaultValue(dpId, origin, v)
	default:
		return nil, err
	}
}

// NewBoolValue creates a DPValue of Type Bool with the given value.
func NewBoolValue(dpId string, origin int64, value bool) (*DPValue, error) {
	cv := &DPValue{DPId: dpId, Origin: origin, Type: commons.BoolType}
	err := encodeValue(cv, value)
	return cv, err
}

// NewStringValue creates a DPValue of Type string with the given value.
func NewStringValue(dpId string, origin int64, value string) (*DPValue, error) {
	return &DPValue{DPId: dpId, Origin: origin, Type: commons.StringType, stringValue: value}, nil
}

// NewInt64Value creates a DPValue of Type Int64 with the given value.
func NewInt64Value(dpId string, origin int64, value int64) (*DPValue, error) {
	cv := &DPValue{DPId: dpId, Origin: origin, Type: commons.ValueType}
	err := encodeValue(cv, value)
	return cv, err
}

func NewEnumValue(dpId string, origin int64, value string) (*DPValue, error) {
	return &DPValue{DPId: dpId, Origin: origin, Type: commons.EnumType, stringValue: value}, nil
}

func NewFaultValue(dpId string, origin int64, value int32) (*DPValue, error) {
	cv := &DPValue{DPId: dpId, Origin: origin, Type: commons.FaultType}
	err := encodeValue(cv, value)
	return cv, err
}

// NewBinaryValue creates a DPValue with binary payload and enforces the memory limit for event readings.
func NewBinaryValue(dpId string, origin int64, value []byte) (*DPValue, error) {
	return &DPValue{
		DPId:     dpId,
		Origin:   origin,
		Type:     commons.RawType,
		BinValue: value,
	}, nil
}

func encodeValue(cv *DPValue, value interface{}) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, value)
	if err == nil {
		cv.NumericValue = buf.Bytes()
	}
	return err
}

func decodeValue(reader io.Reader, value interface{}) error {
	err := binary.Read(reader, binary.BigEndian, value)
	return err
}

// ValueToString returns the string format of the value.
// In EdgeX, float value has two kinds of representation, Base64, and eNotation.
// Users can specify the floatEncoding in the properties value of the device profile, like floatEncoding: "Base64" or floatEncoding: "eNotation".
func (cv *DPValue) ValueToString() (string, error) {
	var str string
	if cv.Type == commons.StringType || cv.Type == commons.EnumType {
		str = cv.stringValue
		return str, nil
	}

	reader := bytes.NewReader(cv.NumericValue)

	switch cv.Type {
	case commons.BoolType:
		var res bool
		err := binary.Read(reader, binary.BigEndian, &res)
		if err != nil {
			return "", err
		}
		str = strconv.FormatBool(res)
	case commons.ValueType:
		var res int64
		err := binary.Read(reader, binary.BigEndian, &res)
		if err != nil {
			return "", err
		}
		str = strconv.FormatInt(res, 10)
	case commons.FaultType:
		var res int32
		err := binary.Read(reader, binary.BigEndian, &res)
		if err != nil {
			str = err.Error()
		}
		res = res & FaultMask
		str = strconv.FormatInt(int64(res), 10)
	default:
		return "", errors.New("nut supported type")
	}
	return str, nil
}

// String returns a string representation of a DPValue instance.
func (cv *DPValue) String() string {
	originStr := fmt.Sprintf("Origin: %d, ", cv.Origin)
	str, err := cv.ValueToString()
	if err != nil {
		return ""
	}
	valueStr := string(cv.Type) + ": " + str
	return originStr + valueStr
}

// BoolValue returns the value in bool data type, and returns error if the Type is not Bool.
func (cv *DPValue) BoolValue() (bool, error) {
	var value bool
	if cv.Type != commons.BoolType {
		return value, fmt.Errorf("the data type is not %T", value)
	}
	err := decodeValue(bytes.NewReader(cv.NumericValue), &value)
	return value, err
}

// StringValue returns the value in string data type, and returns error if the Type is not String.
func (cv *DPValue) StringValue() (string, error) {
	value := cv.stringValue
	if cv.Type != commons.StringType {
		return value, fmt.Errorf("the data type is not %T", value)
	}
	return value, nil
}

// Int64Value returns the value in int64 data type, and returns error if the Type is not Int64.
func (cv *DPValue) Int64Value() (int64, error) {
	var value int64
	if cv.Type != commons.ValueType {
		return value, fmt.Errorf("the data type is not %T", value)
	}
	err := decodeValue(bytes.NewReader(cv.NumericValue), &value)
	return value, err
}

// BinaryValue returns the value in []byte data type, and returns error if the Type is not Binary.
func (cv *DPValue) BinaryValue() ([]byte, error) {
	var value []byte
	if cv.Type != commons.RawType {
		return value, fmt.Errorf("the DPValue (%s) data type (%v) is not binary", cv.String(), cv.Type)
	}
	return cv.BinValue, nil
}

func (cv *DPValue) EnumValue() (string, error) {
	value := cv.stringValue
	if cv.Type != commons.EnumType {
		return value, fmt.Errorf("the data type is not %T", value)
	}
	return value, nil
}

func (cv *DPValue) FaultValue() (int32, error) {
	var value int32
	if cv.Type != commons.FaultType {
		return value, fmt.Errorf("the data type is not %T", value)
	}
	err := decodeValue(bytes.NewReader(cv.NumericValue), &value)
	return value, err
}

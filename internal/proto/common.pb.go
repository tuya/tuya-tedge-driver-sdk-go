// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: common.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EnumTUYAMQTTProtocol int32

const (
	EnumTUYAMQTTProtocol_ENUM_TUYAMQTT_PROTOCOL_UNSPECIFIED EnumTUYAMQTTProtocol = 0
	EnumTUYAMQTTProtocol_ENUM_TUYAMQTT_PROTOCOL_5           EnumTUYAMQTTProtocol = 5
	EnumTUYAMQTTProtocol_ENUM_TUYAMQTT_PROTOCOL_64          EnumTUYAMQTTProtocol = 64
	EnumTUYAMQTTProtocol_ENUM_TUYAMQTT_PROTOCOL_302         EnumTUYAMQTTProtocol = 302
	EnumTUYAMQTTProtocol_ENUM_TUYAMQTT_PROTOCOL_312         EnumTUYAMQTTProtocol = 312
)

// Enum value maps for EnumTUYAMQTTProtocol.
var (
	EnumTUYAMQTTProtocol_name = map[int32]string{
		0:   "ENUM_TUYAMQTT_PROTOCOL_UNSPECIFIED",
		5:   "ENUM_TUYAMQTT_PROTOCOL_5",
		64:  "ENUM_TUYAMQTT_PROTOCOL_64",
		302: "ENUM_TUYAMQTT_PROTOCOL_302",
		312: "ENUM_TUYAMQTT_PROTOCOL_312",
	}
	EnumTUYAMQTTProtocol_value = map[string]int32{
		"ENUM_TUYAMQTT_PROTOCOL_UNSPECIFIED": 0,
		"ENUM_TUYAMQTT_PROTOCOL_5":           5,
		"ENUM_TUYAMQTT_PROTOCOL_64":          64,
		"ENUM_TUYAMQTT_PROTOCOL_302":         302,
		"ENUM_TUYAMQTT_PROTOCOL_312":         312,
	}
)

func (x EnumTUYAMQTTProtocol) Enum() *EnumTUYAMQTTProtocol {
	p := new(EnumTUYAMQTTProtocol)
	*p = x
	return p
}

func (x EnumTUYAMQTTProtocol) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EnumTUYAMQTTProtocol) Descriptor() protoreflect.EnumDescriptor {
	return file_common_proto_enumTypes[0].Descriptor()
}

func (EnumTUYAMQTTProtocol) Type() protoreflect.EnumType {
	return &file_common_proto_enumTypes[0]
}

func (x EnumTUYAMQTTProtocol) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EnumTUYAMQTTProtocol.Descriptor instead.
func (EnumTUYAMQTTProtocol) EnumDescriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

// 日志级别
type EnumLogLevel int32

const (
	EnumLogLevel_ENUM_LOG_LEVEL_UNSPECIFIED EnumLogLevel = 0
	EnumLogLevel_ENUM_LOG_LEVEL_DEBUG       EnumLogLevel = 1
	EnumLogLevel_ENUM_LOG_LEVEL_INFO        EnumLogLevel = 2
	EnumLogLevel_ENUM_LOG_LEVEL_WARNING     EnumLogLevel = 3
	EnumLogLevel_ENUM_LOG_LEVEL_ERROR       EnumLogLevel = 4
)

// Enum value maps for EnumLogLevel.
var (
	EnumLogLevel_name = map[int32]string{
		0: "ENUM_LOG_LEVEL_UNSPECIFIED",
		1: "ENUM_LOG_LEVEL_DEBUG",
		2: "ENUM_LOG_LEVEL_INFO",
		3: "ENUM_LOG_LEVEL_WARNING",
		4: "ENUM_LOG_LEVEL_ERROR",
	}
	EnumLogLevel_value = map[string]int32{
		"ENUM_LOG_LEVEL_UNSPECIFIED": 0,
		"ENUM_LOG_LEVEL_DEBUG":       1,
		"ENUM_LOG_LEVEL_INFO":        2,
		"ENUM_LOG_LEVEL_WARNING":     3,
		"ENUM_LOG_LEVEL_ERROR":       4,
	}
)

func (x EnumLogLevel) Enum() *EnumLogLevel {
	p := new(EnumLogLevel)
	*p = x
	return p
}

func (x EnumLogLevel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EnumLogLevel) Descriptor() protoreflect.EnumDescriptor {
	return file_common_proto_enumTypes[1].Descriptor()
}

func (EnumLogLevel) Type() protoreflect.EnumType {
	return &file_common_proto_enumTypes[1]
}

func (x EnumLogLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EnumLogLevel.Descriptor instead.
func (EnumLogLevel) EnumDescriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{1}
}

type LogLevelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogLevel EnumLogLevel `protobuf:"varint,1,opt,name=log_level,json=logLevel,proto3,enum=common.EnumLogLevel" json:"log_level,omitempty"`
}

func (x *LogLevelRequest) Reset() {
	*x = LogLevelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogLevelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogLevelRequest) ProtoMessage() {}

func (x *LogLevelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogLevelRequest.ProtoReflect.Descriptor instead.
func (*LogLevelRequest) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

func (x *LogLevelRequest) GetLogLevel() EnumLogLevel {
	if x != nil {
		return x.LogLevel
	}
	return EnumLogLevel_ENUM_LOG_LEVEL_UNSPECIFIED
}

type PageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NameLike string `protobuf:"bytes,1,opt,name=name_like,json=nameLike,proto3" json:"name_like,omitempty"`
	Page     int64  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	PageSize int64  `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *PageRequest) Reset() {
	*x = PageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageRequest) ProtoMessage() {}

func (x *PageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageRequest.ProtoReflect.Descriptor instead.
func (*PageRequest) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{1}
}

func (x *PageRequest) GetNameLike() string {
	if x != nil {
		return x.NameLike
	}
	return ""
}

func (x *PageRequest) GetPage() int64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *PageRequest) GetPageSize() int64 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

//
type BaseWithIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id []string `protobuf:"bytes,1,rep,name=id,proto3" json:"id,omitempty"`
}

func (x *BaseWithIdResponse) Reset() {
	*x = BaseWithIdResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BaseWithIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseWithIdResponse) ProtoMessage() {}

func (x *BaseWithIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseWithIdResponse.ProtoReflect.Descriptor instead.
func (*BaseWithIdResponse) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{2}
}

func (x *BaseWithIdResponse) GetId() []string {
	if x != nil {
		return x.Id
	}
	return nil
}

//
type BaseExistResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exist bool `protobuf:"varint,1,opt,name=exist,proto3" json:"exist,omitempty"`
}

func (x *BaseExistResponse) Reset() {
	*x = BaseExistResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BaseExistResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseExistResponse) ProtoMessage() {}

func (x *BaseExistResponse) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseExistResponse.ProtoReflect.Descriptor instead.
func (*BaseExistResponse) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{3}
}

func (x *BaseExistResponse) GetExist() bool {
	if x != nil {
		return x.Exist
	}
	return false
}

// 条件查询
type BaseSearchConditionQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page     int32  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize int32  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Id       string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Ids      string `protobuf:"bytes,4,opt,name=ids,proto3" json:"ids,omitempty"`
	LikeId   string `protobuf:"bytes,5,opt,name=like_id,json=likeId,proto3" json:"like_id,omitempty"`
	Name     string `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	NameLike string `protobuf:"bytes,7,opt,name=name_like,json=nameLike,proto3" json:"name_like,omitempty"`
	IsAll    bool   `protobuf:"varint,8,opt,name=is_all,json=isAll,proto3" json:"is_all,omitempty"`
}

func (x *BaseSearchConditionQuery) Reset() {
	*x = BaseSearchConditionQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BaseSearchConditionQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseSearchConditionQuery) ProtoMessage() {}

func (x *BaseSearchConditionQuery) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseSearchConditionQuery.ProtoReflect.Descriptor instead.
func (*BaseSearchConditionQuery) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{4}
}

func (x *BaseSearchConditionQuery) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *BaseSearchConditionQuery) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *BaseSearchConditionQuery) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BaseSearchConditionQuery) GetIds() string {
	if x != nil {
		return x.Ids
	}
	return ""
}

func (x *BaseSearchConditionQuery) GetLikeId() string {
	if x != nil {
		return x.LikeId
	}
	return ""
}

func (x *BaseSearchConditionQuery) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BaseSearchConditionQuery) GetNameLike() string {
	if x != nil {
		return x.NameLike
	}
	return ""
}

func (x *BaseSearchConditionQuery) GetIsAll() bool {
	if x != nil {
		return x.IsAll
	}
	return false
}

// count
type CountResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count uint32 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *CountResponse) Reset() {
	*x = CountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CountResponse) ProtoMessage() {}

func (x *CountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CountResponse.ProtoReflect.Descriptor instead.
func (*CountResponse) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{5}
}

func (x *CountResponse) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

// pong
type Pong struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp string `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *Pong) Reset() {
	*x = Pong{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pong) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pong) ProtoMessage() {}

func (x *Pong) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pong.ProtoReflect.Descriptor instead.
func (*Pong) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{6}
}

func (x *Pong) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

// secret
type SecretDataKeyValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *SecretDataKeyValue) Reset() {
	*x = SecretDataKeyValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecretDataKeyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretDataKeyValue) ProtoMessage() {}

func (x *SecretDataKeyValue) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecretDataKeyValue.ProtoReflect.Descriptor instead.
func (*SecretDataKeyValue) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{7}
}

func (x *SecretDataKeyValue) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SecretDataKeyValue) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type SecretRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path       string                `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	SecretData []*SecretDataKeyValue `protobuf:"bytes,3,rep,name=secret_data,json=secretData,proto3" json:"secret_data,omitempty"`
}

func (x *SecretRequest) Reset() {
	*x = SecretRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretRequest) ProtoMessage() {}

func (x *SecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecretRequest.ProtoReflect.Descriptor instead.
func (*SecretRequest) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{8}
}

func (x *SecretRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *SecretRequest) GetSecretData() []*SecretDataKeyValue {
	if x != nil {
		return x.SecretData
	}
	return nil
}

// version
type VersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *VersionResponse) Reset() {
	*x = VersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionResponse) ProtoMessage() {}

func (x *VersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionResponse.ProtoReflect.Descriptor instead.
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{9}
}

func (x *VersionResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type VersionSdkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SdkVersion string `protobuf:"bytes,2,opt,name=sdk_version,json=sdkVersion,proto3" json:"sdk_version,omitempty"`
}

func (x *VersionSdkResponse) Reset() {
	*x = VersionSdkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionSdkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionSdkResponse) ProtoMessage() {}

func (x *VersionSdkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionSdkResponse.ProtoReflect.Descriptor instead.
func (*VersionSdkResponse) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{10}
}

func (x *VersionSdkResponse) GetSdkVersion() string {
	if x != nil {
		return x.SdkVersion
	}
	return ""
}

var File_common_proto protoreflect.FileDescriptor

var file_common_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x44, 0x0a, 0x0f, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x31, 0x0a, 0x09, 0x6c, 0x6f, 0x67, 0x5f, 0x6c, 0x65,
	0x76, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x52,
	0x08, 0x6c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x22, 0x5b, 0x0a, 0x0b, 0x50, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65,
	0x5f, 0x6c, 0x69, 0x6b, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x61, 0x6d,
	0x65, 0x4c, 0x69, 0x6b, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67,
	0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x61,
	0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x24, 0x0a, 0x12, 0x42, 0x61, 0x73, 0x65, 0x57, 0x69,
	0x74, 0x68, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x29, 0x0a, 0x11,
	0x42, 0x61, 0x73, 0x65, 0x45, 0x78, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x78, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x05, 0x65, 0x78, 0x69, 0x73, 0x74, 0x22, 0xce, 0x01, 0x0a, 0x18, 0x42, 0x61, 0x73, 0x65,
	0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65,
	0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x6c, 0x69, 0x6b, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x69, 0x6b, 0x65, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x6c, 0x69, 0x6b,
	0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x61, 0x6d, 0x65, 0x4c, 0x69, 0x6b,
	0x65, 0x12, 0x15, 0x0a, 0x06, 0x69, 0x73, 0x5f, 0x61, 0x6c, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x05, 0x69, 0x73, 0x41, 0x6c, 0x6c, 0x22, 0x25, 0x0a, 0x0d, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x24, 0x0a, 0x04, 0x50, 0x6f, 0x6e, 0x67, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x3c, 0x0a, 0x12, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x44,
	0x61, 0x74, 0x61, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x60, 0x0a, 0x0d, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x3b, 0x0a, 0x0b, 0x73, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x73, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x22, 0x2b, 0x0a, 0x0f, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x22, 0x35, 0x0a, 0x12, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x64, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x64, 0x6b, 0x5f,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73,
	0x64, 0x6b, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2a, 0xbd, 0x01, 0x0a, 0x14, 0x45, 0x6e,
	0x75, 0x6d, 0x54, 0x55, 0x59, 0x41, 0x4d, 0x51, 0x54, 0x54, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x12, 0x26, 0x0a, 0x22, 0x45, 0x4e, 0x55, 0x4d, 0x5f, 0x54, 0x55, 0x59, 0x41, 0x4d,
	0x51, 0x54, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4f, 0x4c, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1c, 0x0a, 0x18, 0x45, 0x4e,
	0x55, 0x4d, 0x5f, 0x54, 0x55, 0x59, 0x41, 0x4d, 0x51, 0x54, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54,
	0x4f, 0x43, 0x4f, 0x4c, 0x5f, 0x35, 0x10, 0x05, 0x12, 0x1d, 0x0a, 0x19, 0x45, 0x4e, 0x55, 0x4d,
	0x5f, 0x54, 0x55, 0x59, 0x41, 0x4d, 0x51, 0x54, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4f, 0x4c, 0x5f, 0x36, 0x34, 0x10, 0x40, 0x12, 0x1f, 0x0a, 0x1a, 0x45, 0x4e, 0x55, 0x4d, 0x5f,
	0x54, 0x55, 0x59, 0x41, 0x4d, 0x51, 0x54, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4f,
	0x4c, 0x5f, 0x33, 0x30, 0x32, 0x10, 0xae, 0x02, 0x12, 0x1f, 0x0a, 0x1a, 0x45, 0x4e, 0x55, 0x4d,
	0x5f, 0x54, 0x55, 0x59, 0x41, 0x4d, 0x51, 0x54, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4f, 0x4c, 0x5f, 0x33, 0x31, 0x32, 0x10, 0xb8, 0x02, 0x2a, 0x97, 0x01, 0x0a, 0x0c, 0x45, 0x6e,
	0x75, 0x6d, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x1e, 0x0a, 0x1a, 0x45, 0x4e,
	0x55, 0x4d, 0x5f, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x45, 0x4e,
	0x55, 0x4d, 0x5f, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x44, 0x45, 0x42,
	0x55, 0x47, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x4e, 0x55, 0x4d, 0x5f, 0x4c, 0x4f, 0x47,
	0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x02, 0x12, 0x1a, 0x0a,
	0x16, 0x45, 0x4e, 0x55, 0x4d, 0x5f, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f,
	0x57, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x12, 0x18, 0x0a, 0x14, 0x45, 0x4e, 0x55,
	0x4d, 0x5f, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x10, 0x04, 0x32, 0x76, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x12, 0x2e, 0x0a,
	0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x22, 0x00, 0x12, 0x3c, 0x0a,
	0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x17, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e,
	0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_proto_rawDescOnce sync.Once
	file_common_proto_rawDescData = file_common_proto_rawDesc
)

func file_common_proto_rawDescGZIP() []byte {
	file_common_proto_rawDescOnce.Do(func() {
		file_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_proto_rawDescData)
	})
	return file_common_proto_rawDescData
}

var file_common_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_common_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_common_proto_goTypes = []interface{}{
	(EnumTUYAMQTTProtocol)(0),        // 0: common.EnumTUYAMQTTProtocol
	(EnumLogLevel)(0),                // 1: common.EnumLogLevel
	(*LogLevelRequest)(nil),          // 2: common.LogLevelRequest
	(*PageRequest)(nil),              // 3: common.PageRequest
	(*BaseWithIdResponse)(nil),       // 4: common.BaseWithIdResponse
	(*BaseExistResponse)(nil),        // 5: common.BaseExistResponse
	(*BaseSearchConditionQuery)(nil), // 6: common.BaseSearchConditionQuery
	(*CountResponse)(nil),            // 7: common.CountResponse
	(*Pong)(nil),                     // 8: common.Pong
	(*SecretDataKeyValue)(nil),       // 9: common.SecretDataKeyValue
	(*SecretRequest)(nil),            // 10: common.SecretRequest
	(*VersionResponse)(nil),          // 11: common.VersionResponse
	(*VersionSdkResponse)(nil),       // 12: common.VersionSdkResponse
	(*emptypb.Empty)(nil),            // 13: google.protobuf.Empty
}
var file_common_proto_depIdxs = []int32{
	1,  // 0: common.LogLevelRequest.log_level:type_name -> common.EnumLogLevel
	9,  // 1: common.SecretRequest.secret_data:type_name -> common.SecretDataKeyValue
	13, // 2: common.Common.Ping:input_type -> google.protobuf.Empty
	13, // 3: common.Common.Version:input_type -> google.protobuf.Empty
	8,  // 4: common.Common.Ping:output_type -> common.Pong
	11, // 5: common.Common.Version:output_type -> common.VersionResponse
	4,  // [4:6] is the sub-list for method output_type
	2,  // [2:4] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_common_proto_init() }
func file_common_proto_init() {
	if File_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogLevelRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PageRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BaseWithIdResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BaseExistResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BaseSearchConditionQuery); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CountResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pong); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SecretDataKeyValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SecretRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionSdkResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_common_proto_goTypes,
		DependencyIndexes: file_common_proto_depIdxs,
		EnumInfos:         file_common_proto_enumTypes,
		MessageInfos:      file_common_proto_msgTypes,
	}.Build()
	File_common_proto = out.File
	file_common_proto_rawDesc = nil
	file_common_proto_goTypes = nil
	file_common_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CommonClient is the client API for Common service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CommonClient interface {
	// Ping tests whether the service is working
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Pong, error)
	// Version obtains version information from the target service.
	Version(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*VersionResponse, error)
}

type commonClient struct {
	cc grpc.ClientConnInterface
}

func NewCommonClient(cc grpc.ClientConnInterface) CommonClient {
	return &commonClient{cc}
}

func (c *commonClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Pong, error) {
	out := new(Pong)
	err := c.cc.Invoke(ctx, "/common.Common/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commonClient) Version(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, "/common.Common/Version", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommonServer is the server API for Common service.
type CommonServer interface {
	// Ping tests whether the service is working
	Ping(context.Context, *emptypb.Empty) (*Pong, error)
	// Version obtains version information from the target service.
	Version(context.Context, *emptypb.Empty) (*VersionResponse, error)
}

// UnimplementedCommonServer can be embedded to have forward compatible implementations.
type UnimplementedCommonServer struct {
}

func (*UnimplementedCommonServer) Ping(context.Context, *emptypb.Empty) (*Pong, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedCommonServer) Version(context.Context, *emptypb.Empty) (*VersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Version not implemented")
}

func RegisterCommonServer(s *grpc.Server, srv CommonServer) {
	s.RegisterService(&_Common_serviceDesc, srv)
}

func _Common_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommonServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.Common/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommonServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Common_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommonServer).Version(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.Common/Version",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommonServer).Version(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Common_serviceDesc = grpc.ServiceDesc{
	ServiceName: "common.Common",
	HandlerType: (*CommonServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Common_Ping_Handler,
		},
		{
			MethodName: "Version",
			Handler:    _Common_Version_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "common.proto",
}

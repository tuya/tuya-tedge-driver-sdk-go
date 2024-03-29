// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: driverstorage.proto

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

type PutReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DriverServiceId string `protobuf:"bytes,1,opt,name=driver_service_id,json=driverServiceId,proto3" json:"driver_service_id,omitempty"`
	Data            []*KV  `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *PutReq) Reset() {
	*x = PutReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_driverstorage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutReq) ProtoMessage() {}

func (x *PutReq) ProtoReflect() protoreflect.Message {
	mi := &file_driverstorage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutReq.ProtoReflect.Descriptor instead.
func (*PutReq) Descriptor() ([]byte, []int) {
	return file_driverstorage_proto_rawDescGZIP(), []int{0}
}

func (x *PutReq) GetDriverServiceId() string {
	if x != nil {
		return x.DriverServiceId
	}
	return ""
}

func (x *PutReq) GetData() []*KV {
	if x != nil {
		return x.Data
	}
	return nil
}

type GetReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DriverServiceId string   `protobuf:"bytes,1,opt,name=driver_service_id,json=driverServiceId,proto3" json:"driver_service_id,omitempty"`
	Keys            []string `protobuf:"bytes,2,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *GetReq) Reset() {
	*x = GetReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_driverstorage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReq) ProtoMessage() {}

func (x *GetReq) ProtoReflect() protoreflect.Message {
	mi := &file_driverstorage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetReq.ProtoReflect.Descriptor instead.
func (*GetReq) Descriptor() ([]byte, []int) {
	return file_driverstorage_proto_rawDescGZIP(), []int{1}
}

func (x *GetReq) GetDriverServiceId() string {
	if x != nil {
		return x.DriverServiceId
	}
	return ""
}

func (x *GetReq) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type AllReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DriverServiceId string `protobuf:"bytes,1,opt,name=driver_service_id,json=driverServiceId,proto3" json:"driver_service_id,omitempty"`
}

func (x *AllReq) Reset() {
	*x = AllReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_driverstorage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllReq) ProtoMessage() {}

func (x *AllReq) ProtoReflect() protoreflect.Message {
	mi := &file_driverstorage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllReq.ProtoReflect.Descriptor instead.
func (*AllReq) Descriptor() ([]byte, []int) {
	return file_driverstorage_proto_rawDescGZIP(), []int{2}
}

func (x *AllReq) GetDriverServiceId() string {
	if x != nil {
		return x.DriverServiceId
	}
	return ""
}

type DeleteReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DriverServiceId string   `protobuf:"bytes,1,opt,name=driver_service_id,json=driverServiceId,proto3" json:"driver_service_id,omitempty"`
	Keys            []string `protobuf:"bytes,2,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *DeleteReq) Reset() {
	*x = DeleteReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_driverstorage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteReq) ProtoMessage() {}

func (x *DeleteReq) ProtoReflect() protoreflect.Message {
	mi := &file_driverstorage_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteReq.ProtoReflect.Descriptor instead.
func (*DeleteReq) Descriptor() ([]byte, []int) {
	return file_driverstorage_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteReq) GetDriverServiceId() string {
	if x != nil {
		return x.DriverServiceId
	}
	return ""
}

func (x *DeleteReq) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type KVs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kvs []*KV `protobuf:"bytes,1,rep,name=kvs,proto3" json:"kvs,omitempty"`
}

func (x *KVs) Reset() {
	*x = KVs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_driverstorage_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KVs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KVs) ProtoMessage() {}

func (x *KVs) ProtoReflect() protoreflect.Message {
	mi := &file_driverstorage_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KVs.ProtoReflect.Descriptor instead.
func (*KVs) Descriptor() ([]byte, []int) {
	return file_driverstorage_proto_rawDescGZIP(), []int{4}
}

func (x *KVs) GetKvs() []*KV {
	if x != nil {
		return x.Kvs
	}
	return nil
}

type KV struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *KV) Reset() {
	*x = KV{}
	if protoimpl.UnsafeEnabled {
		mi := &file_driverstorage_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KV) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KV) ProtoMessage() {}

func (x *KV) ProtoReflect() protoreflect.Message {
	mi := &file_driverstorage_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KV.ProtoReflect.Descriptor instead.
func (*KV) Descriptor() ([]byte, []int) {
	return file_driverstorage_proto_rawDescGZIP(), []int{5}
}

func (x *KV) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *KV) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type Keys struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key []string `protobuf:"bytes,1,rep,name=key,proto3" json:"key,omitempty"`
}

func (x *Keys) Reset() {
	*x = Keys{}
	if protoimpl.UnsafeEnabled {
		mi := &file_driverstorage_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Keys) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Keys) ProtoMessage() {}

func (x *Keys) ProtoReflect() protoreflect.Message {
	mi := &file_driverstorage_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Keys.ProtoReflect.Descriptor instead.
func (*Keys) Descriptor() ([]byte, []int) {
	return file_driverstorage_proto_rawDescGZIP(), []int{6}
}

func (x *Keys) GetKey() []string {
	if x != nil {
		return x.Key
	}
	return nil
}

type GetPrefixReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DriverServiceId string `protobuf:"bytes,1,opt,name=driver_service_id,json=driverServiceId,proto3" json:"driver_service_id,omitempty"`
	Prefix          string `protobuf:"bytes,2,opt,name=prefix,proto3" json:"prefix,omitempty"` // 根据前缀获取
}

func (x *GetPrefixReq) Reset() {
	*x = GetPrefixReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_driverstorage_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPrefixReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPrefixReq) ProtoMessage() {}

func (x *GetPrefixReq) ProtoReflect() protoreflect.Message {
	mi := &file_driverstorage_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPrefixReq.ProtoReflect.Descriptor instead.
func (*GetPrefixReq) Descriptor() ([]byte, []int) {
	return file_driverstorage_proto_rawDescGZIP(), []int{7}
}

func (x *GetPrefixReq) GetDriverServiceId() string {
	if x != nil {
		return x.DriverServiceId
	}
	return ""
}

func (x *GetPrefixReq) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

var File_driverstorage_proto protoreflect.FileDescriptor

var file_driverstorage_proto_rawDesc = []byte{
	0x0a, 0x13, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x5b, 0x0a, 0x06, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x12, 0x2a, 0x0a, 0x11, 0x64,
	0x72, 0x69, 0x76, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x4b, 0x56, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x48,
	0x0a, 0x06, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x12, 0x2a, 0x0a, 0x11, 0x64, 0x72, 0x69, 0x76,
	0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x22, 0x34, 0x0a, 0x06, 0x41, 0x6c, 0x6c, 0x52,
	0x65, 0x71, 0x12, 0x2a, 0x0a, 0x11, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64,
	0x72, 0x69, 0x76, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x22, 0x4b,
	0x0a, 0x09, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x12, 0x2a, 0x0a, 0x11, 0x64,
	0x72, 0x69, 0x76, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x22, 0x2a, 0x0a, 0x03, 0x4b,
	0x56, 0x73, 0x12, 0x23, 0x0a, 0x03, 0x6b, 0x76, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x11, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e,
	0x4b, 0x56, 0x52, 0x03, 0x6b, 0x76, 0x73, 0x22, 0x2c, 0x0a, 0x02, 0x4b, 0x56, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x18, 0x0a, 0x04, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22,
	0x52, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x52, 0x65, 0x71, 0x12,
	0x2a, 0x0a, 0x11, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x72, 0x69, 0x76,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x70,
	0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72, 0x65,
	0x66, 0x69, 0x78, 0x32, 0x8b, 0x04, 0x0a, 0x0d, 0x44, 0x72, 0x69, 0x76, 0x65, 0x72, 0x53, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x30, 0x0a, 0x03, 0x41, 0x6c, 0x6c, 0x12, 0x15, 0x2e, 0x64,
	0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x41, 0x6c, 0x6c,
	0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x2e, 0x4b, 0x56, 0x73, 0x12, 0x3b, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x4b, 0x65,
	0x79, 0x73, 0x12, 0x1b, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x52, 0x65, 0x71, 0x1a,
	0x13, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e,
	0x4b, 0x65, 0x79, 0x73, 0x12, 0x30, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x15, 0x2e, 0x64, 0x72,
	0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x1a, 0x12, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2e, 0x4b, 0x56, 0x73, 0x12, 0x34, 0x0a, 0x03, 0x50, 0x75, 0x74, 0x12, 0x15, 0x2e,
	0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x75,
	0x74, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3a, 0x0a, 0x06,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x32, 0x0a, 0x05, 0x47, 0x65, 0x74, 0x56,
	0x32, 0x12, 0x15, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65,
	0x72, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x4b, 0x56, 0x73, 0x12, 0x3d, 0x0a, 0x09,
	0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x73, 0x56, 0x32, 0x12, 0x1b, 0x2e, 0x64, 0x72, 0x69, 0x76,
	0x65, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65,
	0x66, 0x69, 0x78, 0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x36, 0x0a, 0x05, 0x50,
	0x75, 0x74, 0x56, 0x32, 0x12, 0x15, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2e, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x12, 0x3c, 0x0a, 0x08, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x56, 0x32, 0x12,
	0x18, 0x2e, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_driverstorage_proto_rawDescOnce sync.Once
	file_driverstorage_proto_rawDescData = file_driverstorage_proto_rawDesc
)

func file_driverstorage_proto_rawDescGZIP() []byte {
	file_driverstorage_proto_rawDescOnce.Do(func() {
		file_driverstorage_proto_rawDescData = protoimpl.X.CompressGZIP(file_driverstorage_proto_rawDescData)
	})
	return file_driverstorage_proto_rawDescData
}

var file_driverstorage_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_driverstorage_proto_goTypes = []interface{}{
	(*PutReq)(nil),        // 0: driverstorage.PutReq
	(*GetReq)(nil),        // 1: driverstorage.GetReq
	(*AllReq)(nil),        // 2: driverstorage.AllReq
	(*DeleteReq)(nil),     // 3: driverstorage.DeleteReq
	(*KVs)(nil),           // 4: driverstorage.KVs
	(*KV)(nil),            // 5: driverstorage.KV
	(*Keys)(nil),          // 6: driverstorage.Keys
	(*GetPrefixReq)(nil),  // 7: driverstorage.GetPrefixReq
	(*emptypb.Empty)(nil), // 8: google.protobuf.Empty
}
var file_driverstorage_proto_depIdxs = []int32{
	5,  // 0: driverstorage.PutReq.data:type_name -> driverstorage.KV
	5,  // 1: driverstorage.KVs.kvs:type_name -> driverstorage.KV
	2,  // 2: driverstorage.DriverStorage.All:input_type -> driverstorage.AllReq
	7,  // 3: driverstorage.DriverStorage.GetKeys:input_type -> driverstorage.GetPrefixReq
	1,  // 4: driverstorage.DriverStorage.Get:input_type -> driverstorage.GetReq
	0,  // 5: driverstorage.DriverStorage.Put:input_type -> driverstorage.PutReq
	3,  // 6: driverstorage.DriverStorage.Delete:input_type -> driverstorage.DeleteReq
	1,  // 7: driverstorage.DriverStorage.GetV2:input_type -> driverstorage.GetReq
	7,  // 8: driverstorage.DriverStorage.GetKeysV2:input_type -> driverstorage.GetPrefixReq
	0,  // 9: driverstorage.DriverStorage.PutV2:input_type -> driverstorage.PutReq
	3,  // 10: driverstorage.DriverStorage.DeleteV2:input_type -> driverstorage.DeleteReq
	4,  // 11: driverstorage.DriverStorage.All:output_type -> driverstorage.KVs
	6,  // 12: driverstorage.DriverStorage.GetKeys:output_type -> driverstorage.Keys
	4,  // 13: driverstorage.DriverStorage.Get:output_type -> driverstorage.KVs
	8,  // 14: driverstorage.DriverStorage.Put:output_type -> google.protobuf.Empty
	8,  // 15: driverstorage.DriverStorage.Delete:output_type -> google.protobuf.Empty
	4,  // 16: driverstorage.DriverStorage.GetV2:output_type -> driverstorage.KVs
	6,  // 17: driverstorage.DriverStorage.GetKeysV2:output_type -> driverstorage.Keys
	8,  // 18: driverstorage.DriverStorage.PutV2:output_type -> google.protobuf.Empty
	8,  // 19: driverstorage.DriverStorage.DeleteV2:output_type -> google.protobuf.Empty
	11, // [11:20] is the sub-list for method output_type
	2,  // [2:11] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_driverstorage_proto_init() }
func file_driverstorage_proto_init() {
	if File_driverstorage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_driverstorage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutReq); i {
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
		file_driverstorage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetReq); i {
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
		file_driverstorage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllReq); i {
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
		file_driverstorage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteReq); i {
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
		file_driverstorage_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KVs); i {
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
		file_driverstorage_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KV); i {
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
		file_driverstorage_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Keys); i {
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
		file_driverstorage_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPrefixReq); i {
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
			RawDescriptor: file_driverstorage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_driverstorage_proto_goTypes,
		DependencyIndexes: file_driverstorage_proto_depIdxs,
		MessageInfos:      file_driverstorage_proto_msgTypes,
	}.Build()
	File_driverstorage_proto = out.File
	file_driverstorage_proto_rawDesc = nil
	file_driverstorage_proto_goTypes = nil
	file_driverstorage_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// DriverStorageClient is the client API for DriverStorage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DriverStorageClient interface {
	// All 获取全部本地数据(不带备份数据)
	// deprecated 存在性能问题
	All(ctx context.Context, in *AllReq, opts ...grpc.CallOption) (*KVs, error)
	// GetKeys 获取非备份的key
	GetKeys(ctx context.Context, in *GetPrefixReq, opts ...grpc.CallOption) (*Keys, error)
	// Get 本地数据www
	Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*KVs, error)
	// Put 存储在本地
	Put(ctx context.Context, in *PutReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Delete 删除本地数据
	Delete(ctx context.Context, in *DeleteReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// GetV2 获取备份数据
	GetV2(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*KVs, error)
	// GetKeysV2 只获取备份数据的Key
	GetKeysV2(ctx context.Context, in *GetPrefixReq, opts ...grpc.CallOption) (*Keys, error)
	// PutV2 存储备份数据,限制大小，超过会报错
	PutV2(ctx context.Context, in *PutReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeleteV2 删除备份数据
	DeleteV2(ctx context.Context, in *DeleteReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type driverStorageClient struct {
	cc grpc.ClientConnInterface
}

func NewDriverStorageClient(cc grpc.ClientConnInterface) DriverStorageClient {
	return &driverStorageClient{cc}
}

func (c *driverStorageClient) All(ctx context.Context, in *AllReq, opts ...grpc.CallOption) (*KVs, error) {
	out := new(KVs)
	err := c.cc.Invoke(ctx, "/driverstorage.DriverStorage/All", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverStorageClient) GetKeys(ctx context.Context, in *GetPrefixReq, opts ...grpc.CallOption) (*Keys, error) {
	out := new(Keys)
	err := c.cc.Invoke(ctx, "/driverstorage.DriverStorage/GetKeys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverStorageClient) Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*KVs, error) {
	out := new(KVs)
	err := c.cc.Invoke(ctx, "/driverstorage.DriverStorage/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverStorageClient) Put(ctx context.Context, in *PutReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/driverstorage.DriverStorage/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverStorageClient) Delete(ctx context.Context, in *DeleteReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/driverstorage.DriverStorage/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverStorageClient) GetV2(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*KVs, error) {
	out := new(KVs)
	err := c.cc.Invoke(ctx, "/driverstorage.DriverStorage/GetV2", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverStorageClient) GetKeysV2(ctx context.Context, in *GetPrefixReq, opts ...grpc.CallOption) (*Keys, error) {
	out := new(Keys)
	err := c.cc.Invoke(ctx, "/driverstorage.DriverStorage/GetKeysV2", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverStorageClient) PutV2(ctx context.Context, in *PutReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/driverstorage.DriverStorage/PutV2", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *driverStorageClient) DeleteV2(ctx context.Context, in *DeleteReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/driverstorage.DriverStorage/DeleteV2", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DriverStorageServer is the server API for DriverStorage service.
type DriverStorageServer interface {
	// All 获取全部本地数据(不带备份数据)
	// deprecated 存在性能问题
	All(context.Context, *AllReq) (*KVs, error)
	// GetKeys 获取非备份的key
	GetKeys(context.Context, *GetPrefixReq) (*Keys, error)
	// Get 本地数据www
	Get(context.Context, *GetReq) (*KVs, error)
	// Put 存储在本地
	Put(context.Context, *PutReq) (*emptypb.Empty, error)
	// Delete 删除本地数据
	Delete(context.Context, *DeleteReq) (*emptypb.Empty, error)
	// GetV2 获取备份数据
	GetV2(context.Context, *GetReq) (*KVs, error)
	// GetKeysV2 只获取备份数据的Key
	GetKeysV2(context.Context, *GetPrefixReq) (*Keys, error)
	// PutV2 存储备份数据,限制大小，超过会报错
	PutV2(context.Context, *PutReq) (*emptypb.Empty, error)
	// DeleteV2 删除备份数据
	DeleteV2(context.Context, *DeleteReq) (*emptypb.Empty, error)
}

// UnimplementedDriverStorageServer can be embedded to have forward compatible implementations.
type UnimplementedDriverStorageServer struct {
}

func (*UnimplementedDriverStorageServer) All(context.Context, *AllReq) (*KVs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method All not implemented")
}
func (*UnimplementedDriverStorageServer) GetKeys(context.Context, *GetPrefixReq) (*Keys, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeys not implemented")
}
func (*UnimplementedDriverStorageServer) Get(context.Context, *GetReq) (*KVs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedDriverStorageServer) Put(context.Context, *PutReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (*UnimplementedDriverStorageServer) Delete(context.Context, *DeleteReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (*UnimplementedDriverStorageServer) GetV2(context.Context, *GetReq) (*KVs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetV2 not implemented")
}
func (*UnimplementedDriverStorageServer) GetKeysV2(context.Context, *GetPrefixReq) (*Keys, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeysV2 not implemented")
}
func (*UnimplementedDriverStorageServer) PutV2(context.Context, *PutReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutV2 not implemented")
}
func (*UnimplementedDriverStorageServer) DeleteV2(context.Context, *DeleteReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteV2 not implemented")
}

func RegisterDriverStorageServer(s *grpc.Server, srv DriverStorageServer) {
	s.RegisterService(&_DriverStorage_serviceDesc, srv)
}

func _DriverStorage_All_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverStorageServer).All(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverstorage.DriverStorage/All",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverStorageServer).All(ctx, req.(*AllReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverStorage_GetKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPrefixReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverStorageServer).GetKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverstorage.DriverStorage/GetKeys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverStorageServer).GetKeys(ctx, req.(*GetPrefixReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverStorage_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverStorageServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverstorage.DriverStorage/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverStorageServer).Get(ctx, req.(*GetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverStorage_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverStorageServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverstorage.DriverStorage/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverStorageServer).Put(ctx, req.(*PutReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverStorage_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverStorageServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverstorage.DriverStorage/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverStorageServer).Delete(ctx, req.(*DeleteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverStorage_GetV2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverStorageServer).GetV2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverstorage.DriverStorage/GetV2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverStorageServer).GetV2(ctx, req.(*GetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverStorage_GetKeysV2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPrefixReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverStorageServer).GetKeysV2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverstorage.DriverStorage/GetKeysV2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverStorageServer).GetKeysV2(ctx, req.(*GetPrefixReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverStorage_PutV2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverStorageServer).PutV2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverstorage.DriverStorage/PutV2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverStorageServer).PutV2(ctx, req.(*PutReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DriverStorage_DeleteV2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DriverStorageServer).DeleteV2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/driverstorage.DriverStorage/DeleteV2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DriverStorageServer).DeleteV2(ctx, req.(*DeleteReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _DriverStorage_serviceDesc = grpc.ServiceDesc{
	ServiceName: "driverstorage.DriverStorage",
	HandlerType: (*DriverStorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "All",
			Handler:    _DriverStorage_All_Handler,
		},
		{
			MethodName: "GetKeys",
			Handler:    _DriverStorage_GetKeys_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _DriverStorage_Get_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _DriverStorage_Put_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _DriverStorage_Delete_Handler,
		},
		{
			MethodName: "GetV2",
			Handler:    _DriverStorage_GetV2_Handler,
		},
		{
			MethodName: "GetKeysV2",
			Handler:    _DriverStorage_GetKeysV2_Handler,
		},
		{
			MethodName: "PutV2",
			Handler:    _DriverStorage_PutV2_Handler,
		},
		{
			MethodName: "DeleteV2",
			Handler:    _DriverStorage_DeleteV2_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "driverstorage.proto",
}

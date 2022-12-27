// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: gateway.proto

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

// 网关信息
type GateWayInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Env         string `protobuf:"bytes,2,opt,name=env,proto3" json:"env,omitempty"`
	GwId        string `protobuf:"bytes,3,opt,name=gw_id,json=gwId,proto3" json:"gw_id,omitempty"`
	LocalKey    string `protobuf:"bytes,4,opt,name=local_key,json=localKey,proto3" json:"local_key,omitempty"`
	Region      string `protobuf:"bytes,6,opt,name=region,proto3" json:"region,omitempty"`
	IsNewModel  bool   `protobuf:"varint,15,opt,name=is_new_model,json=isNewModel,proto3" json:"is_new_model,omitempty"` // thing model
	CloudState  bool   `protobuf:"varint,17,opt,name=cloud_state,json=cloudState,proto3" json:"cloud_state,omitempty"`
	GatewayName string `protobuf:"bytes,19,opt,name=gateway_name,json=gatewayName,proto3" json:"gateway_name,omitempty"`
}

func (x *GateWayInfoResponse) Reset() {
	*x = GateWayInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GateWayInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GateWayInfoResponse) ProtoMessage() {}

func (x *GateWayInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GateWayInfoResponse.ProtoReflect.Descriptor instead.
func (*GateWayInfoResponse) Descriptor() ([]byte, []int) {
	return file_gateway_proto_rawDescGZIP(), []int{0}
}

func (x *GateWayInfoResponse) GetEnv() string {
	if x != nil {
		return x.Env
	}
	return ""
}

func (x *GateWayInfoResponse) GetGwId() string {
	if x != nil {
		return x.GwId
	}
	return ""
}

func (x *GateWayInfoResponse) GetLocalKey() string {
	if x != nil {
		return x.LocalKey
	}
	return ""
}

func (x *GateWayInfoResponse) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *GateWayInfoResponse) GetIsNewModel() bool {
	if x != nil {
		return x.IsNewModel
	}
	return false
}

func (x *GateWayInfoResponse) GetCloudState() bool {
	if x != nil {
		return x.CloudState
	}
	return false
}

func (x *GateWayInfoResponse) GetGatewayName() string {
	if x != nil {
		return x.GatewayName
	}
	return ""
}

var File_gateway_proto protoreflect.FileDescriptor

var file_gateway_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd7, 0x01, 0x0a,
	0x13, 0x47, 0x61, 0x74, 0x65, 0x57, 0x61, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x65, 0x6e, 0x76, 0x12, 0x13, 0x0a, 0x05, 0x67, 0x77, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x67, 0x77, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6c,
	0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x4b, 0x65, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69,
	0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x12, 0x20, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x6e, 0x65, 0x77, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x18, 0x0f, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x4e, 0x65, 0x77, 0x4d, 0x6f, 0x64,
	0x65, 0x6c, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x5f, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x13, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x32, 0x57, 0x0a, 0x13, 0x52, 0x70, 0x63, 0x47, 0x61, 0x74,
	0x65, 0x77, 0x61, 0x79, 0x46, 0x6f, 0x72, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x40, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x14, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x57, 0x61,
	0x79, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x0a, 0x5a, 0x08, 0x2e, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_gateway_proto_rawDescOnce sync.Once
	file_gateway_proto_rawDescData = file_gateway_proto_rawDesc
)

func file_gateway_proto_rawDescGZIP() []byte {
	file_gateway_proto_rawDescOnce.Do(func() {
		file_gateway_proto_rawDescData = protoimpl.X.CompressGZIP(file_gateway_proto_rawDescData)
	})
	return file_gateway_proto_rawDescData
}

var file_gateway_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_gateway_proto_goTypes = []interface{}{
	(*GateWayInfoResponse)(nil), // 0: GateWayInfoResponse
	(*emptypb.Empty)(nil),       // 1: google.protobuf.Empty
}
var file_gateway_proto_depIdxs = []int32{
	1, // 0: RpcGatewayForDevice.GetGatewayInfo:input_type -> google.protobuf.Empty
	0, // 1: RpcGatewayForDevice.GetGatewayInfo:output_type -> GateWayInfoResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gateway_proto_init() }
func file_gateway_proto_init() {
	if File_gateway_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gateway_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GateWayInfoResponse); i {
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
			RawDescriptor: file_gateway_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gateway_proto_goTypes,
		DependencyIndexes: file_gateway_proto_depIdxs,
		MessageInfos:      file_gateway_proto_msgTypes,
	}.Build()
	File_gateway_proto = out.File
	file_gateway_proto_rawDesc = nil
	file_gateway_proto_goTypes = nil
	file_gateway_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RpcGatewayForDeviceClient is the client API for RpcGatewayForDevice service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RpcGatewayForDeviceClient interface {
	// 获取网关信息
	GetGatewayInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GateWayInfoResponse, error)
}

type rpcGatewayForDeviceClient struct {
	cc grpc.ClientConnInterface
}

func NewRpcGatewayForDeviceClient(cc grpc.ClientConnInterface) RpcGatewayForDeviceClient {
	return &rpcGatewayForDeviceClient{cc}
}

func (c *rpcGatewayForDeviceClient) GetGatewayInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GateWayInfoResponse, error) {
	out := new(GateWayInfoResponse)
	err := c.cc.Invoke(ctx, "/RpcGatewayForDevice/GetGatewayInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RpcGatewayForDeviceServer is the server API for RpcGatewayForDevice service.
type RpcGatewayForDeviceServer interface {
	// 获取网关信息
	GetGatewayInfo(context.Context, *emptypb.Empty) (*GateWayInfoResponse, error)
}

// UnimplementedRpcGatewayForDeviceServer can be embedded to have forward compatible implementations.
type UnimplementedRpcGatewayForDeviceServer struct {
}

func (*UnimplementedRpcGatewayForDeviceServer) GetGatewayInfo(context.Context, *emptypb.Empty) (*GateWayInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGatewayInfo not implemented")
}

func RegisterRpcGatewayForDeviceServer(s *grpc.Server, srv RpcGatewayForDeviceServer) {
	s.RegisterService(&_RpcGatewayForDevice_serviceDesc, srv)
}

func _RpcGatewayForDevice_GetGatewayInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcGatewayForDeviceServer).GetGatewayInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcGatewayForDevice/GetGatewayInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcGatewayForDeviceServer).GetGatewayInfo(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _RpcGatewayForDevice_serviceDesc = grpc.ServiceDesc{
	ServiceName: "RpcGatewayForDevice",
	HandlerType: (*RpcGatewayForDeviceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGatewayInfo",
			Handler:    _RpcGatewayForDevice_GetGatewayInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gateway.proto",
}

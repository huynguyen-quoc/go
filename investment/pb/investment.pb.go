// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: investment/pb/investment.proto

package pb

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type StringMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *StringMessage) Reset() {
	*x = StringMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_investment_pb_investment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringMessage) ProtoMessage() {}

func (x *StringMessage) ProtoReflect() protoreflect.Message {
	mi := &file_investment_pb_investment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringMessage.ProtoReflect.Descriptor instead.
func (*StringMessage) Descriptor() ([]byte, []int) {
	return file_investment_pb_investment_proto_rawDescGZIP(), []int{0}
}

func (x *StringMessage) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_investment_pb_investment_proto protoreflect.FileDescriptor

var file_investment_pb_investment_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x69, 0x6e, 0x76, 0x65, 0x73, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x62, 0x2f,
	0x69, 0x6e, 0x76, 0x65, 0x73, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x25, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32, 0x68,
	0x0a, 0x11, 0x49, 0x6e, 0x76, 0x65, 0x73, 0x74, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x53, 0x0a, 0x04, 0x45, 0x63, 0x68, 0x6f, 0x12, 0x16, 0x2e, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x15, 0x22, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x2f, 0x65, 0x63, 0x68, 0x6f, 0x3a, 0x01, 0x2a, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x75, 0x79, 0x6e, 0x67, 0x75, 0x79, 0x65, 0x6e,
	0x2d, 0x71, 0x75, 0x6f, 0x63, 0x2f, 0x67, 0x6f, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x73, 0x74, 0x6d,
	0x65, 0x6e, 0x74, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_investment_pb_investment_proto_rawDescOnce sync.Once
	file_investment_pb_investment_proto_rawDescData = file_investment_pb_investment_proto_rawDesc
)

func file_investment_pb_investment_proto_rawDescGZIP() []byte {
	file_investment_pb_investment_proto_rawDescOnce.Do(func() {
		file_investment_pb_investment_proto_rawDescData = protoimpl.X.CompressGZIP(file_investment_pb_investment_proto_rawDescData)
	})
	return file_investment_pb_investment_proto_rawDescData
}

var file_investment_pb_investment_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_investment_pb_investment_proto_goTypes = []interface{}{
	(*StringMessage)(nil), // 0: example.StringMessage
}
var file_investment_pb_investment_proto_depIdxs = []int32{
	0, // 0: example.InvestmentService.Echo:input_type -> example.StringMessage
	0, // 1: example.InvestmentService.Echo:output_type -> example.StringMessage
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_investment_pb_investment_proto_init() }
func file_investment_pb_investment_proto_init() {
	if File_investment_pb_investment_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_investment_pb_investment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringMessage); i {
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
			RawDescriptor: file_investment_pb_investment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_investment_pb_investment_proto_goTypes,
		DependencyIndexes: file_investment_pb_investment_proto_depIdxs,
		MessageInfos:      file_investment_pb_investment_proto_msgTypes,
	}.Build()
	File_investment_pb_investment_proto = out.File
	file_investment_pb_investment_proto_rawDesc = nil
	file_investment_pb_investment_proto_goTypes = nil
	file_investment_pb_investment_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// InvestmentServiceClient is the client API for InvestmentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type InvestmentServiceClient interface {
	Echo(ctx context.Context, in *StringMessage, opts ...grpc.CallOption) (*StringMessage, error)
}

type investmentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInvestmentServiceClient(cc grpc.ClientConnInterface) InvestmentServiceClient {
	return &investmentServiceClient{cc}
}

func (c *investmentServiceClient) Echo(ctx context.Context, in *StringMessage, opts ...grpc.CallOption) (*StringMessage, error) {
	out := new(StringMessage)
	err := c.cc.Invoke(ctx, "/example.InvestmentService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InvestmentServiceServer is the server API for InvestmentService service.
type InvestmentServiceServer interface {
	Echo(context.Context, *StringMessage) (*StringMessage, error)
}

// UnimplementedInvestmentServiceServer can be embedded to have forward compatible implementations.
type UnimplementedInvestmentServiceServer struct {
}

func (*UnimplementedInvestmentServiceServer) Echo(context.Context, *StringMessage) (*StringMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}

func RegisterInvestmentServiceServer(s *grpc.Server, srv InvestmentServiceServer) {
	s.RegisterService(&_InvestmentService_serviceDesc, srv)
}

func _InvestmentService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StringMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InvestmentServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/example.InvestmentService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InvestmentServiceServer).Echo(ctx, req.(*StringMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _InvestmentService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "example.InvestmentService",
	HandlerType: (*InvestmentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _InvestmentService_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "investment/pb/investment.proto",
}

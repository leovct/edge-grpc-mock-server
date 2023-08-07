// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.4
// source: grpc/pb/server.proto

package pb

import (
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

// Response message for the GetStatus RPC method.
type StatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Network int64                 `protobuf:"varint,1,opt,name=network,proto3" json:"network,omitempty"`
	Genesis string                `protobuf:"bytes,2,opt,name=genesis,proto3" json:"genesis,omitempty"`
	Current *StatusResponse_Block `protobuf:"bytes,3,opt,name=current,proto3" json:"current,omitempty"`
	P2PAddr string                `protobuf:"bytes,4,opt,name=p2pAddr,proto3" json:"p2pAddr,omitempty"`
}

func (x *StatusResponse) Reset() {
	*x = StatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_pb_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusResponse) ProtoMessage() {}

func (x *StatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_pb_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusResponse.ProtoReflect.Descriptor instead.
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return file_grpc_pb_server_proto_rawDescGZIP(), []int{0}
}

func (x *StatusResponse) GetNetwork() int64 {
	if x != nil {
		return x.Network
	}
	return 0
}

func (x *StatusResponse) GetGenesis() string {
	if x != nil {
		return x.Genesis
	}
	return ""
}

func (x *StatusResponse) GetCurrent() *StatusResponse_Block {
	if x != nil {
		return x.Current
	}
	return nil
}

func (x *StatusResponse) GetP2PAddr() string {
	if x != nil {
		return x.P2PAddr
	}
	return ""
}

// Request message using a block number.
type BlockNumberRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The block number for which the data is requested.
	Number uint64 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
}

func (x *BlockNumberRequest) Reset() {
	*x = BlockNumberRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_pb_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockNumberRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockNumberRequest) ProtoMessage() {}

func (x *BlockNumberRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_pb_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockNumberRequest.ProtoReflect.Descriptor instead.
func (*BlockNumberRequest) Descriptor() ([]byte, []int) {
	return file_grpc_pb_server_proto_rawDescGZIP(), []int{1}
}

func (x *BlockNumberRequest) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

// Response message for the GetTrace RPC method.
type TraceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The trace of the block represented as a byte array.
	Trace []byte `protobuf:"bytes,1,opt,name=trace,proto3" json:"trace,omitempty"`
}

func (x *TraceResponse) Reset() {
	*x = TraceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_pb_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TraceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TraceResponse) ProtoMessage() {}

func (x *TraceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_pb_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TraceResponse.ProtoReflect.Descriptor instead.
func (*TraceResponse) Descriptor() ([]byte, []int) {
	return file_grpc_pb_server_proto_rawDescGZIP(), []int{2}
}

func (x *TraceResponse) GetTrace() []byte {
	if x != nil {
		return x.Trace
	}
	return nil
}

// Response message for the BlockByNumber RPC method.
type BlockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The data of the block represented as a byte array.
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *BlockResponse) Reset() {
	*x = BlockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_pb_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockResponse) ProtoMessage() {}

func (x *BlockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_pb_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockResponse.ProtoReflect.Descriptor instead.
func (*BlockResponse) Descriptor() ([]byte, []int) {
	return file_grpc_pb_server_proto_rawDescGZIP(), []int{3}
}

func (x *BlockResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type StatusResponse_Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The block number for which the data is requested.
	// Note: This is the only field of the StatusResponse message that's used by the prover at the moment.
	Number int64  `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	Hash   string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *StatusResponse_Block) Reset() {
	*x = StatusResponse_Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_pb_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusResponse_Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusResponse_Block) ProtoMessage() {}

func (x *StatusResponse_Block) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_pb_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusResponse_Block.ProtoReflect.Descriptor instead.
func (*StatusResponse_Block) Descriptor() ([]byte, []int) {
	return file_grpc_pb_server_proto_rawDescGZIP(), []int{0, 0}
}

func (x *StatusResponse_Block) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *StatusResponse_Block) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

var File_grpc_pb_server_proto protoreflect.FileDescriptor

var file_grpc_pb_server_proto_rawDesc = []byte{
	0x0a, 0x14, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc7, 0x01, 0x0a, 0x0e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x12, 0x32,
	0x0a, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x07, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x32, 0x70, 0x41, 0x64, 0x64, 0x72, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x32, 0x70, 0x41, 0x64, 0x64, 0x72, 0x1a, 0x33, 0x0a, 0x05,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x12, 0x0a,
	0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73,
	0x68, 0x22, 0x2c, 0x0a, 0x12, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22,
	0x25, 0x0a, 0x0d, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x72, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x05, 0x74, 0x72, 0x61, 0x63, 0x65, 0x22, 0x23, 0x0a, 0x0d, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xb4, 0x01, 0x0a, 0x06,
	0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x37, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x35, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x54, 0x72, 0x61, 0x63, 0x65, 0x12, 0x16, 0x2e, 0x76, 0x31,
	0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x0d, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x42,
	0x79, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x11, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_pb_server_proto_rawDescOnce sync.Once
	file_grpc_pb_server_proto_rawDescData = file_grpc_pb_server_proto_rawDesc
)

func file_grpc_pb_server_proto_rawDescGZIP() []byte {
	file_grpc_pb_server_proto_rawDescOnce.Do(func() {
		file_grpc_pb_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_pb_server_proto_rawDescData)
	})
	return file_grpc_pb_server_proto_rawDescData
}

var file_grpc_pb_server_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_grpc_pb_server_proto_goTypes = []interface{}{
	(*StatusResponse)(nil),       // 0: v1.StatusResponse
	(*BlockNumberRequest)(nil),   // 1: v1.BlockNumberRequest
	(*TraceResponse)(nil),        // 2: v1.TraceResponse
	(*BlockResponse)(nil),        // 3: v1.BlockResponse
	(*StatusResponse_Block)(nil), // 4: v1.StatusResponse.Block
	(*emptypb.Empty)(nil),        // 5: google.protobuf.Empty
}
var file_grpc_pb_server_proto_depIdxs = []int32{
	4, // 0: v1.StatusResponse.current:type_name -> v1.StatusResponse.Block
	5, // 1: v1.System.GetStatus:input_type -> google.protobuf.Empty
	1, // 2: v1.System.GetTrace:input_type -> v1.BlockNumberRequest
	1, // 3: v1.System.BlockByNumber:input_type -> v1.BlockNumberRequest
	0, // 4: v1.System.GetStatus:output_type -> v1.StatusResponse
	2, // 5: v1.System.GetTrace:output_type -> v1.TraceResponse
	3, // 6: v1.System.BlockByNumber:output_type -> v1.BlockResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_grpc_pb_server_proto_init() }
func file_grpc_pb_server_proto_init() {
	if File_grpc_pb_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_pb_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusResponse); i {
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
		file_grpc_pb_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockNumberRequest); i {
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
		file_grpc_pb_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TraceResponse); i {
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
		file_grpc_pb_server_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockResponse); i {
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
		file_grpc_pb_server_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusResponse_Block); i {
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
			RawDescriptor: file_grpc_pb_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_pb_server_proto_goTypes,
		DependencyIndexes: file_grpc_pb_server_proto_depIdxs,
		MessageInfos:      file_grpc_pb_server_proto_msgTypes,
	}.Build()
	File_grpc_pb_server_proto = out.File
	file_grpc_pb_server_proto_rawDesc = nil
	file_grpc_pb_server_proto_goTypes = nil
	file_grpc_pb_server_proto_depIdxs = nil
}

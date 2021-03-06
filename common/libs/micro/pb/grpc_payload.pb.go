// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc_payload.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type JsonRequest struct {
	Params               []byte   `protobuf:"bytes,1,opt,name=Params,proto3" json:"Params,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JsonRequest) Reset()         { *m = JsonRequest{} }
func (m *JsonRequest) String() string { return proto.CompactTextString(m) }
func (*JsonRequest) ProtoMessage()    {}
func (*JsonRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd3e348c9b0cbac2, []int{0}
}

func (m *JsonRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JsonRequest.Unmarshal(m, b)
}
func (m *JsonRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JsonRequest.Marshal(b, m, deterministic)
}
func (m *JsonRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JsonRequest.Merge(m, src)
}
func (m *JsonRequest) XXX_Size() int {
	return xxx_messageInfo_JsonRequest.Size(m)
}
func (m *JsonRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_JsonRequest.DiscardUnknown(m)
}

var xxx_messageInfo_JsonRequest proto.InternalMessageInfo

func (m *JsonRequest) GetParams() []byte {
	if m != nil {
		return m.Params
	}
	return nil
}

type JsonResponse struct {
	Result               []byte   `protobuf:"bytes,1,opt,name=Result,proto3" json:"Result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JsonResponse) Reset()         { *m = JsonResponse{} }
func (m *JsonResponse) String() string { return proto.CompactTextString(m) }
func (*JsonResponse) ProtoMessage()    {}
func (*JsonResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd3e348c9b0cbac2, []int{1}
}

func (m *JsonResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JsonResponse.Unmarshal(m, b)
}
func (m *JsonResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JsonResponse.Marshal(b, m, deterministic)
}
func (m *JsonResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JsonResponse.Merge(m, src)
}
func (m *JsonResponse) XXX_Size() int {
	return xxx_messageInfo_JsonResponse.Size(m)
}
func (m *JsonResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_JsonResponse.DiscardUnknown(m)
}

var xxx_messageInfo_JsonResponse proto.InternalMessageInfo

func (m *JsonResponse) GetResult() []byte {
	if m != nil {
		return m.Result
	}
	return nil
}

func init() {
	proto.RegisterType((*JsonRequest)(nil), "JsonRequest")
	proto.RegisterType((*JsonResponse)(nil), "JsonResponse")
}

func init() { proto.RegisterFile("grpc_payload.proto", fileDescriptor_cd3e348c9b0cbac2) }

var fileDescriptor_cd3e348c9b0cbac2 = []byte{
	// 114 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4a, 0x2f, 0x2a, 0x48,
	0x8e, 0x2f, 0x48, 0xac, 0xcc, 0xc9, 0x4f, 0x4c, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x52,
	0xe5, 0xe2, 0xf6, 0x2a, 0xce, 0xcf, 0x0b, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x12, 0xe3,
	0x62, 0x0b, 0x48, 0x2c, 0x4a, 0xcc, 0x2d, 0x96, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x09, 0x82, 0xf2,
	0x94, 0xd4, 0xb8, 0x78, 0x20, 0xca, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x41, 0xea, 0x82, 0x52,
	0x8b, 0x4b, 0x73, 0x4a, 0x60, 0xea, 0x20, 0x3c, 0x27, 0xd6, 0x28, 0x66, 0xfd, 0x82, 0xa4, 0x24,
	0x36, 0xb0, 0xe1, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd1, 0x92, 0x2a, 0x6b, 0x72, 0x00,
	0x00, 0x00,
}

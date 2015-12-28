// Code generated by protoc-gen-go.
// source: vapi.proto
// DO NOT EDIT!

/*
Package vapi is a generated protocol buffer package.

It is generated from these files:
	vapi.proto

It has these top-level messages:
	File
	PingRequest
	PingResponse
	ListRequest
	ListResponse
*/
package vapi

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// primitives
type File struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Size int64  `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
}

func (m *File) Reset()                    { *m = File{} }
func (m *File) String() string            { return proto.CompactTextString(m) }
func (*File) ProtoMessage()               {}
func (*File) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// RPC Ping
type PingRequest struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *PingRequest) Reset()                    { *m = PingRequest{} }
func (m *PingRequest) String() string            { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()               {}
func (*PingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type PingResponse struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *PingResponse) Reset()                    { *m = PingResponse{} }
func (m *PingResponse) String() string            { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()               {}
func (*PingResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// RPC ListFiles
type ListRequest struct {
	User string `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
	Path string `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
}

func (m *ListRequest) Reset()                    { *m = ListRequest{} }
func (m *ListRequest) String() string            { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()               {}
func (*ListRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type ListResponse struct {
	Request   *ListRequest `protobuf:"bytes,1,opt,name=request" json:"request,omitempty"`
	Layout    string       `protobuf:"bytes,2,opt,name=layout" json:"layout,omitempty"`
	UrlPrefix string       `protobuf:"bytes,3,opt,name=url_prefix" json:"url_prefix,omitempty"`
	Files     []*File      `protobuf:"bytes,4,rep,name=Files" json:"Files,omitempty"`
}

func (m *ListResponse) Reset()                    { *m = ListResponse{} }
func (m *ListResponse) String() string            { return proto.CompactTextString(m) }
func (*ListResponse) ProtoMessage()               {}
func (*ListResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ListResponse) GetRequest() *ListRequest {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *ListResponse) GetFiles() []*File {
	if m != nil {
		return m.Files
	}
	return nil
}

func init() {
	proto.RegisterType((*File)(nil), "vapi.File")
	proto.RegisterType((*PingRequest)(nil), "vapi.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "vapi.PingResponse")
	proto.RegisterType((*ListRequest)(nil), "vapi.ListRequest")
	proto.RegisterType((*ListResponse)(nil), "vapi.ListResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for VApi service

type VApiClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	ListFiles(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
}

type vApiClient struct {
	cc *grpc.ClientConn
}

func NewVApiClient(cc *grpc.ClientConn) VApiClient {
	return &vApiClient{cc}
}

func (c *vApiClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := grpc.Invoke(ctx, "/vapi.VApi/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vApiClient) ListFiles(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := grpc.Invoke(ctx, "/vapi.VApi/ListFiles", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for VApi service

type VApiServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	ListFiles(context.Context, *ListRequest) (*ListResponse, error)
}

func RegisterVApiServer(s *grpc.Server, srv VApiServer) {
	s.RegisterService(&_VApi_serviceDesc, srv)
}

func _VApi_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(VApiServer).Ping(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _VApi_ListFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(VApiServer).ListFiles(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _VApi_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vapi.VApi",
	HandlerType: (*VApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _VApi_Ping_Handler,
		},
		{
			MethodName: "ListFiles",
			Handler:    _VApi_ListFiles_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 250 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x90, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0xc5, 0x5b, 0x62, 0x5a, 0xe5, 0x25, 0x02, 0x71, 0x53, 0x60, 0x80, 0xca, 0x13, 0x2c, 0x45,
	0x2a, 0x7c, 0x01, 0x16, 0x26, 0x06, 0xc4, 0xc0, 0x8a, 0x8c, 0x64, 0x8a, 0xa5, 0xfc, 0x71, 0x7d,
	0x0e, 0x02, 0x3e, 0x3d, 0xb1, 0x9d, 0x8a, 0x48, 0xdd, 0x72, 0x97, 0x77, 0xbf, 0xf7, 0x9e, 0x81,
	0x2f, 0x65, 0xcd, 0xda, 0xba, 0xce, 0x77, 0x24, 0xc2, 0xb7, 0x94, 0x10, 0x8f, 0xa6, 0xd6, 0x54,
	0x42, 0xb4, 0xaa, 0xd1, 0xd5, 0x7c, 0x35, 0xbf, 0xce, 0xc3, 0xc4, 0xe6, 0x57, 0x57, 0x47, 0xc3,
	0x94, 0xc9, 0x4b, 0x14, 0xcf, 0xa6, 0xdd, 0xbe, 0xe8, 0x5d, 0xaf, 0xd9, 0xd3, 0x29, 0x96, 0x8d,
	0x66, 0x56, 0xdb, 0x51, 0x2d, 0xaf, 0x50, 0xa6, 0xff, 0x6c, 0xbb, 0x96, 0xf5, 0xa1, 0xe0, 0x06,
	0xc5, 0x93, 0x61, 0xbf, 0x07, 0x0c, 0xf4, 0x9e, 0xb5, 0xfb, 0xf7, 0xb2, 0xca, 0x7f, 0x46, 0xaf,
	0x5c, 0xee, 0x50, 0x26, 0xe9, 0xc8, 0x92, 0x58, 0xba, 0x74, 0x16, 0xe5, 0xc5, 0xe6, 0x6c, 0x1d,
	0x3b, 0x4c, 0x79, 0x27, 0x58, 0xd4, 0xea, 0xa7, 0xeb, 0x7d, 0x62, 0x10, 0x01, 0xbd, 0xab, 0xdf,
	0xac, 0xd3, 0x1f, 0xe6, 0xbb, 0xca, 0xe2, 0xee, 0x1c, 0xc7, 0xa1, 0x27, 0x57, 0x62, 0x95, 0x0d,
	0x14, 0x24, 0x4a, 0x58, 0x6d, 0x1a, 0x88, 0xd7, 0x07, 0x6b, 0xe8, 0x16, 0x22, 0xd4, 0xa0, 0xd1,
	0x61, 0x52, 0xf9, 0x82, 0xa6, 0xab, 0x94, 0x4c, 0xce, 0xe8, 0x1e, 0x79, 0x88, 0x11, 0xb9, 0x74,
	0x98, 0x6b, 0x7f, 0x35, 0xed, 0x23, 0x67, 0xef, 0x8b, 0xf8, 0xfc, 0x77, 0x7f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x69, 0x55, 0x0b, 0xf0, 0x8c, 0x01, 0x00, 0x00,
}

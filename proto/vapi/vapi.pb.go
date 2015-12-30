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
	ConfigRequest
	ConfigResponse
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
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Size  int64  `protobuf:"varint,2,opt,name=size" json:"size,omitempty"`
	Mtime int64  `protobuf:"varint,3,opt,name=mtime" json:"mtime,omitempty"`
	IsDir bool   `protobuf:"varint,4,opt,name=is_dir" json:"is_dir,omitempty"`
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
	Request      *ListRequest `protobuf:"bytes,1,opt,name=request" json:"request,omitempty"`
	IsDir        bool         `protobuf:"varint,2,opt,name=is_dir" json:"is_dir,omitempty"`
	Layout       string       `protobuf:"bytes,3,opt,name=layout" json:"layout,omitempty"`
	UrlPrefix    string       `protobuf:"bytes,4,opt,name=url_prefix" json:"url_prefix,omitempty"`
	RemoteServer string       `protobuf:"bytes,5,opt,name=remote_server" json:"remote_server,omitempty"`
	RelPath      string       `protobuf:"bytes,6,opt,name=rel_path" json:"rel_path,omitempty"`
	LocalPath    string       `protobuf:"bytes,7,opt,name=local_path" json:"local_path,omitempty"`
	Files        []*File      `protobuf:"bytes,8,rep,name=Files" json:"Files,omitempty"`
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

// RPC GetConfig
type ConfigRequest struct {
	User string `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
}

func (m *ConfigRequest) Reset()                    { *m = ConfigRequest{} }
func (m *ConfigRequest) String() string            { return proto.CompactTextString(m) }
func (*ConfigRequest) ProtoMessage()               {}
func (*ConfigRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type ConfigResponse struct {
	Allow   []string          `protobuf:"bytes,1,rep,name=allow" json:"allow,omitempty"`
	Alias   map[string]string `protobuf:"bytes,2,rep,name=alias" json:"alias,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Servers map[string]string `protobuf:"bytes,3,rep,name=servers" json:"servers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *ConfigResponse) Reset()                    { *m = ConfigResponse{} }
func (m *ConfigResponse) String() string            { return proto.CompactTextString(m) }
func (*ConfigResponse) ProtoMessage()               {}
func (*ConfigResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ConfigResponse) GetAlias() map[string]string {
	if m != nil {
		return m.Alias
	}
	return nil
}

func (m *ConfigResponse) GetServers() map[string]string {
	if m != nil {
		return m.Servers
	}
	return nil
}

func init() {
	proto.RegisterType((*File)(nil), "vapi.File")
	proto.RegisterType((*PingRequest)(nil), "vapi.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "vapi.PingResponse")
	proto.RegisterType((*ListRequest)(nil), "vapi.ListRequest")
	proto.RegisterType((*ListResponse)(nil), "vapi.ListResponse")
	proto.RegisterType((*ConfigRequest)(nil), "vapi.ConfigRequest")
	proto.RegisterType((*ConfigResponse)(nil), "vapi.ConfigResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for VApi service

type VApiClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	ListFiles(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	GetConfig(ctx context.Context, in *ConfigRequest, opts ...grpc.CallOption) (*ConfigResponse, error)
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

func (c *vApiClient) GetConfig(ctx context.Context, in *ConfigRequest, opts ...grpc.CallOption) (*ConfigResponse, error) {
	out := new(ConfigResponse)
	err := grpc.Invoke(ctx, "/vapi.VApi/GetConfig", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for VApi service

type VApiServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	ListFiles(context.Context, *ListRequest) (*ListResponse, error)
	GetConfig(context.Context, *ConfigRequest) (*ConfigResponse, error)
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

func _VApi_GetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(ConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(VApiServer).GetConfig(ctx, in)
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
		{
			MethodName: "GetConfig",
			Handler:    _VApi_GetConfig_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 437 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x53, 0xcd, 0xce, 0xd2, 0x40,
	0x14, 0xa5, 0xb4, 0xd0, 0xf6, 0xb6, 0xa0, 0x8e, 0x9a, 0x54, 0x12, 0x05, 0x67, 0xa5, 0x89, 0xa9,
	0x09, 0xba, 0x20, 0xee, 0x88, 0x51, 0x37, 0x2e, 0x8c, 0x26, 0x6e, 0xc9, 0xa8, 0x03, 0x4e, 0x6c,
	0x3b, 0x75, 0x66, 0x8a, 0xe2, 0x0b, 0xf9, 0x16, 0xdf, 0xbb, 0x7c, 0x6f, 0xf2, 0xcd, 0x4f, 0x1b,
	0x4a, 0x20, 0xf9, 0x76, 0xdc, 0x3b, 0xe7, 0x1c, 0xce, 0xb9, 0xf7, 0x16, 0x60, 0x4f, 0x6a, 0x96,
	0xd7, 0x82, 0x2b, 0x8e, 0x02, 0xf3, 0x1b, 0xaf, 0x21, 0x78, 0xcf, 0x0a, 0x8a, 0x52, 0x08, 0x2a,
	0x52, 0xd2, 0xcc, 0x5b, 0x78, 0xcf, 0x62, 0x53, 0x49, 0xf6, 0x8f, 0x66, 0x43, 0x5d, 0xf9, 0x68,
	0x02, 0xa3, 0x52, 0x31, 0xfd, 0xe8, 0xdb, 0x72, 0x0a, 0x63, 0x26, 0x37, 0x3f, 0x98, 0xc8, 0x02,
	0x5d, 0x47, 0xf8, 0x09, 0x24, 0x9f, 0x58, 0xb5, 0xfb, 0x4c, 0x7f, 0x37, 0x54, 0x2a, 0x74, 0x07,
	0xc2, 0x92, 0x4a, 0x49, 0x76, 0xad, 0x18, 0x9e, 0x43, 0xea, 0xde, 0x65, 0xcd, 0x2b, 0x49, 0xcf,
	0x01, 0xcf, 0x21, 0xf9, 0xc8, 0xa4, 0xea, 0x04, 0xf4, 0x9f, 0x37, 0x92, 0x8a, 0xa3, 0x95, 0x9a,
	0xa8, 0x9f, 0xd6, 0x4a, 0x8c, 0xaf, 0x3c, 0x48, 0x1d, 0xb6, 0x15, 0xc3, 0x10, 0x0a, 0xc7, 0xb3,
	0xf8, 0x64, 0x79, 0x2f, 0xb7, 0x19, 0xfb, 0x82, 0x47, 0xc3, 0x46, 0x24, 0x32, 0x75, 0x41, 0x0e,
	0xbc, 0x51, 0x36, 0x50, 0x8c, 0x10, 0x40, 0x23, 0x8a, 0x4d, 0x2d, 0xe8, 0x96, 0xfd, 0xb5, 0xa1,
	0x62, 0xf4, 0x10, 0x26, 0x82, 0x96, 0x5c, 0xd1, 0x8d, 0xb6, 0xb2, 0xd7, 0x6e, 0x46, 0xb6, 0x7d,
	0x17, 0x22, 0x41, 0x35, 0xd4, 0x38, 0x1a, 0x77, 0xe4, 0x82, 0x7f, 0x27, 0x6d, 0x2f, 0xb4, 0xbd,
	0x47, 0x30, 0x32, 0x43, 0x95, 0x59, 0xb4, 0xf0, 0xb5, 0x25, 0x70, 0x96, 0x4c, 0x0b, 0x3f, 0x86,
	0xc9, 0x5b, 0x5e, 0x6d, 0xd9, 0xee, 0x62, 0x5a, 0x7c, 0xed, 0xc1, 0xb4, 0x7b, 0x6f, 0x13, 0xea,
	0xe9, 0x93, 0xa2, 0xe0, 0x7f, 0x34, 0xc2, 0xd7, 0xda, 0xb9, 0x29, 0x19, 0x91, 0x3a, 0x8b, 0xd1,
	0x9e, 0x3b, 0xed, 0x53, 0x4e, 0xbe, 0x36, 0x88, 0x77, 0x95, 0x12, 0x07, 0xb4, 0x84, 0xd0, 0x25,
	0x90, 0x3a, 0xad, 0x61, 0x3c, 0xbd, 0xc8, 0xf8, 0xe2, 0x30, 0x96, 0x33, 0x7b, 0x01, 0xd0, 0x53,
	0x48, 0xc0, 0xff, 0x45, 0x0f, 0xed, 0x3a, 0xb4, 0x9b, 0x3d, 0x29, 0x1a, 0x77, 0x1a, 0xf1, 0x9b,
	0xe1, 0xca, 0x9b, 0xe5, 0x90, 0xf6, 0xd9, 0xb7, 0xe1, 0x97, 0xff, 0x3d, 0x08, 0xbe, 0xae, 0x6b,
	0x86, 0x5e, 0x42, 0x60, 0x0e, 0x03, 0xb5, 0x2b, 0xeb, 0x1d, 0xd1, 0x0c, 0xf5, 0x5b, 0xce, 0x22,
	0x1e, 0xa0, 0xd7, 0x10, 0x9b, 0xbd, 0xda, 0xd9, 0xa2, 0xf3, 0x45, 0x77, 0xac, 0xfe, 0x81, 0x68,
	0xd6, 0x0a, 0xe2, 0x0f, 0x54, 0xb9, 0xbc, 0xe8, 0xfe, 0x69, 0x7a, 0xc7, 0x7b, 0x70, 0x69, 0x24,
	0x78, 0xf0, 0x6d, 0x6c, 0xbf, 0x94, 0x57, 0x37, 0x01, 0x00, 0x00, 0xff, 0xff, 0xb1, 0x4a, 0xa1,
	0xf8, 0x37, 0x03, 0x00, 0x00,
}

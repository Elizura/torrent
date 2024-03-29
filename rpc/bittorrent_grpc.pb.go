// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: bittorrent.proto

package rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DownloadFileClient is the client API for DownloadFile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DownloadFileClient interface {
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadReply, error)
}

type downloadFileClient struct {
	cc grpc.ClientConnInterface
}

func NewDownloadFileClient(cc grpc.ClientConnInterface) DownloadFileClient {
	return &downloadFileClient{cc}
}

func (c *downloadFileClient) Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadReply, error) {
	out := new(DownloadReply)
	err := c.cc.Invoke(ctx, "/bittorrent.DownloadFile/Download", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DownloadFileServer is the server API for DownloadFile service.
// All implementations must embed UnimplementedDownloadFileServer
// for forward compatibility
type DownloadFileServer interface {
	Download(context.Context, *DownloadRequest) (*DownloadReply, error)
	mustEmbedUnimplementedDownloadFileServer()
}

// UnimplementedDownloadFileServer must be embedded to have forward compatible implementations.
type UnimplementedDownloadFileServer struct {
}

func (UnimplementedDownloadFileServer) Download(context.Context, *DownloadRequest) (*DownloadReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedDownloadFileServer) mustEmbedUnimplementedDownloadFileServer() {}

// UnsafeDownloadFileServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DownloadFileServer will
// result in compilation errors.
type UnsafeDownloadFileServer interface {
	mustEmbedUnimplementedDownloadFileServer()
}

func RegisterDownloadFileServer(s grpc.ServiceRegistrar, srv DownloadFileServer) {
	s.RegisterService(&DownloadFile_ServiceDesc, srv)
}

func _DownloadFile_Download_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DownloadFileServer).Download(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bittorrent.DownloadFile/Download",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DownloadFileServer).Download(ctx, req.(*DownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DownloadFile_ServiceDesc is the grpc.ServiceDesc for DownloadFile service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DownloadFile_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bittorrent.DownloadFile",
	HandlerType: (*DownloadFileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Download",
			Handler:    _DownloadFile_Download_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bittorrent.proto",
}

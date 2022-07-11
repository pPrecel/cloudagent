// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: pkg/agent/proto/route.proto

package cluster_agent

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

// AgentClient is the client API for Agent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentClient interface {
	GardenerShoots(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GardenerResponse, error)
	GCPClusters(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ClusterList, error)
}

type agentClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentClient(cc grpc.ClientConnInterface) AgentClient {
	return &agentClient{cc}
}

func (c *agentClient) GardenerShoots(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GardenerResponse, error) {
	out := new(GardenerResponse)
	err := c.cc.Invoke(ctx, "/Agent/GardenerShoots", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) GCPClusters(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ClusterList, error) {
	out := new(ClusterList)
	err := c.cc.Invoke(ctx, "/Agent/GCPClusters", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentServer is the server API for Agent service.
// All implementations must embed UnimplementedAgentServer
// for forward compatibility
type AgentServer interface {
	GardenerShoots(context.Context, *Empty) (*GardenerResponse, error)
	GCPClusters(context.Context, *Empty) (*ClusterList, error)
	mustEmbedUnimplementedAgentServer()
}

// UnimplementedAgentServer must be embedded to have forward compatible implementations.
type UnimplementedAgentServer struct {
}

func (UnimplementedAgentServer) GardenerShoots(context.Context, *Empty) (*GardenerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GardenerShoots not implemented")
}
func (UnimplementedAgentServer) GCPClusters(context.Context, *Empty) (*ClusterList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GCPClusters not implemented")
}
func (UnimplementedAgentServer) mustEmbedUnimplementedAgentServer() {}

// UnsafeAgentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentServer will
// result in compilation errors.
type UnsafeAgentServer interface {
	mustEmbedUnimplementedAgentServer()
}

func RegisterAgentServer(s grpc.ServiceRegistrar, srv AgentServer) {
	s.RegisterService(&Agent_ServiceDesc, srv)
}

func _Agent_GardenerShoots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).GardenerShoots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Agent/GardenerShoots",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).GardenerShoots(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_GCPClusters_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).GCPClusters(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Agent/GCPClusters",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).GCPClusters(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Agent_ServiceDesc is the grpc.ServiceDesc for Agent service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Agent_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Agent",
	HandlerType: (*AgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GardenerShoots",
			Handler:    _Agent_GardenerShoots_Handler,
		},
		{
			MethodName: "GCPClusters",
			Handler:    _Agent_GCPClusters_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/agent/proto/route.proto",
}

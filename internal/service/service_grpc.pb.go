// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package service

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

// MatchSimulatorClient is the client API for MatchSimulator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MatchSimulatorClient interface {
	SimulateGame(ctx context.Context, in *GameRequest, opts ...grpc.CallOption) (*MatchResults, error)
}

type matchSimulatorClient struct {
	cc grpc.ClientConnInterface
}

func NewMatchSimulatorClient(cc grpc.ClientConnInterface) MatchSimulatorClient {
	return &matchSimulatorClient{cc}
}

func (c *matchSimulatorClient) SimulateGame(ctx context.Context, in *GameRequest, opts ...grpc.CallOption) (*MatchResults, error) {
	out := new(MatchResults)
	err := c.cc.Invoke(ctx, "/MatchSimulator/SimulateGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MatchSimulatorServer is the server API for MatchSimulator service.
// All implementations must embed UnimplementedMatchSimulatorServer
// for forward compatibility
type MatchSimulatorServer interface {
	SimulateGame(context.Context, *GameRequest) (*MatchResults, error)
	mustEmbedUnimplementedMatchSimulatorServer()
}

// UnimplementedMatchSimulatorServer must be embedded to have forward compatible implementations.
type UnimplementedMatchSimulatorServer struct {
}

func (UnimplementedMatchSimulatorServer) SimulateGame(context.Context, *GameRequest) (*MatchResults, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SimulateGame not implemented")
}
func (UnimplementedMatchSimulatorServer) mustEmbedUnimplementedMatchSimulatorServer() {}

// UnsafeMatchSimulatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MatchSimulatorServer will
// result in compilation errors.
type UnsafeMatchSimulatorServer interface {
	mustEmbedUnimplementedMatchSimulatorServer()
}

func RegisterMatchSimulatorServer(s grpc.ServiceRegistrar, srv MatchSimulatorServer) {
	s.RegisterService(&MatchSimulator_ServiceDesc, srv)
}

func _MatchSimulator_SimulateGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MatchSimulatorServer).SimulateGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MatchSimulator/SimulateGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MatchSimulatorServer).SimulateGame(ctx, req.(*GameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MatchSimulator_ServiceDesc is the grpc.ServiceDesc for MatchSimulator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MatchSimulator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MatchSimulator",
	HandlerType: (*MatchSimulatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SimulateGame",
			Handler:    _MatchSimulator_SimulateGame_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

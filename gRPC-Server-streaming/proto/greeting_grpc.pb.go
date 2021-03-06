// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: proto/greeting.proto

package gRPC_Server_streaming

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

// GreetingClient is the client API for Greeting service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GreetingClient interface {
	GreetingMany(ctx context.Context, in *GreetingRequest, opts ...grpc.CallOption) (Greeting_GreetingManyClient, error)
	CalculatePrimeNum(ctx context.Context, in *Number, opts ...grpc.CallOption) (Greeting_CalculatePrimeNumClient, error)
}

type greetingClient struct {
	cc grpc.ClientConnInterface
}

func NewGreetingClient(cc grpc.ClientConnInterface) GreetingClient {
	return &greetingClient{cc}
}

func (c *greetingClient) GreetingMany(ctx context.Context, in *GreetingRequest, opts ...grpc.CallOption) (Greeting_GreetingManyClient, error) {
	stream, err := c.cc.NewStream(ctx, &Greeting_ServiceDesc.Streams[0], "/greeting.Greeting/GreetingMany", opts...)
	if err != nil {
		return nil, err
	}
	x := &greetingGreetingManyClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Greeting_GreetingManyClient interface {
	Recv() (*GreetingResponse, error)
	grpc.ClientStream
}

type greetingGreetingManyClient struct {
	grpc.ClientStream
}

func (x *greetingGreetingManyClient) Recv() (*GreetingResponse, error) {
	m := new(GreetingResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *greetingClient) CalculatePrimeNum(ctx context.Context, in *Number, opts ...grpc.CallOption) (Greeting_CalculatePrimeNumClient, error) {
	stream, err := c.cc.NewStream(ctx, &Greeting_ServiceDesc.Streams[1], "/greeting.Greeting/CalculatePrimeNum", opts...)
	if err != nil {
		return nil, err
	}
	x := &greetingCalculatePrimeNumClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Greeting_CalculatePrimeNumClient interface {
	Recv() (*PrimeResponse, error)
	grpc.ClientStream
}

type greetingCalculatePrimeNumClient struct {
	grpc.ClientStream
}

func (x *greetingCalculatePrimeNumClient) Recv() (*PrimeResponse, error) {
	m := new(PrimeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GreetingServer is the server API for Greeting service.
// All implementations must embed UnimplementedGreetingServer
// for forward compatibility
type GreetingServer interface {
	GreetingMany(*GreetingRequest, Greeting_GreetingManyServer) error
	CalculatePrimeNum(*Number, Greeting_CalculatePrimeNumServer) error
	mustEmbedUnimplementedGreetingServer()
}

// UnimplementedGreetingServer must be embedded to have forward compatible implementations.
type UnimplementedGreetingServer struct {
}

func (UnimplementedGreetingServer) GreetingMany(*GreetingRequest, Greeting_GreetingManyServer) error {
	return status.Errorf(codes.Unimplemented, "method GreetingMany not implemented")
}
func (UnimplementedGreetingServer) CalculatePrimeNum(*Number, Greeting_CalculatePrimeNumServer) error {
	return status.Errorf(codes.Unimplemented, "method CalculatePrimeNum not implemented")
}
func (UnimplementedGreetingServer) mustEmbedUnimplementedGreetingServer() {}

// UnsafeGreetingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GreetingServer will
// result in compilation errors.
type UnsafeGreetingServer interface {
	mustEmbedUnimplementedGreetingServer()
}

func RegisterGreetingServer(s grpc.ServiceRegistrar, srv GreetingServer) {
	s.RegisterService(&Greeting_ServiceDesc, srv)
}

func _Greeting_GreetingMany_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GreetingRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GreetingServer).GreetingMany(m, &greetingGreetingManyServer{stream})
}

type Greeting_GreetingManyServer interface {
	Send(*GreetingResponse) error
	grpc.ServerStream
}

type greetingGreetingManyServer struct {
	grpc.ServerStream
}

func (x *greetingGreetingManyServer) Send(m *GreetingResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Greeting_CalculatePrimeNum_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Number)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GreetingServer).CalculatePrimeNum(m, &greetingCalculatePrimeNumServer{stream})
}

type Greeting_CalculatePrimeNumServer interface {
	Send(*PrimeResponse) error
	grpc.ServerStream
}

type greetingCalculatePrimeNumServer struct {
	grpc.ServerStream
}

func (x *greetingCalculatePrimeNumServer) Send(m *PrimeResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Greeting_ServiceDesc is the grpc.ServiceDesc for Greeting service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Greeting_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "greeting.Greeting",
	HandlerType: (*GreetingServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GreetingMany",
			Handler:       _Greeting_GreetingMany_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "CalculatePrimeNum",
			Handler:       _Greeting_CalculatePrimeNum_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/greeting.proto",
}

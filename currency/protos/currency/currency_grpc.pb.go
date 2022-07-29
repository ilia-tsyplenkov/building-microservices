// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package currency

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

// CurrencyClient is the client API for Currency service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CurrencyClient interface {
	// GetRate returns the exchange rate for the two provided currency codes
	GetRate(ctx context.Context, in *RateRequest, opts ...grpc.CallOption) (*RateResponse, error)
	// SubscribeRates allows a client to subscribe for changes in an exchange rate
	// when the rate changes a response will be send
	SubscribeRates(ctx context.Context, opts ...grpc.CallOption) (Currency_SubscribeRatesClient, error)
}

type currencyClient struct {
	cc grpc.ClientConnInterface
}

func NewCurrencyClient(cc grpc.ClientConnInterface) CurrencyClient {
	return &currencyClient{cc}
}

func (c *currencyClient) GetRate(ctx context.Context, in *RateRequest, opts ...grpc.CallOption) (*RateResponse, error) {
	out := new(RateResponse)
	err := c.cc.Invoke(ctx, "/currency.Currency/GetRate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *currencyClient) SubscribeRates(ctx context.Context, opts ...grpc.CallOption) (Currency_SubscribeRatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Currency_ServiceDesc.Streams[0], "/currency.Currency/SubscribeRates", opts...)
	if err != nil {
		return nil, err
	}
	x := &currencySubscribeRatesClient{stream}
	return x, nil
}

type Currency_SubscribeRatesClient interface {
	Send(*RateRequest) error
	Recv() (*RateResponse, error)
	grpc.ClientStream
}

type currencySubscribeRatesClient struct {
	grpc.ClientStream
}

func (x *currencySubscribeRatesClient) Send(m *RateRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *currencySubscribeRatesClient) Recv() (*RateResponse, error) {
	m := new(RateResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CurrencyServer is the server API for Currency service.
// All implementations must embed UnimplementedCurrencyServer
// for forward compatibility
type CurrencyServer interface {
	// GetRate returns the exchange rate for the two provided currency codes
	GetRate(context.Context, *RateRequest) (*RateResponse, error)
	// SubscribeRates allows a client to subscribe for changes in an exchange rate
	// when the rate changes a response will be send
	SubscribeRates(Currency_SubscribeRatesServer) error
	mustEmbedUnimplementedCurrencyServer()
}

// UnimplementedCurrencyServer must be embedded to have forward compatible implementations.
type UnimplementedCurrencyServer struct {
}

func (UnimplementedCurrencyServer) GetRate(context.Context, *RateRequest) (*RateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRate not implemented")
}
func (UnimplementedCurrencyServer) SubscribeRates(Currency_SubscribeRatesServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeRates not implemented")
}
func (UnimplementedCurrencyServer) mustEmbedUnimplementedCurrencyServer() {}

// UnsafeCurrencyServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CurrencyServer will
// result in compilation errors.
type UnsafeCurrencyServer interface {
	mustEmbedUnimplementedCurrencyServer()
}

func RegisterCurrencyServer(s grpc.ServiceRegistrar, srv CurrencyServer) {
	s.RegisterService(&Currency_ServiceDesc, srv)
}

func _Currency_GetRate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CurrencyServer).GetRate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/currency.Currency/GetRate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CurrencyServer).GetRate(ctx, req.(*RateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Currency_SubscribeRates_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CurrencyServer).SubscribeRates(&currencySubscribeRatesServer{stream})
}

type Currency_SubscribeRatesServer interface {
	Send(*RateResponse) error
	Recv() (*RateRequest, error)
	grpc.ServerStream
}

type currencySubscribeRatesServer struct {
	grpc.ServerStream
}

func (x *currencySubscribeRatesServer) Send(m *RateResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *currencySubscribeRatesServer) Recv() (*RateRequest, error) {
	m := new(RateRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Currency_ServiceDesc is the grpc.ServiceDesc for Currency service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Currency_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "currency.Currency",
	HandlerType: (*CurrencyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRate",
			Handler:    _Currency_GetRate_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeRates",
			Handler:       _Currency_SubscribeRates_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "currency.proto",
}

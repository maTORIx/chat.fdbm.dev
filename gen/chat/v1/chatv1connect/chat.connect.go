// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: chat/v1/chat.proto

package chatv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/matorix/chat.fdbm.dev/gen/chat/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// ChatServiceName is the fully-qualified name of the ChatService service.
	ChatServiceName = "chat.v1.ChatService"
)

// ChatServiceClient is a client for the chat.v1.ChatService service.
type ChatServiceClient interface {
	Greet(context.Context, *connect_go.Request[v1.GreetRequest]) (*connect_go.Response[v1.GreetResponse], error)
	SendChat(context.Context, *connect_go.Request[v1.SendChatRequest]) (*connect_go.Response[v1.SendChatResponse], error)
	GetChats(context.Context, *connect_go.Request[v1.GetChatsRequest]) (*connect_go.Response[v1.GetChatsResponse], error)
	GetChatsStream(context.Context, *connect_go.Request[v1.GetChatsStreamRequest]) (*connect_go.ServerStreamForClient[v1.GetChatsStreamResponse], error)
}

// NewChatServiceClient constructs a client for the chat.v1.ChatService service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewChatServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ChatServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &chatServiceClient{
		greet: connect_go.NewClient[v1.GreetRequest, v1.GreetResponse](
			httpClient,
			baseURL+"/chat.v1.ChatService/Greet",
			opts...,
		),
		sendChat: connect_go.NewClient[v1.SendChatRequest, v1.SendChatResponse](
			httpClient,
			baseURL+"/chat.v1.ChatService/SendChat",
			opts...,
		),
		getChats: connect_go.NewClient[v1.GetChatsRequest, v1.GetChatsResponse](
			httpClient,
			baseURL+"/chat.v1.ChatService/GetChats",
			opts...,
		),
		getChatsStream: connect_go.NewClient[v1.GetChatsStreamRequest, v1.GetChatsStreamResponse](
			httpClient,
			baseURL+"/chat.v1.ChatService/GetChatsStream",
			opts...,
		),
	}
}

// chatServiceClient implements ChatServiceClient.
type chatServiceClient struct {
	greet          *connect_go.Client[v1.GreetRequest, v1.GreetResponse]
	sendChat       *connect_go.Client[v1.SendChatRequest, v1.SendChatResponse]
	getChats       *connect_go.Client[v1.GetChatsRequest, v1.GetChatsResponse]
	getChatsStream *connect_go.Client[v1.GetChatsStreamRequest, v1.GetChatsStreamResponse]
}

// Greet calls chat.v1.ChatService.Greet.
func (c *chatServiceClient) Greet(ctx context.Context, req *connect_go.Request[v1.GreetRequest]) (*connect_go.Response[v1.GreetResponse], error) {
	return c.greet.CallUnary(ctx, req)
}

// SendChat calls chat.v1.ChatService.SendChat.
func (c *chatServiceClient) SendChat(ctx context.Context, req *connect_go.Request[v1.SendChatRequest]) (*connect_go.Response[v1.SendChatResponse], error) {
	return c.sendChat.CallUnary(ctx, req)
}

// GetChats calls chat.v1.ChatService.GetChats.
func (c *chatServiceClient) GetChats(ctx context.Context, req *connect_go.Request[v1.GetChatsRequest]) (*connect_go.Response[v1.GetChatsResponse], error) {
	return c.getChats.CallUnary(ctx, req)
}

// GetChatsStream calls chat.v1.ChatService.GetChatsStream.
func (c *chatServiceClient) GetChatsStream(ctx context.Context, req *connect_go.Request[v1.GetChatsStreamRequest]) (*connect_go.ServerStreamForClient[v1.GetChatsStreamResponse], error) {
	return c.getChatsStream.CallServerStream(ctx, req)
}

// ChatServiceHandler is an implementation of the chat.v1.ChatService service.
type ChatServiceHandler interface {
	Greet(context.Context, *connect_go.Request[v1.GreetRequest]) (*connect_go.Response[v1.GreetResponse], error)
	SendChat(context.Context, *connect_go.Request[v1.SendChatRequest]) (*connect_go.Response[v1.SendChatResponse], error)
	GetChats(context.Context, *connect_go.Request[v1.GetChatsRequest]) (*connect_go.Response[v1.GetChatsResponse], error)
	GetChatsStream(context.Context, *connect_go.Request[v1.GetChatsStreamRequest], *connect_go.ServerStream[v1.GetChatsStreamResponse]) error
}

// NewChatServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewChatServiceHandler(svc ChatServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/chat.v1.ChatService/Greet", connect_go.NewUnaryHandler(
		"/chat.v1.ChatService/Greet",
		svc.Greet,
		opts...,
	))
	mux.Handle("/chat.v1.ChatService/SendChat", connect_go.NewUnaryHandler(
		"/chat.v1.ChatService/SendChat",
		svc.SendChat,
		opts...,
	))
	mux.Handle("/chat.v1.ChatService/GetChats", connect_go.NewUnaryHandler(
		"/chat.v1.ChatService/GetChats",
		svc.GetChats,
		opts...,
	))
	mux.Handle("/chat.v1.ChatService/GetChatsStream", connect_go.NewServerStreamHandler(
		"/chat.v1.ChatService/GetChatsStream",
		svc.GetChatsStream,
		opts...,
	))
	return "/chat.v1.ChatService/", mux
}

// UnimplementedChatServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedChatServiceHandler struct{}

func (UnimplementedChatServiceHandler) Greet(context.Context, *connect_go.Request[v1.GreetRequest]) (*connect_go.Response[v1.GreetResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("chat.v1.ChatService.Greet is not implemented"))
}

func (UnimplementedChatServiceHandler) SendChat(context.Context, *connect_go.Request[v1.SendChatRequest]) (*connect_go.Response[v1.SendChatResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("chat.v1.ChatService.SendChat is not implemented"))
}

func (UnimplementedChatServiceHandler) GetChats(context.Context, *connect_go.Request[v1.GetChatsRequest]) (*connect_go.Response[v1.GetChatsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("chat.v1.ChatService.GetChats is not implemented"))
}

func (UnimplementedChatServiceHandler) GetChatsStream(context.Context, *connect_go.Request[v1.GetChatsStreamRequest], *connect_go.ServerStream[v1.GetChatsStreamResponse]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("chat.v1.ChatService.GetChatsStream is not implemented"))
}
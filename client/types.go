package client

import (
	"context"
	"go-geektime3/client/user_service"
)

type Request struct {
	ServiceName string `json:"serviceName"`
	MethodName  string `json:"methodName"`
	Args        []byte `json:"args"`
}

type Response struct {
	Data []byte `json:"data"`
}

type proxy interface {
	invoke(ctx context.Context, request *Request) (*Response, error)
}

var _ proxy = &Client{}

type service interface {
	Name() string
}

var _ service = &user_service.UserService{}

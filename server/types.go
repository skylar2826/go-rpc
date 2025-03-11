package server

import (
	"context"
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

var _ proxy = &Server{}

type service interface {
	Name() string
}

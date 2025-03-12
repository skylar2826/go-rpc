package common

import "context"

type Proxy interface {
	Invoke(ctx context.Context, request *Request) (*Response, error)
}

type Service interface {
	Name() string
}

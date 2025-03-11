package client

import (
	"context"
	"encoding/json"
	"github.com/silenceper/pool"
	"net"
	"time"
)

type Client struct {
	connMsg
	pool pool.Pool
}

func (c *Client) invoke(ctx context.Context, request *Request) (*Response, error) {
	reqBs, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	var respBs []byte
	respBs, err = c.Send(ctx, reqBs)
	if err != nil {
		return nil, err
	}
	resp := &Response{}
	err = json.Unmarshal(respBs, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Send(ctx context.Context, req []byte) ([]byte, error) {
	p, err := c.pool.Get()
	if err != nil {
		_ = c.pool.Close(p)
		return nil, err
	}

	conn := p.(net.Conn)
	err = c.writeMsg(ctx, conn, req)
	if err != nil {
		_ = c.pool.Close(p)
		return nil, err
	}

	var respBs []byte
	respBs, err = c.readMsg(ctx, conn)
	if err != nil {
		_ = c.pool.Close(p)
		return nil, err
	}
	return respBs, nil
}

func NewClient(addr string, timeout time.Duration) (*Client, error) {
	conn, err := pool.NewChannelPool(&pool.Config{
		InitialCap: 1,
		MaxCap:     30,
		MaxIdle:    10,
		Factory: func() (interface{}, error) {
			return net.DialTimeout("tcp", addr, timeout)
		},
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		pool: conn,
	}, nil
}

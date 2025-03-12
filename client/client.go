package client

import (
	"context"
	"encoding/json"
	"github.com/silenceper/pool"
	"go-geektime3/common"
	"go-geektime3/serializer"
	"net"
	"reflect"
	"time"
)

type Client struct {
	common.ConnMsg
	pool       pool.Pool
	serializer serializer.Serializer
}

func (c *Client) Invoke(ctx context.Context, request *common.Request) (*common.Response, error) {
	reqBs, err := request.Encode()
	if err != nil {
		return nil, err
	}
	var respBs []byte
	respBs, err = c.Send(ctx, reqBs)
	if err != nil {
		return nil, err
	}
	resp := &common.Response{}
	err = resp.Decode(respBs)
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
	err = c.WriteMsg(ctx, conn, req)
	if err != nil {
		_ = c.pool.Close(p)
		return nil, err
	}

	var respBs []byte
	respBs, err = c.ReadMsg(ctx, conn)
	if err != nil {
		_ = c.pool.Close(p)
		return nil, err
	}
	return respBs, nil
}

func (c *Client) BindProxy(s common.Service) error {
	typ := reflect.TypeOf(s)
	val := reflect.ValueOf(s)

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		fTyp := typ.Field(i)
		fVal := val.Field(i)

		if fVal.CanSet() {
			fn := reflect.MakeFunc(fTyp.Type, func(args []reflect.Value) (results []reflect.Value) {
				ctx := args[0].Interface().(context.Context)
				reqData, err := json.Marshal(args[1].Interface())
				resVal := reflect.New(fTyp.Type.Out(0).Elem())
				if err != nil {
					return []reflect.Value{
						resVal,
						reflect.ValueOf(err),
					}
				}

				req := &common.Request{
					Serializer:  c.serializer.Code(),
					ServiceName: s.Name(),
					MethodName:  fTyp.Name,
					Data:        reqData,
				}

				var res *common.Response
				res, err = c.Invoke(ctx, req)
				if err != nil {
					return []reflect.Value{
						resVal,
						reflect.ValueOf(err),
					}
				}

				err = json.Unmarshal(res.Data, resVal.Interface())
				if err != nil {
					return []reflect.Value{
						resVal,
						reflect.ValueOf(err),
					}
				}

				return []reflect.Value{
					resVal,
					reflect.Zero(reflect.TypeOf(new(error)).Elem()),
				}
			})

			fVal.Set(fn)
		}
	}

	return nil
}

type Opt func(c *Client)

func WithSerializer(serializer serializer.Serializer) Opt {
	return func(c *Client) {
		c.serializer = serializer
	}
}

func NewClient(addr string, timeout time.Duration, opts ...Opt) (*Client, error) {
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

	c := &Client{
		pool:       conn,
		serializer: &serializer.Json{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

var _ common.Proxy = &Client{}

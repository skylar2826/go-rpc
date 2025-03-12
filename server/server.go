package server

import (
	"context"
	"fmt"
	"go-geektime3/common"
	"go-geektime3/serializer"
	"log"
	"net"
	"reflect"
	"strconv"
)

type Server struct {
	common.ConnMsg
	addr        string
	services    map[string]common.Service
	serializers map[string]serializer.Serializer
}

func (s *Server) Invoke(ctx context.Context, request *common.Request) (*common.Response, error) {
	si, ok := s.services[request.ServiceName]
	if !ok {
		return nil, fmt.Errorf("服务不存在 %s\n", request.ServiceName)
	}

	typ := reflect.TypeOf(si)
	var method reflect.Method
	method, ok = typ.MethodByName(request.MethodName)
	if !ok {
		return nil, fmt.Errorf("服务 %s 上的方法 %s 不存在 \n", request.ServiceName, request.MethodName)
	}

	in := make([]reflect.Value, 3)
	in[0] = reflect.ValueOf(si)
	in[1] = reflect.ValueOf(ctx)
	req := reflect.New(method.Type.In(2).Elem())

	// 获取序列化协议
	serializerCode := request.Serializer
	sl := s.serializers[strconv.Itoa(int(serializerCode))]
	// 解码请求参数
	err := sl.DeCode(request.Data, req.Interface())
	if err != nil {
		return nil, err
	}
	in[2] = req
	result := method.Func.Call(in)
	if result[1].Interface() != nil {
		return nil, result[1].Interface().(error)
	}
	var data []byte
	// 编码数据
	data, err = sl.Encode(result[0].Interface())
	if err != nil {
		return nil, err
	}
	return &common.Response{
		Data: data,
	}, nil
}

func (s *Server) RegisterService(services ...common.Service) *Server {
	if s.services == nil {
		s.services = make(map[string]common.Service, len(services))
	}
	for _, si := range services {
		s.services[si.Name()] = si
	}
	return s
}

func (s *Server) Start() {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer listen.Close()

	for {
		var conn net.Conn
		conn, err = listen.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		ctx := context.Background()

		var resp *common.Response
		var respBs []byte

		var handleErr = func() {
			resp.Error = []byte("error")
			respBs, err = resp.Encode()
			if err != nil {
				log.Println(err)
				return
			}
			err = s.WriteMsg(ctx, conn, respBs)
			if err != nil {
				log.Println(err)
				return
			}
		}

		var reqBs []byte
		// 读数据包
		reqBs, err = s.ReadMsg(ctx, conn)
		if err != nil {
			handleErr()
		}
		req := &common.Request{}
		// 对请求解码
		err = req.Decode(reqBs)
		if err != nil {
			handleErr()
		}

		resp, err = s.Invoke(ctx, req)
		if err != nil {
			handleErr()
		}

		respBs, err = resp.Encode()
		if err != nil {
			handleErr()
		}
		err = s.WriteMsg(ctx, conn, respBs)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

type Opt func(c *Server)

func WithSerializer(serializer serializer.Serializer) Opt {
	return func(c *Server) {
		c.serializers[strconv.Itoa(int(serializer.Code()))] = serializer
	}
}

func NewServer(addr string, opts ...Opt) *Server {
	j := &serializer.Json{}
	s := &Server{
		addr: addr,
		serializers: map[string]serializer.Serializer{
			strconv.Itoa(int(j.Code())): j,
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

var _ common.Proxy = &Server{}

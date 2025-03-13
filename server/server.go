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
	"time"
)

type Server struct {
	common.ConnMsg
	addr        string
	services    map[string]*ReflectionStub
	serializers map[string]serializer.Serializer
}

func (s *Server) Invoke(ctx context.Context, request *common.Request) (*common.Response, error) {
	si, ok := s.services[request.ServiceName]
	if !ok {
		return nil, fmt.Errorf("服务不存在 %s\n", request.ServiceName)
	}

	resp := &common.Response{
		Serializer: request.Serializer,
	}

	var cancel context.CancelFunc = func() {

	}
	deadline, err :=
		strconv.ParseInt(request.Meta["deadline"], 10, 64)
	if err != nil {
		log.Println("超时时间解析错误", err)
	} else {
		ctx, cancel = context.WithDeadline(ctx, time.UnixMilli(deadline))
	}

	var hasOneWay string
	hasOneWay, ok = request.Meta["one-way"]
	if ok && hasOneWay == "true" {
		go func() {
			_, err = si.Invoke(ctx, request)
			if err != nil {
				log.Println(err)
				return
			}
		}()
		resp.Error = []byte("已开启异步调用，无返回值")
		cancel()
	} else {
		var respData []byte
		respData, err = si.Invoke(ctx, request)
		if err != nil {
			resp.Error = []byte(err.Error())
			cancel()
			return resp, nil
		}
		resp.Data = respData
	}
	defer cancel()

	return resp, nil
}

func (s *Server) RegisterService(services ...common.Service) *Server {
	if s.services == nil {
		s.services = make(map[string]*ReflectionStub, len(services))
	}
	for _, si := range services {
		s.services[si.Name()] = &ReflectionStub{
			service:     si,
			serializers: s.serializers,
		}
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

		var handleErr = func(err error) {
			resp.Error = []byte(err.Error())
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
			handleErr(err)
		}
		req := &common.Request{}
		// 对请求解码
		err = req.Decode(reqBs)
		if err != nil {
			handleErr(err)
		}

		resp, err = s.Invoke(ctx, req)
		if err != nil {
			handleErr(err)
		}

		respBs, err = resp.Encode()
		if err != nil {
			handleErr(err)
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

type ReflectionStub struct {
	serializers map[string]serializer.Serializer
	service     common.Service
}

func (r *ReflectionStub) Invoke(ctx context.Context, request *common.Request) ([]byte, error) {
	typ := reflect.TypeOf(r.service)
	method, ok := typ.MethodByName(request.MethodName)
	if !ok {
		return nil, fmt.Errorf("服务 %s 上的方法 %s 不存在 \n", request.ServiceName, request.MethodName)
	}

	in := make([]reflect.Value, 3)
	in[0] = reflect.ValueOf(r.service)
	in[1] = reflect.ValueOf(ctx)
	methodIn := reflect.New(method.Type.In(2).Elem())

	// 获取序列化协议
	serializerCode := request.Serializer
	sl := r.serializers[strconv.Itoa(int(serializerCode))]
	// 解码请求参数
	err := sl.DeCode(request.Data, methodIn.Interface())
	if err != nil {
		return nil, err
	}
	in[2] = methodIn
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
	return data, nil
}

var _ common.Proxy = &Server{}

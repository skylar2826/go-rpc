package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
)

type Server struct {
	connMsg
	addr     string
	services map[string]service
}

func (s *Server) invoke(ctx context.Context, request *Request) (*Response, error) {
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
	err := json.Unmarshal(request.Args, req.Interface())
	if err != nil {
		return nil, err
	}
	in[2] = req
	result := method.Func.Call(in)
	if result[1].Interface() != nil {
		return nil, result[1].Interface().(error)
	}
	var data []byte
	data, err = json.Marshal(result[0].Interface())
	if err != nil {
		return nil, err
	}
	return &Response{
		Data: data,
	}, nil
}

func (s *Server) RegisterService(services ...service) *Server {
	if s.services == nil {
		s.services = make(map[string]service, len(services))
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

		var resp *Response
		var respBs []byte

		var handleErr = func() {
			resp.Data = []byte("error")
			respBs, err = json.Marshal(resp)
			err = s.writeMsg(ctx, conn, respBs)
			if err != nil {
				log.Println(err)
				return
			}
		}

		var reqBs []byte
		reqBs, err = s.readMsg(ctx, conn)
		if err != nil {
			handleErr()
		}
		req := &Request{}
		err = json.Unmarshal(reqBs, req)
		if err != nil {
			handleErr()
		}

		resp, err = s.invoke(ctx, req)
		if err != nil {
			handleErr()
		}

		respBs, err = json.Marshal(resp)
		if err != nil {
			handleErr()
		}
		err = s.writeMsg(ctx, conn, respBs)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

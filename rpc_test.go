package main

import (
	"testing"
)

func TestRPC(t *testing.T) {
	//go func() {
	//	s := server.NewServer("127.0.0.1:8080")
	//	userServiceServer := &user_service2.UserService{}
	//	s.RegisterService(userServiceServer)
	//
	//	s.Start()
	//}()
	//
	//time.Sleep(time.Second * 3)
	//
	//for i := 0; i < 3; i++ {
	//	c, err := client.NewClient("127.0.0.1:8080", time.Minute)
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//	userServiceClient := &user_service.UserService{}
	//	err = client.BindProxy(userServiceClient, c)
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//
	//	req := &user_service.UserServiceReq{Id: i}
	//	ctx := context.Background()
	//	var resp *user_service.UserServiceResp
	//	resp, err = userServiceClient.GetUserById(ctx, req)
	//	fmt.Println(resp)
	//}
}

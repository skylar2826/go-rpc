package user_service

import (
	"context"
	"time"
)

type User struct {
	Id        int    `json:"id"`
	Birthdate string `json:"birthdate"`
	Email     string `json:"email"`
}

type UserService struct{}

func (u *UserService) Name() string {
	return "user_service"
}

func (u *UserService) GetUserById(ctx context.Context, UserServiceReq *UserServiceReq) (*UserServiceResp, error) {
	time.Sleep(time.Second)
	return &UserServiceResp{
		Data: &User{
			Id:        UserServiceReq.Id,
			Birthdate: "2020-05-05",
			Email:     "135",
		},
	}, nil
}

package user_service

import (
	"context"
	"go-geektime3/common"
)

type UserService struct {
	GetUserById func(ctx context.Context, request *UserServiceReq) (*UserServiceResp, error)
}

func (u *UserService) Name() string {
	return "user_service"
}

var _ common.Service = &UserService{}

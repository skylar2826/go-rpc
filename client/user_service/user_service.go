package user_service

import (
	"context"
)

type UserService struct {
	GetUserById func(ctx context.Context, request *UserServiceReq) (*UserServiceResp, error)
}

func (u *UserService) Name() string {
	return "user_service"
}

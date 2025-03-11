package user_service

import "context"

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
	return &UserServiceResp{
		Data: &User{
			Id:        UserServiceReq.Id,
			Birthdate: "2020-05-05",
			Email:     "135",
		},
	}, nil
}

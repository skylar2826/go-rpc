package user_service

type UserServiceReq struct {
	Id int `json:"id"`
}

type UserServiceResp struct {
	Data   *User  `json:"data"`
	ErrMsg string `json:"errMsg"`
}

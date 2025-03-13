package user_service

type User struct {
	Id        int    `json:"id"`
	Birthdate string `json:"birthdate"`
	Email     string `json:"email"`
}

type UserServiceReq struct {
	Id int `json:"id"`
}

type UserServiceResp struct {
	Data *User `json:"data"`
}

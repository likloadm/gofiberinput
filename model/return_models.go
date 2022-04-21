package model

type MessageModel struct {
	Message string `json:"message" example:"error"`
	Success bool   `json:"success" example:"false"`
}

type AccessTokenJWT struct {
	Token string `json:"token"`
}

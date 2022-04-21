package model

type MessageModel struct {
	Message string `json:"message" example:"error"`
	Success bool   `json:"success" example:"false"`
}

type MessageOk struct {
	Message string `json:"message" example:"ok"`
	Success bool   `json:"success" example:"true"`
}

type AccessTokenJWT struct {
	Token string `json:"token"`
}

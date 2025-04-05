package models

type LoginInput struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginOutput struct {
	ErrorCode int    `json:"error_code"`
	Token     string `json:"token"`
}

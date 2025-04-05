package models

type RegisterInput struct {
	Username string   `json:"username" binding:"required,min=3"`
	Password string   `json:"password" binding:"required,min=6"`
	Scopes   []string `json:"scopes"`
}

type RegisterOutput struct {
	ErrorCode int `json:"error_code"`
}

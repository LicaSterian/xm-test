package models

type RegisterInput struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Scopes   []string `json:"scopes"`
}

type RegisterOutput struct {
	ErrorCode int `json:"error_code"`
}

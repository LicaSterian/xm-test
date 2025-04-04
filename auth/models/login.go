package models

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginOutput struct {
	ErrorCode int    `json:"errorCode"`
	Token     string `json:"token"`
}

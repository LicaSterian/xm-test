package handlers

import "errors"

const (
	errMessageInvalidInput          string = "invalid input"
	errMessageAuthenticationFailed  string = "authentication failed"
	errMessageCouldNotGenerateToken string = "could not generate token"
	errMessageRegistrationFailed    string = "registration failed"
)

var (
	ErrInvalidInput          = errors.New(errMessageInvalidInput)
	ErrAuthFailed            = errors.New(errMessageAuthenticationFailed)
	ErrCouldNotGenerateToken = errors.New(errMessageCouldNotGenerateToken)
	ErrRegistrationFailed    = errors.New(errMessageRegistrationFailed)
)

const (
	ErrCodeInvalidInput          int = 1
	ErrCodeAuthFailed            int = 2
	ErrCodeCouldNotGenerateToken int = 3
	ErrCodeRegistrationFailed    int = 4
)

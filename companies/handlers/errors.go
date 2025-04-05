package handlers

import "errors"

const (
	errMessageInvalidInput          string = "invalid input"
	errMessageCouldNotCreateCompany string = "could not create company"
)

var (
	ErrInvalidInput          = errors.New(errMessageInvalidInput)
	ErrCouldNotCreateCompany = errors.New(errMessageCouldNotCreateCompany)
)

const (
	ErrCodeInvalidInput          int = 1
	ErrCodeCouldNotCreateCompany int = 2
)

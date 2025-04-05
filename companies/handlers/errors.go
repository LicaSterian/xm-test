package handlers

import "errors"

const (
	errMessageInvalidInput          string = "invalid input"
	errMessageCouldNotCreateCompany string = "could not create company"
	errMessageInvalidId             string = "invalid id"
	errMessageGetCompany            string = "error while getting company"
)

var (
	ErrInvalidInput          = errors.New(errMessageInvalidInput)
	ErrCouldNotCreateCompany = errors.New(errMessageCouldNotCreateCompany)
	ErrInvalidId             = errors.New(errMessageInvalidId)
	ErrGetCompany            = errors.New(errMessageGetCompany)
)

const (
	ErrCodeInvalidInput          int = 1
	ErrCodeCouldNotCreateCompany int = 2
	ErrCodeInvalidId             int = 3
	ErrCodeGetCompany            int = 4
)

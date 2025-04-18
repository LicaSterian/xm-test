// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	gin "github.com/gin-gonic/gin"

	mock "github.com/stretchr/testify/mock"
)

// CompanyHandler is an autogenerated mock type for the CompanyHandler type
type CompanyHandler struct {
	mock.Mock
}

// CreateCompany provides a mock function with given fields: c
func (_m *CompanyHandler) CreateCompany(c *gin.Context) {
	_m.Called(c)
}

// DeleteCompany provides a mock function with given fields: c
func (_m *CompanyHandler) DeleteCompany(c *gin.Context) {
	_m.Called(c)
}

// GetCompany provides a mock function with given fields: c
func (_m *CompanyHandler) GetCompany(c *gin.Context) {
	_m.Called(c)
}

// PatchCompany provides a mock function with given fields: c
func (_m *CompanyHandler) PatchCompany(c *gin.Context) {
	_m.Called(c)
}

// NewCompanyHandler creates a new instance of CompanyHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCompanyHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *CompanyHandler {
	mock := &CompanyHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

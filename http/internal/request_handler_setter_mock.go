package internal

import (
	"net/http"
	"reflect"

	"github.com/stretchr/testify/mock"
)

var _ RequestHandlerSetter = (*RequestHandlerSetterMock)(nil)

type RequestHandlerSetterMock struct {
	mock.Mock
}

func (r *RequestHandlerSetterMock) Body(
	value reflect.Value, field reflect.Type, request *http.Request) error {
	return r.Called(value, field, request).Error(0)
}

func (r *RequestHandlerSetterMock) Header(
	value reflect.Value, field reflect.Type, request *http.Request, str string) error {
	return r.Called(value, field, request, str).Error(0)
}

func (r *RequestHandlerSetterMock) Path(
	value reflect.Value, field reflect.Type, request *http.Request, str string) error {
	return r.Called(value, field, request, str).Error(0)
}

func (r *RequestHandlerSetterMock) Query(
	value reflect.Value, field reflect.Type, request *http.Request, str string) error {
	return r.Called(value, field, request, str).Error(0)
}

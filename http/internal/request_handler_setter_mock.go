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
	value reflect.Value, request *http.Request,
) error {
	return r.Called(value, request).Error(0)
}

func (r *RequestHandlerSetterMock) Header(
	value reflect.Value, request *http.Request, str string,
) error {
	return r.Called(value, request, str).Error(0)
}

func (r *RequestHandlerSetterMock) Path(
	value reflect.Value, request *http.Request, str string,
) error {
	return r.Called(value, request, str).Error(0)
}

func (r *RequestHandlerSetterMock) Query(
	value reflect.Value, request *http.Request, str string,
) error {
	return r.Called(value, request, str).Error(0)
}

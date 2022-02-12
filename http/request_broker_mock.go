package http

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ RequestBroker = (*RequestMock)(nil)

type RequestMock struct {
	mock.Mock
}

func (r *RequestMock) Parameter(key, value string) RequestBroker {
	r.Called(key, value)

	return r
}

func (r *RequestMock) Header(key, value string) RequestBroker {
	r.Called(key, value)

	return r
}

func (r *RequestMock) Go(ctx context.Context, v interface{}) error {
	return r.Called(ctx, v).Error(0)
}

func (r *RequestMock) Post() RequestBroker {
	r.Called()

	return r
}

func (r *RequestMock) Get() RequestBroker {
	r.Called()

	return r
}

func (r *RequestMock) Put() RequestBroker {
	r.Called()

	return r
}

func (r *RequestMock) Delete() RequestBroker {
	r.Called()

	return r
}

func (r *RequestMock) Domain(s string) RequestBroker {
	r.Called(s)

	return r
}

func (r *RequestMock) Path(s string) RequestBroker {
	r.Called(s)

	return r
}

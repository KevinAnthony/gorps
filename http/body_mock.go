// Code generated by mockery v2.35.4. DO NOT EDIT.

package http

import mock "github.com/stretchr/testify/mock"

// BodyMock is an autogenerated mock type for the ReadCloser type
type BodyMock struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *BodyMock) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Read provides a mock function with given fields: p
func (_m *BodyMock) Read(p []byte) (int, error) {
	ret := _m.Called(p)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(p)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBodyMock creates a new instance of BodyMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBodyMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *BodyMock {
	mock := &BodyMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

package internal_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/kevinanthony/gorps/encoder"
	"github.com/kevinanthony/gorps/http/internal"
	"github.com/kevinanthony/gorps/internal/testx"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestNewRequestHandlerSetter(t *testing.T) {
	t.Parallel()

	Convey("NewRequestHandlerSetter", t, func() {
		factory := &encoder.FactoryMock{}

		Convey("should not panic", func() {
			So(func() { internal.NewRequestHandlerSetter(factory) }, ShouldNotPanic)
		})
		Convey("should panic", func() {
			So(func() { internal.NewRequestHandlerSetter(nil) }, ShouldPanicWith, "encoder factory is required")
		})
	})
}

func TestRequestHandlerSetter_Body(t *testing.T) {
	t.Parallel()

	Convey("Body", t, func() {
		Convey("should set everything", func() {
			factory := &encoder.FactoryMock{}

			setter := internal.NewRequestHandlerSetter(factory)

			expected := testx.GetTestStruct()
			actual := testx.TestStruct{}

			valueOf := reflect.ValueOf(&actual).Elem().Field(32).Addr()
			typeOf := reflect.TypeOf(&actual).Elem().Field(32).Type

			Convey("should set structure", func() {
				Convey("request is json", func() {
					req := httptest.NewRequest(http.MethodGet, "/", testx.ToReadCloser(encoder.NewJSON(), expected.Body))
					factory.On("CreateFromRequest", req).Return(encoder.NewJSON()).Once()

					err := setter.Body(valueOf, typeOf, req)

					So(err, ShouldBeNil)
					So(actual.Body, ShouldResemble, expected.Body)
					mock.AssertExpectationsForObjects(t, factory)
				})
				Convey("request is XML", func() {
					req := httptest.NewRequest(http.MethodGet, "/", testx.ToReadCloser(encoder.NewXML(), expected.Body))
					factory.On("CreateFromRequest", req).Return(encoder.NewXML()).Once()

					err := setter.Body(valueOf, typeOf, req)

					So(err, ShouldBeNil)
					So(actual.Body, ShouldResemble, expected.Body)
					mock.AssertExpectationsForObjects(t, factory)
				})
			})
			Convey("should return error when", func() {
				Convey("when body value is not valid", func() {
					valueOf := reflect.ValueOf(actual).Field(32).Elem()
					typeOf := reflect.TypeOf(&actual).Elem().Field(32).Type

					req := httptest.NewRequest(http.MethodGet, "/", testx.ToReadCloser(encoder.NewXML(), expected.Body))

					err := setter.Body(valueOf, typeOf, req)

					So(err, ShouldBeError, "bad body value")
					So(actual.Body, ShouldBeNil)
					mock.AssertExpectationsForObjects(t, factory)
				})
				Convey("input type is not valid", func() {
					valueOf := reflect.ValueOf(actual).Field(32)

					req := httptest.NewRequest(http.MethodGet, "/", testx.ToReadCloser(encoder.NewJSON(), expected.Body))

					err := setter.Body(valueOf, typeOf, req)

					So(err, ShouldBeError, "cannot set value to type")
					So(actual.Body, ShouldBeNil)
					mock.AssertExpectationsForObjects(t, factory)
				})
				Convey("when body type fails to unmarshal", func() {
					req := httptest.NewRequest(http.MethodGet, "/", ioutil.NopCloser(strings.NewReader("{")))
					factory.On("CreateFromRequest", req).Return(encoder.NewJSON()).Once()

					err := setter.Body(valueOf, typeOf, req)

					So(err, ShouldBeError)
					So(err.Error(), ShouldStartWith, "decode application/json: testx.JSONGambit")
					So(actual.Body, ShouldBeNil)
					mock.AssertExpectationsForObjects(t, factory)
				})
			})
		})
	})
}

func TestRequestHandlerSetter_Header(t *testing.T) {
	t.Parallel()
}

func TestRequestHandlerSetter_Path(t *testing.T) {
	t.Parallel()
}

func TestRequestHandlerSetter_Query(t *testing.T) {
	t.Parallel()
}

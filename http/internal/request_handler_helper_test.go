package internal_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/kevinanthony/gorps/v2/encoder"
	"github.com/kevinanthony/gorps/v2/http"
	"github.com/kevinanthony/gorps/v2/http/internal"
	"github.com/kevinanthony/gorps/v2/internal/testx"

	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestNewRequestHandlerHelper(t *testing.T) {
	t.Parallel()

	Convey("NewRequestHandlerHelper", t, func() {
		setter := &internal.RequestHandlerSetterMock{}
		Convey("should not panic", func() {
			So(func() { internal.NewRequestHandlerHelper(setter) }, ShouldNotPanic)
		})
		Convey("should panic", func() {
			So(func() { internal.NewRequestHandlerHelper(nil) }, ShouldPanicWith, "request handler setter is required")
		})
	})
}

func TestRequestHandlerHelper_Fill(t *testing.T) {
	t.Parallel()

	Convey("Fill", t, func() {
		// because reflection is hard, lets test that we are setting correctly
		setter := &internal.RequestHandlerSetterMock{}

		expected := testx.GetTestStruct()
		actual := testx.TestStruct{}

		helper := internal.NewRequestHandlerHelper(setter)

		cctx := &chi.Context{
			URLParams: chi.RouteParams{
				Keys: []string{
					"string",
				},
				Values: []string{
					expected.PathString,
				},
			},
		}

		req := httptest.
			NewRequest(http.MethodGet, "/", testx.ToReadCloser(encoder.NewJSON(), expected.Body)).
			WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, cctx))

		req.Header.Set("string", expected.HeaderString)
		q := req.URL.Query()
		q.Add("string", expected.QueryString)
		req.URL.RawQuery = q.Encode()

		setPathCall := setter.On("Path", mock.AnythingOfType("reflect.Value"), req, "string").Maybe()
		setHeaderCall := setter.On("Header", mock.AnythingOfType("reflect.Value"), req, "string").Maybe()
		setQueryCall := setter.On("Query", mock.AnythingOfType("reflect.Value"), req, "string").Maybe()
		setBodycall := setter.On("Body", mock.AnythingOfType("reflect.Value"), req).Maybe()
		extraPath := setter.On("Path", mock.AnythingOfType("reflect.Value"), req, mock.Anything).Return(nil).Maybe()
		extraHeader := setter.On("Header", mock.AnythingOfType("reflect.Value"), req, mock.Anything).Return(nil).Maybe()
		extraQuery := setter.On("Query", mock.AnythingOfType("reflect.Value"), req, mock.Anything).Return(nil).Maybe()

		Convey("should fill everything when there are no errors", func() {
			setPathCall.Return(nil).Once()
			setHeaderCall.Return(nil).Once()
			setQueryCall.Return(nil).Once()
			setBodycall.Return(nil).Once()
			extraPath.Times(4)
			extraHeader.Times(5)
			extraQuery.Times(5)

			err := helper.Fill(req, &actual)

			So(err, ShouldBeNil)
			mock.AssertExpectationsForObjects(t, setter)
		})
		Convey("should return error when", func() {
			Convey("set path returns an error", func() {
				setPathCall.Return(errors.New("bad path")).Once()
				err := helper.Fill(req, &actual)

				So(err, ShouldBeError, "bad path")
				mock.AssertExpectationsForObjects(t, setter)
			})
			Convey("set header returns an error", func() {
				setPathCall.Return(nil).Once()
				setHeaderCall.Return(errors.New("bad header")).Once()
				extraPath.Times(4)

				err := helper.Fill(req, &actual)

				So(err, ShouldBeError, "bad header")
				mock.AssertExpectationsForObjects(t, setter)
			})
			Convey("set query returns an error", func() {
				setPathCall.Return(nil).Once()
				setHeaderCall.Return(nil).Once()
				setQueryCall.Return(errors.New("bad query")).Once()
				extraPath.Times(4)
				extraHeader.Times(5)
				err := helper.Fill(req, &actual)

				So(err, ShouldBeError, "bad query")
				mock.AssertExpectationsForObjects(t, setter)
			})
			Convey("set body returns an error", func() {
				setPathCall.Return(nil).Once()
				setHeaderCall.Return(nil).Once()
				setQueryCall.Return(nil).Once()
				setBodycall.Return(errors.New("bad body")).Once()
				extraPath.Times(4)
				extraHeader.Times(5)
				extraQuery.Times(5)

				err := helper.Fill(req, &actual)

				So(err, ShouldBeError, "bad body")
				mock.AssertExpectationsForObjects(t, setter)
			})
		})
	})
}

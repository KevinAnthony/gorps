package internal_test

import (
	"context"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/kevinanthony/gorps/encoder"
	mocks "github.com/kevinanthony/gorps/http"
	"github.com/kevinanthony/gorps/http/internal"
	"github.com/kevinanthony/gorps/internal/testx"

	"github.com/go-chi/chi"
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
			reader := &mocks.BodyMock{}

			bag := []interface{}{factory, reader}

			setter := internal.NewRequestHandlerSetter(factory)

			expected := testx.GetTestStruct()
			actual := testx.TestStruct{}

			valueOf := getFields(&actual, "Body")

			Convey("should set structure", func() {
				Convey("request is json", func() {
					req := httptest.NewRequest(http.MethodGet, "/", testx.ToReadCloser(encoder.NewJSON(), expected.Body))
					factory.On("CreateFromRequest", req).Return(encoder.NewJSON()).Once()

					err := setter.Body(valueOf, req)

					So(err, ShouldBeNil)
					So(actual.Body, ShouldResemble, expected.Body)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("request is XML", func() {
					req := httptest.NewRequest(http.MethodGet, "/", testx.ToReadCloser(encoder.NewXML(), expected.Body))
					factory.On("CreateFromRequest", req).Return(encoder.NewXML()).Once()

					err := setter.Body(valueOf, req)

					So(err, ShouldBeNil)
					So(actual.Body, ShouldResemble, expected.Body)
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("should return error when", func() {
				Convey("io reader fails", func() {
					reader.On("Read", mock.Anything).Return(0, errors.New("everybody body mock"))

					req := httptest.NewRequest(http.MethodGet, "/", reader)

					err := setter.Body(valueOf, req)

					So(err, ShouldBeError, "everybody body mock")
					So(actual.Body, ShouldBeNil)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when body value is not valid", func() {
					valueOf := reflect.ValueOf(actual).FieldByName("Body").Elem()

					req := httptest.NewRequest(http.MethodGet, "/", testx.ToReadCloser(encoder.NewXML(), expected.Body))
					factory.On("CreateFromRequest", req).Return(encoder.NewXML()).Once()

					err := setter.Body(valueOf, req)

					So(err, ShouldBeError, "bad body value")
					So(actual.Body, ShouldBeNil)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("input type is not valid", func() {
					valueOf := reflect.ValueOf(actual).FieldByName("Body")

					req := httptest.NewRequest(http.MethodGet, "/", testx.ToReadCloser(encoder.NewJSON(), expected.Body))
					factory.On("CreateFromRequest", req).Return(encoder.NewJSON()).Once()

					err := setter.Body(valueOf, req)

					So(err, ShouldBeError, "cannot set value to type")
					So(actual.Body, ShouldBeNil)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when body type fails to unmarshal", func() {
					req := httptest.NewRequest(http.MethodGet, "/", ioutil.NopCloser(strings.NewReader("{")))
					factory.On("CreateFromRequest", req).Return(encoder.NewJSON()).Once()

					err := setter.Body(valueOf, req)

					So(err, ShouldBeError)
					So(err.Error(), ShouldStartWith, "decode application/json: testx.JSONGambit")
					So(actual.Body, ShouldBeNil)
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
		})
	})
}

func TestRequestHandlerSetter_Header(t *testing.T) {
	t.Parallel()

	Convey("Header", t, func() {
		factory := &encoder.FactoryMock{}

		bag := []interface{}{factory}

		setter := internal.NewRequestHandlerSetter(factory)

		expected := testx.GetTestStruct()
		actual := testx.TestStruct{}

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		jsonHeader, err := encoder.NewJSON().Encode(expected.HeaderJSON)
		So(err, ShouldBeNil)

		req.Header.Set("string", expected.HeaderString)
		req.Header.Set("int", "-2")
		req.Header.Set("uint", "2")
		req.Header.Set("float", ".2")
		req.Header.Set("bool", "false")
		req.Header.Set("json", string(jsonHeader))

		Convey("for type", func() {
			Convey("string", func() {
				Convey("when header is string", func() {
					valueOf := getFields(&actual, "HeaderString")

					err := setter.Header(valueOf, req, "string")

					So(err, ShouldBeNil)
					So(actual.HeaderString, ShouldResemble, expected.HeaderString)
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("int", func() {
				Convey("when header is int", func() {
					valueOf := getFields(&actual, "HeaderInt")

					err := setter.Header(valueOf, req, "int")

					So(err, ShouldBeNil)
					So(actual.HeaderInt, ShouldResemble, expected.HeaderInt)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when header not a valid number", func() {
					req.Header.Set("int", "NaN")
					valueOf := getFields(&actual, "HeaderInt")

					err := setter.Header(valueOf, req, "int")

					So(err, ShouldBeError, "strconv.ParseInt: parsing \"NaN\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("uint", func() {
				Convey("when header is int", func() {
					valueOf := getFields(&actual, "HeaderUInt")

					err := setter.Header(valueOf, req, "uint")

					So(err, ShouldBeNil)
					So(actual.HeaderUInt, ShouldResemble, expected.HeaderUInt)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when header is not a valid number", func() {
					req.Header.Set("uint", "NaN")
					valueOf := getFields(&actual, "HeaderUInt")

					err := setter.Header(valueOf, req, "uint")

					So(err, ShouldBeError, "strconv.ParseUint: parsing \"NaN\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("float", func() {
				Convey("when header is float", func() {
					valueOf := getFields(&actual, "HeaderFloat")

					err := setter.Header(valueOf, req, "float")

					So(err, ShouldBeNil)
					So(actual.HeaderFloat, ShouldResemble, expected.HeaderFloat)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when header is NaN", func() {
					req.Header.Set("float", "NaN")
					valueOf := getFields(&actual, "HeaderFloat")

					err := setter.Header(valueOf, req, "float")

					So(err, ShouldBeNil)
					So(math.IsNaN(actual.HeaderFloat), ShouldBeTrue)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when header is not a valid number", func() {
					req.Header.Set("float", "not a float")
					valueOf := getFields(&actual, "HeaderFloat")

					err := setter.Header(valueOf, req, "float")

					So(err, ShouldBeError, "strconv.ParseFloat: parsing \"not a float\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("bool", func() {
				Convey("when header is bool", func() {
					valueOf := getFields(&actual, "HeaderBool")

					err := setter.Header(valueOf, req, "bool")

					So(err, ShouldBeNil)
					So(actual.HeaderBool, ShouldResemble, expected.HeaderBool)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when header is not a bool", func() {
					req.Header.Set("bool", "maybe")
					valueOf := getFields(&actual, "HeaderBool")

					err := setter.Header(valueOf, req, "bool")

					So(err, ShouldBeError, "strconv.ParseBool: parsing \"maybe\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			mock.AssertExpectationsForObjects(t, bag...)
		})
	})
}

func TestRequestHandlerSetter_Path(t *testing.T) {
	t.Parallel()

	Convey("Path", t, func() {
		factory := &encoder.FactoryMock{}

		bag := []interface{}{factory}

		setter := internal.NewRequestHandlerSetter(factory)

		expected := testx.GetTestStruct()
		actual := testx.TestStruct{}

		cctx := &chi.Context{
			URLParams: chi.RouteParams{
				Keys: []string{
					"string",
					"int",
					"uint",
					"float",
					"bool",
				},
				Values: []string{
					expected.PathString,
					"-1",
					"1",
					".1",
					"true",
				},
			},
		}

		req := httptest.
			NewRequest(http.MethodGet, "/", nil).
			WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, cctx))

		Convey("for type", func() {
			Convey("string", func() {
				Convey("when header is string", func() {
					valueOf := getFields(&actual, "PathString")

					err := setter.Path(valueOf, req, "string")

					So(err, ShouldBeNil)
					So(actual.PathString, ShouldResemble, expected.PathString)
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("int", func() {
				Convey("when header is int", func() {
					valueOf := getFields(&actual, "PathInt")

					err := setter.Path(valueOf, req, "int")

					So(err, ShouldBeNil)
					So(actual.PathInt, ShouldResemble, expected.PathInt)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when header not a valid number", func() {
					cctx.URLParams.Values[1] = "NaN"
					req := httptest.
						NewRequest(http.MethodGet, "/", nil).
						WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, cctx))
					valueOf := getFields(&actual, "PathInt")

					err := setter.Path(valueOf, req, "int")

					So(err, ShouldBeError, "strconv.ParseInt: parsing \"NaN\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("uint", func() {
				Convey("when header is int", func() {
					valueOf := getFields(&actual, "PathUInt")

					err := setter.Path(valueOf, req, "uint")

					So(err, ShouldBeNil)
					So(actual.PathUInt, ShouldResemble, expected.PathUInt)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when header is not a valid number", func() {
					cctx.URLParams.Values[2] = "NaN"
					valueOf := getFields(&actual, "PathUInt")

					err := setter.Path(valueOf, req, "uint")

					So(err, ShouldBeError, "strconv.ParseUint: parsing \"NaN\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("float", func() {
				Convey("when header is float", func() {
					valueOf := getFields(&actual, "PathFloat")

					err := setter.Path(valueOf, req, "float")

					So(err, ShouldBeNil)
					So(actual.PathFloat, ShouldResemble, expected.PathFloat)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when header is NaN", func() {
					cctx.URLParams.Values[3] = "NaN"
					valueOf := getFields(&actual, "PathFloat")

					err := setter.Path(valueOf, req, "float")

					So(err, ShouldBeNil)
					So(math.IsNaN(actual.PathFloat), ShouldBeTrue)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when header is not a valid number", func() {
					cctx.URLParams.Values[3] = "not a float"
					valueOf := getFields(&actual, "PathFloat")

					err := setter.Path(valueOf, req, "float")

					So(err, ShouldBeError, "strconv.ParseFloat: parsing \"not a float\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("bool", func() {
				Convey("when header is bool", func() {
					valueOf := getFields(&actual, "PathBool")

					err := setter.Path(valueOf, req, "bool")

					So(err, ShouldBeNil)
					So(actual.PathBool, ShouldResemble, expected.PathBool)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when header is not a bool", func() {
					cctx.URLParams.Values[4] = "maybe"
					valueOf := getFields(&actual, "PathBool")

					err := setter.Path(valueOf, req, "bool")

					So(err, ShouldBeError, "strconv.ParseBool: parsing \"maybe\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			mock.AssertExpectationsForObjects(t, bag...)
		})
		Convey("should return nil when", func() {
			Convey("context is not set", func() {
				req := httptest.
					NewRequest(http.MethodGet, "/", nil).
					WithContext(context.Background())

				valueOf := getFields(&actual, "PathInt")

				err := setter.Path(valueOf, req, "int")

				So(err, ShouldBeNil)
				So(actual.PathInt, ShouldBeZeroValue)
				mock.AssertExpectationsForObjects(t, bag...)
			})
			Convey("context value is wrong type", func() {
				req := httptest.
					NewRequest(http.MethodGet, "/", nil).
					WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, map[string]string{}))

				valueOf := getFields(&actual, "PathInt")

				err := setter.Path(valueOf, req, "int")

				So(err, ShouldBeNil)
				So(actual.PathInt, ShouldBeZeroValue)
				mock.AssertExpectationsForObjects(t, bag...)
			})
			Convey("invalid path param is passed", func() {
				valueOf := getFields(&actual, "PathInt")

				err := setter.Path(valueOf, req, "InT")

				So(err, ShouldBeNil)
				So(actual.PathInt, ShouldBeZeroValue)
				mock.AssertExpectationsForObjects(t, bag...)
			})
		})
	})
}

func TestRequestHandlerSetter_Query(t *testing.T) {
	t.Parallel()

	Convey("Query", t, func() {
		factory := &encoder.FactoryMock{}

		bag := []interface{}{factory}

		setter := internal.NewRequestHandlerSetter(factory)

		expected := testx.GetTestStruct()
		actual := testx.TestStruct{}

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		jsonQuery, err := encoder.NewJSON().Encode(expected.QueryJSON)
		So(err, ShouldBeNil)

		q := req.URL.Query()
		q.Add("string", expected.QueryString)
		q.Add("int", "-3")
		q.Add("uint", "3")
		q.Add("float", ".3")
		q.Add("bool", "true")
		q.Add("json", string(jsonQuery))
		req.URL.RawQuery = q.Encode()

		Convey("for type", func() {
			Convey("string", func() {
				Convey("when Query is string", func() {
					valueOf := getFields(&actual, "QueryString")

					err := setter.Query(valueOf, req, "string")

					So(err, ShouldBeNil)
					So(actual.QueryString, ShouldResemble, expected.QueryString)
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("int", func() {
				Convey("when Query is int", func() {
					valueOf := getFields(&actual, "QueryInt")

					err := setter.Query(valueOf, req, "int")

					So(err, ShouldBeNil)
					So(actual.QueryInt, ShouldResemble, expected.QueryInt)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when Query not a valid number", func() {
					q.Set("int", "NaN")
					req.URL.RawQuery = q.Encode()

					valueOf := getFields(&actual, "QueryInt")

					err := setter.Query(valueOf, req, "int")

					So(err, ShouldBeError, "strconv.ParseInt: parsing \"NaN\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("uint", func() {
				Convey("when Query is int", func() {
					valueOf := getFields(&actual, "QueryUInt")

					err := setter.Query(valueOf, req, "uint")

					So(err, ShouldBeNil)
					So(actual.QueryUInt, ShouldResemble, expected.QueryUInt)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when Query is not a valid number", func() {
					q.Set("uint", "NaN")
					req.URL.RawQuery = q.Encode()

					valueOf := getFields(&actual, "QueryUInt")

					err := setter.Query(valueOf, req, "uint")

					So(err, ShouldBeError, "strconv.ParseUint: parsing \"NaN\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("float", func() {
				Convey("when Query is float", func() {
					valueOf := getFields(&actual, "QueryFloat")

					err := setter.Query(valueOf, req, "float")

					So(err, ShouldBeNil)
					So(actual.QueryFloat, ShouldResemble, expected.QueryFloat)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when Query is NaN", func() {
					q.Set("float", "NaN")
					req.URL.RawQuery = q.Encode()

					valueOf := getFields(&actual, "QueryFloat")

					err := setter.Query(valueOf, req, "float")

					So(err, ShouldBeNil)
					So(math.IsNaN(actual.QueryFloat), ShouldBeTrue)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when Query is not a valid number", func() {
					q.Set("float", "not a float")
					req.URL.RawQuery = q.Encode()

					valueOf := getFields(&actual, "QueryFloat")

					err := setter.Query(valueOf, req, "float")

					So(err, ShouldBeError, "strconv.ParseFloat: parsing \"not a float\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			Convey("bool", func() {
				Convey("when Query is bool", func() {
					valueOf := getFields(&actual, "QueryBool")

					err := setter.Query(valueOf, req, "bool")

					So(err, ShouldBeNil)
					So(actual.QueryBool, ShouldResemble, expected.QueryBool)
					mock.AssertExpectationsForObjects(t, bag...)
				})
				Convey("when Query is not a bool", func() {
					q.Set("bool", "maybe")
					req.URL.RawQuery = q.Encode()

					valueOf := getFields(&actual, "QueryBool")

					err := setter.Query(valueOf, req, "bool")

					So(err, ShouldBeError, "strconv.ParseBool: parsing \"maybe\": invalid syntax")
					mock.AssertExpectationsForObjects(t, bag...)
				})
			})
			mock.AssertExpectationsForObjects(t, bag...)
		})
		Convey("structure", func() {
			Convey("when query is encoded structure", func() {
				valueOf := getFields(&actual, "QueryJSON")

				err := setter.Query(valueOf, req, "json")

				So(err, ShouldBeNil)
				So(actual.QueryJSON, ShouldResemble, expected.QueryJSON)
				mock.AssertExpectationsForObjects(t, bag...)
			})
		})
	})
}

func getFields(t *testx.TestStruct, name string) reflect.Value {
	valueOf := reflect.ValueOf(t).Elem().FieldByName(name)

	return valueOf
}

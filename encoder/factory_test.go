package encoder_test

import (
	"net/http"
	"testing"

	"github.com/kevinanthony/gorps/v2/encoder"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewFactory(t *testing.T) {
	t.Parallel()

	Convey("NewFactory", t, func() {
		Convey("should return new factory", func() {
			f := encoder.NewFactory()

			So(f, ShouldImplement, (*encoder.Factory)(nil))
		})
	})
}

func TestFactoryMock_CreateFromResponse(t *testing.T) {
	t.Parallel()

	Convey("CreateFromResponse", t, func() {
		resp := &http.Response{
			Header: http.Header{},
		}
		factory := encoder.NewFactory()

		Convey("should return json encoder", func() {
			Convey("when content-type is empty", func() {
				actual := factory.CreateFromResponse(resp)

				So(actual, ShouldHaveSameTypeAs, encoder.NewJSON())
			})
			Convey("when content-type is application/json", func() {
				resp.Header.Add("content-type", encoder.ApplicationJSON)

				actual := factory.CreateFromResponse(resp)

				So(actual, ShouldHaveSameTypeAs, encoder.NewJSON())
			})
		})
		Convey("should return xml encoder", func() {
			Convey("when content-type is application/xml", func() {
				resp.Header.Add("content-type", encoder.ApplicationXML)

				actual := factory.CreateFromResponse(resp)

				So(actual, ShouldHaveSameTypeAs, encoder.NewXML())
			})
		})
	})
}

func TestFactory_CreateFromRequest(t *testing.T) {
	t.Parallel()

	Convey("CreateFromResponse", t, func() {
		resp := &http.Request{
			Header: http.Header{},
		}
		factory := encoder.NewFactory()

		Convey("should return json encoder", func() {
			Convey("when accept is empty", func() {
				actual := factory.CreateFromRequest(resp)

				So(actual, ShouldHaveSameTypeAs, encoder.NewJSON())
			})
			Convey("when accept is application/json", func() {
				resp.Header.Add("accept", encoder.ApplicationJSON)

				actual := factory.CreateFromRequest(resp)

				So(actual, ShouldHaveSameTypeAs, encoder.NewJSON())
			})
		})
		Convey("should return xml encoder", func() {
			Convey("when accept is application/xml", func() {
				resp.Header.Add("accept", encoder.ApplicationXML)

				actual := factory.CreateFromRequest(resp)

				So(actual, ShouldHaveSameTypeAs, encoder.NewXML())
			})
		})
	})
}

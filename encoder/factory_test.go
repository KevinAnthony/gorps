package encoder_test

import (
	"net/http"
	"testing"

	encoder2 "github.com/kevinanthony/gorps/encoder"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewFactory(t *testing.T) {
	t.Parallel()

	Convey("NewFactory", t, func() {
		Convey("should return new factory", func() {
			f := encoder2.NewFactory()

			So(f, ShouldImplement, (*encoder2.Factory)(nil))
		})
	})
}

func TestFactoryMock_Create(t *testing.T) {
	t.Parallel()

	Convey("CreateFromResponse", t, func() {
		resp := &http.Response{
			Header: http.Header{},
		}
		factory := encoder2.NewFactory()

		Convey("should return json encoder", func() {
			Convey("when content-type is empty", func() {
				actual := factory.CreateFromResponse(resp)

				So(actual, ShouldHaveSameTypeAs, encoder2.NewJSON())
			})
			Convey("when content-type is application/json", func() {
				resp.Header.Add("content-type", encoder2.ApplicationJSON)

				actual := factory.CreateFromResponse(resp)

				So(actual, ShouldHaveSameTypeAs, encoder2.NewJSON())
			})
		})
		Convey("should return xml encoder", func() {
			Convey("when content-type is application/xml", func() {
				resp.Header.Add("content-type", encoder2.ApplicationXML)

				actual := factory.CreateFromResponse(resp)

				So(actual, ShouldHaveSameTypeAs, encoder2.NewXML())
			})
		})
	})
}

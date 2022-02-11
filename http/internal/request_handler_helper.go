package internal

import (
	"net/http"
	"reflect"
)

func NewRequestHandlerHelper() RequestHandlerHelper {
	return &requestHandlerHelper{
		setter: NewRequestHandlerSetter(nil),
	}
}

type RequestHandlerHelper interface {
	Fill(r *http.Request, dst interface{}) error
}

type requestHandlerHelper struct {
	setter RequestHandlerSetter
}

func (h requestHandlerHelper) Fill(r *http.Request, dst interface{}) error {
	elementsValue := reflect.ValueOf(dst).Elem()
	elementsType := reflect.TypeOf(dst).Elem()

	for i := 0; i < elementsValue.NumField(); i++ {
		value := elementsValue.Field(i)
		typeOf := elementsType.Field(i)

		var err error
		if tag, found := typeOf.Tag.Lookup("header"); found {
			err = h.setter.Header(value, typeOf.Type, r, tag)
		}

		if tag, found := typeOf.Tag.Lookup("query"); found {
			err = h.setter.Query(value, typeOf.Type, r, tag)
		}

		if tag, found := typeOf.Tag.Lookup("path"); found {
			err = h.setter.Path(value, typeOf.Type, r, tag)
		}

		if _, found := typeOf.Tag.Lookup("body"); found {
			err = h.setter.Body(value, typeOf.Type, r)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

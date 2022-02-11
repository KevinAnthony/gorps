package internal

import (
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/kevinanthony/gorps/encoder"

	"github.com/go-chi/chi"
)

const (
	base10  = 10
	bitSize = 64
)

type RequestHandlerSetter interface {
	Body(reflect.Value, reflect.Type, *http.Request) error
	Header(reflect.Value, reflect.Type, *http.Request, string) error
	Path(reflect.Value, reflect.Type, *http.Request, string) error
	Query(reflect.Value, reflect.Type, *http.Request, string) error
}

type requestHandlerSetter struct {
	factory encoder.Factory
}

func NewRequestHandlerSetter(factory encoder.Factory) RequestHandlerSetter {
	if factory == nil {
		panic("encoder factory is required")
	}

	return &requestHandlerSetter{
		factory: factory,
	}
}

func (r requestHandlerSetter) Body(value reflect.Value, typeOf reflect.Type, req *http.Request) error {
	bts, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return r.setStruct(typeOf, value, req, bts)
}

func (r requestHandlerSetter) Header(value reflect.Value, typeOf reflect.Type, req *http.Request, header string) error {
	headerStr := req.Header.Get(header)

	return r.set(value, typeOf, req, headerStr)
}

func (r requestHandlerSetter) Path(value reflect.Value, typeOf reflect.Type,
	req *http.Request, pathParam string) error {
	val, ok := req.Context().Value(chi.RouteCtxKey).(*chi.Context)
	if !ok {
		return nil
	}

	for i := range val.URLParams.Keys {
		if val.URLParams.Keys[i] == pathParam {
			return r.set(value, typeOf, req, val.URLParams.Values[i])
		}
	}

	return nil
}

func (r requestHandlerSetter) Query(value reflect.Value, typeOf reflect.Type, req *http.Request, query string) error {
	queryStr := req.URL.Query().Get(query)

	return r.set(value, typeOf, req, queryStr)
}

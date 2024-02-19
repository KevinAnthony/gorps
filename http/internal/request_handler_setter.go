package internal

import (
	"io"
	"net/http"
	"reflect"

	"github.com/kevinanthony/gorps/v2/encoder"

	"github.com/go-chi/chi/v5"
)

const (
	base10  = 10
	bitSize = 64
)

type RequestHandlerSetter interface {
	Body(value reflect.Value, req *http.Request) error
	Header(value reflect.Value, req *http.Request, header string) error
	Path(value reflect.Value, req *http.Request, path string) error
	Query(value reflect.Value, req *http.Request, query string) error
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

func (r requestHandlerSetter) Body(value reflect.Value, req *http.Request) error {
	bts, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	return r.setStruct(value, r.factory.CreateFromRequest(req), bts)
}

func (r requestHandlerSetter) Header(value reflect.Value, req *http.Request, header string) error {
	headerStr := req.Header.Get(header)

	return r.set(value, headerStr)
}

func (r requestHandlerSetter) Path(value reflect.Value,
	req *http.Request, pathParam string,
) error {
	chiContext, ok := req.Context().Value(chi.RouteCtxKey).(*chi.Context)
	if !ok {
		return nil
	}

	str := chiContext.URLParam(pathParam)
	if len(str) == 0 {
		return nil
	}

	return r.set(value, str)
}

func (r requestHandlerSetter) Query(value reflect.Value, req *http.Request, query string) error {
	queryStr := req.URL.Query().Get(query)

	return r.set(value, queryStr)
}

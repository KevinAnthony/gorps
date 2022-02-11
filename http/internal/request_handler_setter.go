package internal

import (
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/kevinanthony/gorps/encoder"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
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

//nolint: funlen, cyclop // this is just a big switch, nothing complex
func (r requestHandlerSetter) set(value reflect.Value, typeOf reflect.Type, req *http.Request, str string) error {
	if len(str) == 0 {
		return nil
	}

	switch typeOf.Kind() {
	case reflect.Bool:
		return r.setBool(str, value)
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
		return r.setInt(str, value)
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		return r.setUint(str, value)
	case reflect.Float32:
	case reflect.Float64:
		return r.setFloat(str, value)
	case reflect.Array:
	case reflect.Slice:
		return errors.New("arrays/slices are currently unsupported")
	case reflect.String:
		return r.setString(str, value)
	case reflect.Struct:
		return r.setStruct(typeOf, value, req, []byte(str))
	case reflect.Chan:
		return errors.New("unsupported type: channel")
	case reflect.Complex128:
	case reflect.Complex64:
		return errors.New("unsupported type: complex")
	case reflect.Func:
		return errors.New("unsupported type: function")
	case reflect.Map:
		return errors.New("unsupported type: map")
	case reflect.Ptr:
		return errors.New("unsupported type: pointer")
	case reflect.Uintptr:
		return errors.New("unsupported type: int pointer")
	case reflect.UnsafePointer:
		return errors.New("unsupported type: unsafe pointer")
	case reflect.Interface:
		return errors.New("unsupported type: interface")
	case reflect.Invalid:
	default:
		return errors.New("unsupported type: unknown/invalid")
	}

	return nil
}

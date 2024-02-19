package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/kevinanthony/gorps/v2/encoder"
	"github.com/kevinanthony/gorps/v2/header"
	"github.com/kevinanthony/gorps/v2/http/internal"
)

type RequestHandlerFunc func(ctx context.Context, r *http.Request) (interface{}, error)

type RequestHandler interface {
	Handle(f RequestHandlerFunc) http.HandlerFunc
	MarshalAndVerify(r *http.Request, dst interface{}) error
}

type requestHandler struct {
	helper internal.RequestHandlerHelper
}

func NewRequestHandler(helper internal.RequestHandlerHelper) RequestHandler {
	if helper == nil {
		panic("request handler helper is required")
	}

	return &requestHandler{
		helper: helper,
	}
}

func NewRequestHandlerHelper() internal.RequestHandlerHelper {
	return internal.NewRequestHandlerHelper(
		internal.NewRequestHandlerSetter(
			encoder.NewFactory()))
}

func (rh requestHandler) MarshalAndVerify(r *http.Request, dst interface{}) error {
	return rh.helper.Fill(r, dst)
}

func (rh requestHandler) Handle(f RequestHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := f(r.Context(), r)
		if err != nil {
			write(w, r, http.StatusBadRequest, err.Error())

			return
		}

		// TODO create interface to get StatusCode
		write(w, r, http.StatusOK, resp)
	}
}

func write(w http.ResponseWriter, r *http.Request, statusCode int, src interface{}) {
	enc := encoder.NewFactory().CreateFromRequest(r)

	bts, err := enc.Encode(src)
	if err != nil {
		// TODO standard error object
		w.WriteHeader(http.StatusBadRequest)

		_, _ = w.Write([]byte("could not encode data"))

		return
	}

	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Add(header.ContentType, enc.GetMime())
	w.Header().Add(header.ContentLength, strconv.Itoa(len(bts)))

	w.WriteHeader(statusCode)

	_, _ = w.Write(bts)
}

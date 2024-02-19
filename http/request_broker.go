package http

import (
	"context"
	"fmt"
	"io"
	native "net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

//nolint:interfacebloat
//go:generate mockery --name=RequestBroker --structname=RequestBrokerMock --filename=request_broker_mock.go --inpackage
type RequestBroker interface {
	DoAndUnmarshal(ctx context.Context, v interface{}) error
	Do(ctx context.Context) ([]byte, error)

	Post() RequestBroker
	Get() RequestBroker
	Put() RequestBroker
	Delete() RequestBroker

	URL(url string, v ...any) RequestBroker
	Query(key string, value string) RequestBroker
	Header(key string, value string) RequestBroker
	Body(body string) RequestBroker

	CreateRequest(ctx context.Context) (*native.Request, error)
}

type requestBroker struct {
	err    error
	client Client
	url    string

	method  MethodType
	headers map[string]string
	query   map[string]string
	body    string
}

func NewRequest(client Client) RequestBroker {
	r := &requestBroker{
		method:  MethodGet,
		headers: map[string]string{},
		query:   map[string]string{},
	}

	if client == nil {
		r.setErrStr("native client is nil")
	}

	r.client = client

	return r
}

func (r *requestBroker) Post() RequestBroker {
	r.method = MethodPost

	return r
}

func (r *requestBroker) Get() RequestBroker {
	r.method = MethodGet

	return r
}

func (r *requestBroker) Put() RequestBroker {
	r.method = MethodPut

	return r
}

func (r *requestBroker) Delete() RequestBroker {
	r.method = MethodDelete

	return r
}

func (r *requestBroker) URL(s string, v ...any) RequestBroker {
	r.url = fmt.Sprintf(s, v...)

	return r
}

func (r *requestBroker) Query(pattern, value string) RequestBroker {
	r.query[pattern] = value

	return r
}

func (r *requestBroker) Header(header, value string) RequestBroker {
	r.headers[header] = value

	return r
}

func (r *requestBroker) Body(bytes string) RequestBroker {
	r.body = bytes

	return r
}

func (r *requestBroker) DoAndUnmarshal(ctx context.Context, out interface{}) error {
	req, err := r.CreateRequest(ctx)
	if err != nil {
		return err
	}

	return r.client.DoAndUnmarshal(req, out)
}

func (r *requestBroker) Do(ctx context.Context) ([]byte, error) {
	req, err := r.CreateRequest(ctx)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req)
}

func (r *requestBroker) setErrStr(s string) {
	if r.err != nil {
		r.err = errors.New(s)
	}
}

func (r *requestBroker) CreateRequest(ctx context.Context) (*native.Request, error) {
	if r.err != nil {
		return nil, r.err
	}

	url, err := url.Parse(r.url)
	if err != nil {
		return nil, err
	}

	req := &native.Request{
		Method:     string(r.method),
		URL:        url,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(native.Header),
		Body:       nil,
		Host:       url.Host,
	}

	if len(r.body) > 0 {
		req.ContentLength = int64(len(r.body))
		req.Body = io.NopCloser(strings.NewReader(r.body))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(r.body)), nil
		}
	}

	for k, v := range r.headers {
		req.Header.Add(k, v)
	}

	query := req.URL.Query()

	for k, v := range r.query {
		query.Set(k, v)
	}

	req.URL.RawQuery = query.Encode()

	return req.WithContext(ctx), nil
}

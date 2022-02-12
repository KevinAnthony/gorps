package http

import (
	"context"
	"fmt"
	native "net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type RequestBroker interface {
	Go(ctx context.Context, v interface{}) error

	Post() RequestBroker
	Get() RequestBroker
	Put() RequestBroker
	Delete() RequestBroker

	Domain(string) RequestBroker
	Path(string) RequestBroker
	Parameter(string, string) RequestBroker

	Header(string, string) RequestBroker
}

type requestBroker struct {
	err    error
	client Client
	domain string
	path   string
	method MethodType

	parameters map[string]string
	headers    map[string]string
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

func (r *requestBroker) Domain(s string) RequestBroker {
	r.domain = s

	return r
}

func (r *requestBroker) Path(s string) RequestBroker {
	r.path = s

	return r
}

func (r *requestBroker) Parameter(pattern, value string) RequestBroker {
	r.parameters[pattern] = value

	return r
}

func (r *requestBroker) Header(header, value string) RequestBroker {
	r.headers[header] = value

	return r
}

func NewRequest(client Client) RequestBroker {
	r := &requestBroker{
		method:     MethodGet,
		parameters: map[string]string{},
		headers:    map[string]string{},
	}

	if client == nil {
		r.setErrStr("native client is nil")
	}

	r.client = client

	return r
}

func (r *requestBroker) Go(ctx context.Context, out interface{}) error {
	if r.err != nil {
		return r.err
	}

	for k, v := range r.parameters {
		if !strings.Contains(r.path, k) {
			return fmt.Errorf("missing parameter %s in path", k)
		}

		r.path = strings.ReplaceAll(r.path, k, v)
	}

	req := &native.Request{
		Method: string(r.method),
		URL: &url.URL{
			Scheme: "https",
			Host:   r.domain,
			Path:   r.path,
		},
		Proto:  "https",
		Header: native.Header{},
		Body:   nil,
	}

	for k, v := range r.headers {
		req.Header.Add(k, v)
	}

	req = req.WithContext(ctx)

	return r.client.Do(req, out)
}

func (r *requestBroker) setErrStr(s string) {
	if r.err != nil {
		r.err = errors.New(s)
	}
}

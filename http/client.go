package http

import (
	"io/ioutil"
	native "net/http"
	"strings"

	"github.com/kevinanthony/gorps/encoder"

	"github.com/pkg/errors"
)

var (
	errBadRequest = errors.New("bad requestBroker")
)

type Client interface {
	Do(req *native.Request, dst interface{}) error
}

type Native interface {
	Do(req *native.Request) (*native.Response, error)
}

type client struct {
	encFactory encoder.Factory
	client     Native
}

func NewNativeClient() Native {
	return &native.Client{}
}

func NewClient(nativeClient Native, enc encoder.Factory) Client {
	if nativeClient == nil {
		panic("http client is required")
	}

	if enc == nil {
		panic("encoding factory is required")
	}

	return &client{
		encFactory: enc,
		client:     nativeClient,
	}
}

func (c client) Do(req *native.Request, dst interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= native.StatusBadRequest {
		return errors.Wrapf(errBadRequest, "%d: %s",
			resp.StatusCode, strings.Trim(string(bts), "\""))
	}

	if len(bts) == 0 {
		return nil
	}

	return c.encFactory.CreateFromResponse(resp).Decode(bts, dst)
}

package http

import (
	"io"
	native "net/http"
	"strings"

	"github.com/kevinanthony/gorps/v2/encoder"

	"github.com/pkg/errors"
)

var errBadRequest = errors.New("bad requestBroker")

//go:generate mockery --srcpkg=io --name=ReadCloser --structname=BodyMock --filename=body_mock.go --output . --outpkg=http

//go:generate mockery --name=Client --structname=ClientMock --filename=client_mock.go --inpackage
type Client interface {
	DoAndUnmarshal(req *native.Request, v interface{}) error
	Do(req *native.Request) (io.Reader, error)
}

//go:generate mockery --name=Native --structname=NativeMock --filename=native_mock.go --inpackage
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

func (c client) DoAndUnmarshal(req *native.Request, dst interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	bts, err := io.ReadAll(resp.Body)
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

func (c client) Do(req *native.Request) (io.Reader, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode >= native.StatusBadRequest {
		bts, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.Wrapf(errBadRequest, "%d: %s",
			resp.StatusCode, strings.Trim(string(bts), "\""))
	}

	return resp.Body, nil
}

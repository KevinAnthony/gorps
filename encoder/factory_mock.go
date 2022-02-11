package encoder

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

var _ Factory = (*FactoryMock)(nil)

type FactoryMock struct {
	mock.Mock
}

func (f *FactoryMock) CreateFromRequest(req *http.Request) Encoder {
	args := f.Called(req)

	var enc Encoder
	if item := args.Get(0); item != nil {
		enc = item.(Encoder)
	}

	return enc
}

func (f *FactoryMock) FromMime(mediaType string) Encoder {
	args := f.Called(mediaType)

	var enc Encoder
	if item := args.Get(0); item != nil {
		enc = item.(Encoder)
	}

	return enc
}

func (f *FactoryMock) CreateFromResponse(resp *http.Response) Encoder {
	args := f.Called(resp)

	var enc Encoder
	if item := args.Get(0); item != nil {
		enc = item.(Encoder)
	}

	return enc
}

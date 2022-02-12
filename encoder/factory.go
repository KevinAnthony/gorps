package encoder

import (
	"mime"
	"net/http"

	"github.com/kevinanthony/gorps/header"
)

type Factory interface {
	CreateFromResponse(resp *http.Response) Encoder
	CreateFromRequest(req *http.Request) Encoder
	FromMime(mediaType string) Encoder
}

type factory struct{}

func NewFactory() Factory {
	return factory{}
}

func (f factory) CreateFromResponse(resp *http.Response) Encoder {
	mediaType, _, _ := mime.ParseMediaType(resp.Header.Get(header.ContentType))

	return f.FromMime(mediaType)
}

func (f factory) CreateFromRequest(req *http.Request) Encoder {
	mediaType, _, _ := mime.ParseMediaType(req.Header.Get(header.Accept))

	return f.FromMime(mediaType)
}

func (f factory) FromMime(mediaType string) Encoder {
	switch mediaType {
	case ApplicationXML:
		return NewXML()
	case ApplicationJSON:
		return NewJSON()
	default:
		return NewJSON()
	}
}

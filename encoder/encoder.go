package encoder

type (
	AcceptType = string
)

const (
	ApplicationJSON AcceptType = "application/json"
	ApplicationXML  AcceptType = "application/xml"
	TextXML         AcceptType = "text/xml"
)

type Encoder interface {
	Encode(data interface{}) ([]byte, error)
	Decode(data []byte, dst interface{}) error
	GetMime() string
}

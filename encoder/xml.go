package encoder

import "encoding/xml"

type xmlEncoder struct{}

func NewXML() Encoder {
	return xmlEncoder{}
}

func (x xmlEncoder) Encode(data any) ([]byte, error) {
	return xml.Marshal(data)
}

func (x xmlEncoder) Decode(data []byte, dst any) error {
	return xml.Unmarshal(data, dst)
}

func (x xmlEncoder) GetMime() string {
	return ApplicationXML
}

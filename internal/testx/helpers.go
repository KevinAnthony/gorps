package testx

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/kevinanthony/gorps/encoder"
)

func ToReadCloser(enc encoder.Encoder, in interface{}) io.ReadCloser {
	b, err := enc.Encode(in)
	if err != nil {
		panic(err)
	}

	return ioutil.NopCloser(bytes.NewReader(b))
}

func GetTestStruct() TestStruct {
	return TestStruct{
		PathString:    "Path",
		PathStringP:   s2p("Path"),
		PathInt:       -1,
		PathIntP:      i2p(-1),
		PathUInt:      1,
		PathUIntP:     u2p(1),
		PathFloat:     .1,
		PathFloatP:    f2p(.1),
		PathBool:      true,
		PathBoolP:     b2p(true),
		HeaderString:  "header",
		HeaderStringP: s2p("header"),
		HeaderInt:     -2,
		HeaderIntP:    i2p(-2),
		HeaderUInt:    2,
		HeaderUIntP:   u2p(2),
		HeaderFloat:   .2,
		HeaderFloatP:  f2p(.2),
		HeaderBool:    false,
		HeaderBoolP:   b2p(false),
		HeaderJSON:    getJSONGambit(),
		QueryString:   "query",
		QueryStringP:  s2p("query"),
		QueryInt:      -3,
		QueryIntP:     i2p(-3),
		QueryUInt:     3,
		QueryUIntP:    u2p(3),
		QueryFloat:    .3,
		QueryFloatP:   f2p(.3),
		QueryBool:     true,
		QueryBoolP:    b2p(true),
		QueryJSON:     getJSONGambit(),
		Body:          getJSONGambitPtr(),
	}
}

func getJSONGambit() JSONGambit {
	return JSONGambit{
		String:  "json",
		StringP: s2p("json"),
		Int:     -4,
		IntP:    i2p(-4),
		UInt:    4,
		UIntP:   u2p(4),
		Float:   44.4,
		FloatP:  f2p(44.4),
		Bool:    false,
		BoolP:   b2p(true),
	}
}
func getJSONGambitPtr() *JSONGambit {
	return &JSONGambit{
		String:  "json",
		StringP: s2p("json"),
		Int:     -4,
		IntP:    i2p(-4),
		UInt:    4,
		UIntP:   u2p(4),
		Float:   44.4,
		FloatP:  f2p(44.4),
		Bool:    false,
		BoolP:   b2p(true),
	}
}

type TestStruct struct {
	PathString  string   `path:"string"`
	PathStringP *string  `path:"string"`
	PathInt     int      `path:"int"`
	PathIntP    *int     `path:"int"`
	PathUInt    uint     `path:"uint"`
	PathUIntP   *uint    `path:"uint"`
	PathFloat   float64  `path:"float"`
	PathFloatP  *float64 `path:"float"`
	PathBool    bool     `path:"bool"`
	PathBoolP   *bool    `path:"bool"`

	HeaderString  string     `header:"string"`
	HeaderStringP *string    `header:"string"`
	HeaderInt     int        `header:"int"`
	HeaderIntP    *int       `header:"int"`
	HeaderUInt    uint       `header:"uint"`
	HeaderUIntP   *uint      `header:"uint"`
	HeaderFloat   float64    `header:"float"`
	HeaderFloatP  *float64   `header:"float"`
	HeaderBool    bool       `header:"bool"`
	HeaderBoolP   *bool      `header:"bool"`
	HeaderJSON    JSONGambit `header:"json"`

	QueryString  string     `query:"string"`
	QueryStringP *string    `query:"string"`
	QueryInt     int        `query:"int"`
	QueryIntP    *int       `query:"int"`
	QueryUInt    uint       `query:"uint"`
	QueryUIntP   *uint      `query:"uint"`
	QueryFloat   float64    `query:"float"`
	QueryFloatP  *float64   `query:"float"`
	QueryBool    bool       `query:"bool"`
	QueryBoolP   *bool      `query:"bool"`
	QueryJSON    JSONGambit `query:"json"`

	Body *JSONGambit `body:"request"`
}

type JSONGambit struct {
	String  string   `json:"string"`
	StringP *string  `json:"string_2"`
	Int     int      `json:"int"`
	IntP    *int     `json:"int_2"`
	UInt    uint     `json:"uint"`
	UIntP   *uint    `json:"uint_2"`
	Float   float64  `json:"float"`
	FloatP  *float64 `json:"float_2"`
	Bool    bool     `json:"bool"`
	BoolP   *bool    `json:"bool_2"`
}

func s2p(s string) *string {
	return &s
}

func i2p(i int) *int {
	return &i
}

func u2p(u uint) *uint {
	return &u
}

func b2p(b bool) *bool {
	return &b
}

func f2p(f float64) *float64 {
	return &f
}

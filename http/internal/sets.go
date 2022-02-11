package internal

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

func (r requestHandlerSetter) setStruct(typeOf reflect.Type, value reflect.Value, req *http.Request, bts []byte) error {
	if !value.IsValid() {
		return errors.New("bad body value")
	}

	if !value.Elem().CanSet() {
		return errors.New("cannot set value to type")
	}

	dst := reflect.New(typeOf.Elem()).Interface()
	enc := r.factory.CreateFromRequest(req)

	if err := enc.Decode(bts, &dst); err != nil {
		return errors.Wrapf(err, "decode %s", enc.GetMime())
	}

	v := reflect.ValueOf(dst)
	value.Elem().Set(v)

	return nil
}

func (r requestHandlerSetter) setBool(str string, value reflect.Value) error {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return err
	}

	value.SetBool(b)

	return nil
}

func (r requestHandlerSetter) setInt(str string, value reflect.Value) error {
	i, err := strconv.ParseInt(str, base10, bitSize)
	if err != nil {
		return err
	}

	value.SetInt(i)

	return nil
}

func (r requestHandlerSetter) setString(str string, value reflect.Value) error {
	value.SetString(str)

	return nil
}

func (r requestHandlerSetter) setFloat(str string, value reflect.Value) error {
	i, err := strconv.ParseFloat(str, bitSize)
	if err != nil {
		return err
	}

	value.SetFloat(i)

	return nil
}

func (r requestHandlerSetter) setUint(str string, value reflect.Value) error {
	i, err := strconv.ParseUint(str, base10, bitSize)
	if err != nil {
		return err
	}

	value.SetUint(i)

	return nil
}

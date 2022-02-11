package internal

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

//nolint: funlen, cyclop // this is just a big switch, nothing complex
func (r requestHandlerSetter) set(value reflect.Value, typeOf reflect.Type, req *http.Request, str string) error {
	if len(str) == 0 {
		return nil
	}

	switch typeOf.Kind() {
	case reflect.Bool:
		return r.setBool(str, value)
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return r.setInt(str, value)
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return r.setUint(str, value)
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return r.setFloat(str, value)
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return errors.New("arrays/slices are currently unsupported")
	case reflect.String:
		return r.setString(str, value)
	case reflect.Struct:
		return r.setStruct(typeOf, value, req, []byte(str))
	case reflect.Chan:
		return errors.New("unsupported type: channel")
	case reflect.Complex128:
		fallthrough
	case reflect.Complex64:
		return errors.New("unsupported type: complex")
	case reflect.Func:
		return errors.New("unsupported type: function")
	case reflect.Map:
		return errors.New("unsupported type: map")
	case reflect.Ptr:
		return errors.New("unsupported type: pointer")
	case reflect.Uintptr:
		return errors.New("unsupported type: int pointer")
	case reflect.UnsafePointer:
		return errors.New("unsupported type: unsafe pointer")
	case reflect.Interface:
		return errors.New("unsupported type: interface")
	case reflect.Invalid:
	default:
		return errors.New("unsupported type: unknown/invalid")
	}

	return nil
}

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

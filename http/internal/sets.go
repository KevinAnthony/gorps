package internal

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/kevinanthony/gorps/encoder"

	"github.com/pkg/errors"
)

//nolint: cyclop // this is just a big switch, nothing complex
func (r requestHandlerSetter) set(value reflect.Value, str string) error {
	if len(str) == 0 {
		return nil
	}

	switch value.Elem().Interface().(type) {
	case int, int8, int16, int32, int64:
		return r.setInt(value, str)
	case string:
		return r.setString(value, str)
	case float32, float64:
		return r.setFloat(value, str)
	case uint, uint8, uint16, uint32, uint64:
		return r.setUint(value, str)
	case bool:
		return r.setBool(value, str)
	default:
		//nolint: exhaustive
		switch value.Kind() {
		case reflect.Struct, reflect.Map:
			return r.setStruct(value, encoder.NewJSON(), []byte(str))
		case reflect.Ptr:
			if value.Type().Elem().Kind() == reflect.Struct {
				return r.setStruct(value, encoder.NewJSON(), []byte(str))
			}
		default:
			return fmt.Errorf("unsupported kind: %s", value.Kind())
		}

		return fmt.Errorf("unsupported type: %T", value.Interface())
	}
}

func (r requestHandlerSetter) setStruct(value reflect.Value, enc encoder.Encoder, bts []byte) error {
	if !value.IsValid() {
		return errors.New("bad body value")
	}

	if !value.Elem().CanSet() {
		return errors.New("cannot set value to type")
	}

	isPtr := value.Elem().Type().Kind() == reflect.Ptr

	typeOf := value.Elem().Type()
	if isPtr {
		typeOf = value.Elem().Type().Elem()
	}

	dst := reflect.New(typeOf).Interface()

	if err := enc.Decode(bts, &dst); err != nil {
		return errors.Wrapf(err, "decode %s", enc.GetMime())
	}

	dstValue := reflect.ValueOf(dst)
	if isPtr {
		value.Elem().Set(dstValue)

		return nil
	}

	value.Elem().Set(dstValue.Elem())

	return nil
}

func (r requestHandlerSetter) setBool(value reflect.Value, str string) error {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return err
	}

	value.Elem().SetBool(b)

	return nil
}

func (r requestHandlerSetter) setInt(value reflect.Value, str string) error {
	i, err := strconv.ParseInt(str, base10, bitSize)
	if err != nil {
		return err
	}

	value.Elem().SetInt(i)

	return nil
}

func (r requestHandlerSetter) setString(value reflect.Value, str string) error {
	value.Elem().SetString(str)

	return nil
}

func (r requestHandlerSetter) setFloat(value reflect.Value, str string) error {
	i, err := strconv.ParseFloat(str, bitSize)
	if err != nil {
		return err
	}

	value.Elem().SetFloat(i)

	return nil
}

func (r requestHandlerSetter) setUint(value reflect.Value, str string) error {
	i, err := strconv.ParseUint(str, base10, bitSize)
	if err != nil {
		return err
	}

	value.Elem().SetUint(i)

	return nil
}

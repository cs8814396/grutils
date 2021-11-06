package fasthttp

// decode query string into struct
// borrow from https://raw.githubusercontent.com/labstack/echo/master/bind.go

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	tag = "json" // we use json tag to decode
)

// QueryString2Struct QueryString2Struct
func QueryString2Struct(query string, ptr interface{}) error {
	values, e := url.ParseQuery(query)
	if e != nil {
		return e
	}

	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()

	// !struct
	if typ.Kind() != reflect.Struct {
		return errors.New("binding element must be a struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)

		structFieldKind := structField.Kind()
		tagStr := typeField.Tag.Get(tag)
		inputFieldName := strings.Split(tagStr, ",")[0]

		if len(inputFieldName) > 0 {
			inputValue, exists := values[inputFieldName]
			if exists {
				numElems := len(inputValue)
				if structFieldKind == reflect.Slice && numElems > 0 {
					sliceOf := structField.Type().Elem().Kind()
					slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
					for j := 0; j < numElems; j++ {
						if err := setWithProperType(sliceOf, inputValue[j], slice.Index(j)); err != nil {
							return err
						}
					}
					val.Field(i).Set(slice)
				} else if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
					return err

				}
			}
		}

	}

	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Ptr:
		return setWithProperType(structField.Elem().Kind(), val, structField.Elem())
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("unknown type")
	}
	return nil
}

func setIntField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

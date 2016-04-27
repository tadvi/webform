package webform

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

// If duplicate field name if located in nested structs - panic.
var PanicOnDuplicateName = true

// Decode from map[string]string into the spec.
func Decode(dst interface{}, m map[string][]string) error {
	v := reflect.ValueOf(dst)
	mustBe(v, reflect.Ptr)
	v = v.Elem()
	return decode(v, m)
}

// decode from map[string]string into the spec.
func decode(v reflect.Value, m map[string][]string) error {
	mustBe(v, reflect.Struct)
	t := v.Type()
	done := make(map[string]bool)

	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i)
		tField := t.Field(i)
		// ignore not settable fields
		if !vField.CanSet() {
			continue
		}

		if vField.Kind() == reflect.Struct {
			// embedded struct
			if err := decode(vField, m); err != nil {
				return err
			}
			continue
		}

		if vField.Kind() == reflect.Ptr {
			if vField.IsNil() {
				// create instance of the struct if it is a nil pointer
				vField.Set(reflect.New(vField.Type().Elem()))
			}
			if vField.Elem().Kind() == reflect.Struct {
				// embedded struct pointer
				if err := decode(vField.Elem(), m); err != nil {
					return err
				}
			}
			continue
		}

		name := tField.Tag.Get("webform")
		if name == "-" {
			// ignore '-' tags
			continue
		}
		if name == "" {
			name = tField.Name
		}

		item, ok := m[name]
		if !ok || len(item) == 0 {
			// form did not have this field
			continue
		}
		// at this point supports only one value coming from form.
		value := item[0]

		// handle duplicates
		if ok := done[name]; ok {
			if PanicOnDuplicateName {
				panic("duplicate field name: " + name)
			} else {
				// ignore duplicates but do not set them
				continue
			}
		}
		done[name] = true

		switch vField.Kind() {
		case reflect.String:
			vField.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			intValue, err := strconv.ParseInt(value, 0, vField.Type().Bits())
			if err != nil {
				return fmt.Errorf("Field: %s, Type: %v, Value: %v", name, tField, value)
			}
			vField.SetInt(intValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("Field: %s, Type: %v, Value: %v", name, tField, value)
			}
			vField.SetBool(boolValue)
		case reflect.Float32, reflect.Float64:
			floatValue, err := strconv.ParseFloat(value, vField.Type().Bits())
			if err != nil {
				return fmt.Errorf("Field: %s, Type: %v, Value: %v", name, tField, value)
			}
			vField.SetFloat(floatValue)
		}

	}
	return nil
}

// -- helpers & utilities --

type kinder interface {
	Kind() reflect.Kind
}

// methodName is returns the caller of the function calling methodName
func methodName() string {
	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}
	return f.Name()
}

// mustBe checks a value against a kind, panicing with a reflect.ValueError
// if the kind isn't that which is required.
func mustBe(v kinder, expected reflect.Kind) {
	k := v.Kind()
	if k != expected {
		panic(&reflect.ValueError{Method: methodName(), Kind: k})
	}
}

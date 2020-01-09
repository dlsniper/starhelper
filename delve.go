// The functions defined in this file are taken from the Delve project
// and are under the MIT license of the project.

package starhelper

import (
	"fmt"
	"reflect"

	"go.starlark.net/starlark"
)

// DecorateError adds additional information to the error making it easier to debug
func DecorateError(thread *starlark.Thread, err error) error {
	if err == nil {
		return nil
	}
	pos := thread.CallFrame(1).Pos
	if pos.Col > 0 {
		return fmt.Errorf("%s:%d:%d: %v", pos.Filename(), pos.Line, pos.Col, err)
	}
	return fmt.Errorf("%s:%d: %v", pos.Filename(), pos.Line, err)
}

// ToStarlarkValue handles the conversion of a Go data value to a Starlark value
func ToStarlarkValue(v interface{}) starlark.Value {
	switch v := v.(type) {
	case uint8:
		return starlark.MakeUint64(uint64(v))
	case uint16:
		return starlark.MakeUint64(uint64(v))
	case uint32:
		return starlark.MakeUint64(uint64(v))
	case uint64:
		return starlark.MakeUint64(v)
	case uintptr:
		return starlark.MakeUint64(uint64(v))
	case uint:
		return starlark.MakeUint64(uint64(v))
	case int8:
		return starlark.MakeInt64(int64(v))
	case int16:
		return starlark.MakeInt64(int64(v))
	case int32:
		return starlark.MakeInt64(int64(v))
	case int64:
		return starlark.MakeInt64(v)
	case int:
		return starlark.MakeInt64(int64(v))
	case string:
		return starlark.String(v)
	case map[string]uint64:
		var r starlark.Dict
		for k, v := range v {
			_ = r.SetKey(starlark.String(k), starlark.MakeUint64(v))
		}
		return &r
	case map[string]string:
		var r starlark.Dict
		for k, v := range v {
			_ = r.SetKey(starlark.String(k), starlark.String(v))
		}
		return &r
	case nil:
		return starlark.None
	case error:
		return starlark.String(v.Error())
	default:
		vval := reflect.ValueOf(v)
		switch vval.Type().Kind() {
		case reflect.Ptr:
			if vval.IsNil() {
				return starlark.None
			}
			return starlark.None
		}
		return starlark.String(fmt.Sprintf("%v", v))
	}
}

// unmarshalStarlarkValue unmarshals a starlark.Value 'val' into a Go variable 'dst'.
// This works similarly to encoding/json.Unmarshal and similar functions,
// but instead of getting its input from a byte buffer, it uses a
// starlark.Value.
func UnmarshalStarlarkValue(val starlark.Value, dst interface{}, path string) error {
	return unmarshalStarlarkValueIntl(val, reflect.ValueOf(dst), path)
}

func unmarshalStarlarkValueIntl(val starlark.Value, dst reflect.Value, path string) (err error) {
	defer func() {
		// catches reflect panics
		ierr := recover()
		if ierr != nil {
			err = fmt.Errorf("error setting argument %q to %s: %v", path, val, ierr)
		}
	}()

	convErr := func(args ...string) error {
		if len(args) > 0 {
			return fmt.Errorf("error setting argument %q: can not convert %s to %s: %s", path, val, dst.Type().String(), args[0])
		}
		return fmt.Errorf("error setting argument %q: can not convert %s to %s", path, val, dst.Type().String())
	}

	if _, isnone := val.(starlark.NoneType); isnone {
		return nil
	}

	for dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = dst.Elem()
	}

	switch val := val.(type) {
	case starlark.Bool:
		dst.SetBool(bool(val))
	case starlark.Int:
		switch dst.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, ok := val.Uint64()
			if !ok {
				return convErr()
			}
			dst.SetUint(n)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, ok := val.Int64()
			if !ok {
				return convErr()
			}
			dst.SetInt(n)
		default:
			return convErr()
		}
	case starlark.Float:
		dst.SetFloat(float64(val))
	case starlark.String:
		dst.SetString(string(val))
	case *starlark.List:
		if dst.Kind() != reflect.Slice {
			return convErr()
		}
		r := []reflect.Value{}
		for i := 0; i < val.Len(); i++ {
			cur := reflect.New(dst.Type().Elem())
			err := unmarshalStarlarkValueIntl(val.Index(i), cur, path)
			if err != nil {
				return err
			}
			r = append(r, cur)
		}
	case *starlark.Dict:
		if dst.Kind() != reflect.Struct {
			return convErr()
		}
		for _, k := range val.Keys() {
			if _, ok := k.(starlark.String); !ok {
				return convErr(fmt.Sprintf("non-string key %q", k.String()))
			}
			fieldName := string(k.(starlark.String))
			dstfield := dst.FieldByName(fieldName)
			if dstfield == (reflect.Value{}) {
				return convErr(fmt.Sprintf("unknown field %s", fieldName))
			}
			valfield, _, _ := val.Get(starlark.String(fieldName))
			err := unmarshalStarlarkValueIntl(valfield, dstfield, path+"."+fieldName)
			if err != nil {
				return err
			}
		}
	default:
		return convErr()
	}
	return nil
}

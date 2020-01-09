// The functions defined in this file are taken from the Delve project
// and are under the MIT license of the project.

package starhelper

import (
	"fmt"
	"reflect"

	"go.starlark.net/starlark"
)

// Decorates the error using additional information about the error itself
func (env *App) decorateError(thread *starlark.Thread, err error) error {
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

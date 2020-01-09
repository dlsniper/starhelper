package starhelper

import (
	"errors"

	"go.starlark.net/starlark"
)

var neaError = errors.New("not enough arguments for function")

// ArgBool handles the conversion of function call arguments from Starlark bool to Go bool
func ArgBool(thread *starlark.Thread, args starlark.Tuple, idx int) (bool, error) {
	var res bool
	if len(args) < idx {
		return false, DecorateError(thread, neaError)
	}

	err := UnmarshalStarlarkValue(args.Index(idx), &res, "")
	if err != nil {
		return false, nil
	}
	return res, nil
}

// ArgInt handles the conversion of function call arguments from Starlark int to Go int
func ArgInt(thread *starlark.Thread, args starlark.Tuple, idx int) (int, error) {
	var res int
	if len(args) < idx {
		return 0, DecorateError(thread, neaError)
	}

	err := UnmarshalStarlarkValue(args.Index(idx), &res, "")
	if err != nil {
		return 0, nil
	}
	return res, nil
}

// ArgFloat handles the conversion of function call arguments from Starlark float to Go float64
func ArgFloat(thread *starlark.Thread, args starlark.Tuple, idx int) (float64, error) {
	var res float64
	if len(args) < idx {
		return 0, DecorateError(thread, neaError)
	}

	err := UnmarshalStarlarkValue(args.Index(idx), &res, "")
	if err != nil {
		return 0, nil
	}
	return res, nil
}

// ArgString handles the conversion of function call arguments from Starlark string to Go string
func ArgString(thread *starlark.Thread, args starlark.Tuple, idx int) (string, error) {
	var res string
	if len(args) < idx {
		return "", DecorateError(thread, neaError)
	}

	err := UnmarshalStarlarkValue(args.Index(idx), &res, "")
	if err != nil {
		return "", nil
	}
	return res, nil
}

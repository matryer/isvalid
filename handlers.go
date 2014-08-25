package is

import (
	"errors"
	"reflect"
	"strings"
)

var defaultHandlers = map[string]HandlerFunc{
	"required": func(v interface{}) (interface{}, error) {
		// ensures value is present

		var ok bool = true

		if v == nil {
			ok = false
		} else {
			switch val := v.(type) {
			case string:
				ok = len(val) > 0
			}
		}

		if !ok {
			return nil, errors.New("is required")
		}

		return v, nil
	},
	"nonzero": func(v interface{}) (interface{}, error) {
		// ensures values aren't zero or ""

		if v != nil {
			switch val := v.(type) {
			case string:
				if len(val) > 0 {
					return v, nil
				}
				return nil, errors.New("cannot be empty")
			case int, int8, int16, int32, int64,
				uint, uint8, uint16, uint32, uint64,
				float32, float64, complex64, complex128:
				if val != reflect.Zero(reflect.TypeOf(val)).Interface() {
					return v, nil
				}
			}
		}
		return nil, errors.New("cannot be zero")
	},
	"email": func(v interface{}) (interface{}, error) {
		// simple and quick email syntax checking

		var email string
		var ok bool
		if email, ok = v.(string); !ok {
			return nil, errors.New("should be a string")
		}
		if len(email) == 0 {
			// no value to check - skip
			return v, nil
		}

		atI := strings.Index(email, "@")
		ok = atI > 0 && atI < len(email)-1
		if ok {
			dotI := strings.LastIndex(email, ".")
			ok = dotI > atI && dotI < len(email)-1
		}
		if !ok {
			return nil, errors.New("is not a valid email address")
		}

		return v, nil
	},
	"lower": func(v interface{}) (interface{}, error) {
		var ok bool
		var s string
		if s, ok = v.(string); !ok {
			return nil, errors.New("is not a string")
		}
		return strings.ToLower(s), nil
	},
}

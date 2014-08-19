package is

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"
)

var Handlers = map[string]func(interface{}, bool) error{
	"required": func(v interface{}, present bool) error {
		// ensures value is present

		if !present || v == nil {
			return errors.New("is required")
		}
		return nil
	},
	"nonzero": func(v interface{}, present bool) error {
		// ensures values aren't zero or ""

		if present && v != nil {
			switch val := v.(type) {
			case string:
				if len(val) > 0 {
					return nil
				}
			case int, int8, int16, int32, int64,
				uint, uint8, uint16, uint32, uint64,
				float32, float64, complex64, complex128:
				if val != 0 {
					return nil
				}
			}
		}
		return errors.New("cannot be empty")
	},
	"email": func(v interface{}, _ bool) error {
		// simple and quick email syntax checking

		var email string
		var ok bool
		if email, ok = v.(string); !ok {
			return errors.New("should be a string")
		}
		atI := strings.Index(email, "@")
		ok = atI > 0 && atI < len(email)-1
		if ok {
			dotI := strings.LastIndex(email, ".")
			ok = dotI > atI && dotI < len(email)-1
		}
		if !ok {
			return errors.New("is not a valid email address")
		}

		return nil
	},
}

var ErrInvalidTarget = errors.New("target must be non-nil pointer to struct")

type ErrInvalidHandler struct {
	Handler string
}

func (e ErrInvalidHandler) Error() string {
	return e.Handler + " is not a valid handler"
}

type Decoder interface {
	Decode(interface{}) (map[string]error, error)
}

type msiDecoder struct {
	lock     sync.Mutex
	src      map[string]interface{}
	handlers map[string]func(interface{}, bool) error
}

func (m *msiDecoder) Decode(target interface{}) (map[string]error, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	problems := make(map[string]error)

	s := reflect.ValueOf(target)
	if s.Kind() != reflect.Ptr || s.IsNil() {
		return nil, ErrInvalidTarget
	}

	s = s.Elem()
	if s.Kind() != reflect.Struct {
		return nil, ErrInvalidTarget
	}

	t := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		switch f.Kind() {
		case reflect.Ptr:
			if f.Elem().Kind() != reflect.Struct {
				break
			}

			f = f.Elem()
			fallthrough

		case reflect.Struct:
			ss := f.Addr().Interface()
			m.Decode(ss)
		}

		if !f.CanSet() {
			continue
		}

		tag := t.Field(i).Tag.Get("is")
		if len(tag) == 0 {
			continue
		}

		segs := strings.Split(tag, ",")
		field := segs[0]
		hs := segs[1:]
		value, hasValue := m.src[field]

		// run each handler
		var ok bool = true
		for _, h := range hs {
			if handler, ok := m.handlers[h]; !ok {
				problems[field] = &ErrInvalidHandler{Handler: h}
			} else {
				// run it
				if err := handler(value, hasValue); err != nil {
					problems[field] = err
					ok = false
					break // skip to next field
				}
			}
		}
		if ok {
			rVal := reflect.ValueOf(value)
			if !rVal.IsValid() {
				continue
			}
			f.Set(rVal)
		}
	}

	if len(problems) > 0 {
		return problems, nil
	}

	return nil, nil

}

func NewMSIDecoder(src map[string]interface{}) Decoder {
	return &msiDecoder{src: src, handlers: Handlers}
}

func NewJsonDecoder(decoder *json.Decoder) Decoder {
	var data map[string]interface{}
	decoder.Decode(&data)
	return NewMSIDecoder(data)
}

package is

import (
	"errors"
	"reflect"
	"strings"
)

var ErrInvalidTarget = errors.New("target must be non-nil pointer to struct")

type Problem struct {
	Field string
	Err   error
}

func (e Problem) Error() string {
	return e.Field + " " + e.Err.Error()
}

type Problems map[string]*Problem
type HandlerFunc func(interface{}) (interface{}, error)

var DefaultValidator = NewValidator()

func Valid(value interface{}) (Problems, error) {
	return DefaultValidator.Valid(value)
}

type Validator struct {
	Handlers map[string]HandlerFunc
}

func NewValidator() *Validator {
	handlers := make(map[string]HandlerFunc)
	for k, v := range defaultHandlers {
		handlers[k] = v
	}
	return &Validator{Handlers: handlers}
}

func (v *Validator) Valid(value interface{}) (Problems, error) {

	problems := make(Problems)

	s := reflect.ValueOf(value)
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
		f2 := t.Field(i)
		tag := f2.Tag.Get("is")
		if len(tag) == 0 {
			continue
		}
		hs := strings.Split(tag, ",")
		for _, h := range hs {
			var handlerFunc HandlerFunc
			var ok bool
			if handlerFunc, ok = v.Handlers[h]; !ok {
				panic("is: no such handler " + h)
			}
			newVal, err := handlerFunc(f.Interface())
			if err != nil {
				problems[f2.Name] = &Problem{Field: f2.Name, Err: err}
				break // next field
			} else {
				newValV := reflect.ValueOf(newVal)
				f.Set(newValV)
			}
		}
	}

	return problems, nil

}

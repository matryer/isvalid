package is

import (
	"errors"
	"reflect"
	"strings"
)

// ErrInvalidTarget is returned when Valid is called on an
// invalid object. Usually fixed by passing a pointer of the
// object in, like Valid(&obj).
var ErrInvalidTarget = errors.New("target must be non-nil pointer to struct")

// Problem represents a single problem or failed validation.
type Problem struct {
	Field string
	Err   error
}

// Error gets the string describing this problem.
func (e Problem) Error() string {
	return e.Field + " " + e.Err.Error()
}

// Problems represents a map of Problems to the appropriate
// fields.
type Problems map[string]*Problem

// HandlerFunc is the a function that will be given a value,
// which it may mutate and return, or else return an error if
// something is wrong with it.
type HandlerFunc func(handler string, v interface{}) (interface{}, error)

// DefaultValidator is the default Validator type.
var DefaultValidator = NewValidator()

// Valid checks to see if the specified object is valid or
// not.
//
//     probs, err := Valid(&obj)
//     if err != nil { /* something went seriously wrong */ }
//     if len(probs) > 0 { /* some problems with the object */ }
func Valid(value interface{}) (Problems, error) {
	return DefaultValidator.Valid(value)
}

// GetValue processes the handler and gets the value or an
// error.
func GetValue(handler string, value interface{}) (interface{}, error) {
	return DefaultValidator.GetValue(handler, value)
}

// Validator represents a type capable of validating
// objects.
type Validator struct {
	Handlers map[string]HandlerFunc
}

// NewValidator makes a new Validator, configured with
// the default handlers.
func NewValidator() *Validator {
	handlers := make(map[string]HandlerFunc)
	for k, v := range defaultHandlers {
		handlers[k] = v
	}
	return &Validator{Handlers: handlers}
}

// GetValue processes the handler and gets the value or an
// error.
func (v *Validator) GetValue(handler string, value interface{}) (interface{}, error) {
	for handlerKey, hf := range v.Handlers {
		if strings.HasPrefix(handler, handlerKey) {
			return hf(handler, value)
		}
	}
	panic("is: do not understand " + handler)
}

// Valid checks to see if the specified object is valid or
// not. See Valid() for more information.
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
		fieldname := f2.Name
		// use json tag field name if available
		jsonTag := f2.Tag.Get("json")
		if len(jsonTag) > 0 {
			fieldname = strings.Split(jsonTag, ",")[0]
		}
		// process is tag
		hs := strings.Split(tag, ",")
		for _, h := range hs {
			newVal, err := v.GetValue(h, f.Interface())
			if err != nil {
				problems[fieldname] = &Problem{Field: fieldname, Err: err}
				break // next field
			} else {
				newValV := reflect.ValueOf(newVal)
				f.Set(newValV)
			}
		}
	}

	return problems, nil

}

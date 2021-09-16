package types

import (
	"reflect"
	"strings"
)

// The validation error type
type ValidationError struct {
	// Field name
	Field string
	// The kind of the value
	Kind reflect.Kind
	// The rulename
	Rule string
	// The value of the field that fails
	Value interface{}
	// The error code
	ErrorCode string
	// The error message
	Message string
}

// Message container
type MessageBag struct {
	Errors map[string][]ValidationError
}

// Creates a new message bag
func NewMessageBag() *MessageBag {
	return &MessageBag{
		Errors: make(map[string][]ValidationError),
	}
}

// Error implement
func (bag *MessageBag) Error() string {
	var errors []string
	if bag != nil {
		for k, messages := range bag.Errors {
			var currErrors []string
			for _, m := range messages {
				currErrors = append(currErrors, m.Message)
			}
			errors = append(errors, k+":["+strings.Join(currErrors, ",")+"]")
		}
	}
	return strings.Join(errors, ", ")
}

// Add a message to the bag
func (bag *MessageBag) Add(key string, err ValidationError) *MessageBag {
	_, ok := bag.Errors[key]
	if !ok {
		bag.Errors[key] = make([]ValidationError, 0)
	}
	bag.Errors[key] = append(bag.Errors[key], err)
	return bag
}

// Obtains the first error for a key (if any)
// It returns a ponter because it could be null
func (bag *MessageBag) First(key string) *ValidationError {
	m, ok := bag.Errors[key]
	if ok {
		return &m[0]
	}
	return nil
}

// Determines if the bag has a key
func (bag *MessageBag) Has(key string) bool {
	return bag.First(key) != nil
}

// Returns all the errors
func (bag *MessageBag) Messages() map[string][]string {
	if bag == nil {
		return nil
	}

	ret := make(map[string][]string)
	if bag.Errors == nil {
		return ret
	}
	for k, v := range bag.Errors {
		for _, e := range v {
			_, ok := ret[k]
			if !ok {
				ret[k] = make([]string, 0)
			}
			ret[k] = append(ret[k], e.Message)
		}
	}
	return ret
}

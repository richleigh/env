package env // github.com/richleigh/env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ErrNotAStructPtr is returned if you pass something that is not a pointer to a
// Struct to Parse
var ErrNotAStructPtr = errors.New("Expected a pointer to a Struct")

// ErrUnsuportedType if the struct field type is not supported by env
var ErrUnsuportedType = errors.New("Type is not supported")

// Parse parses a struct containing `env` tags and loads its values from
// environment variables.
func Parse(val interface{}) error {
	ptrRef := reflect.ValueOf(val)
	if ptrRef.Kind() != reflect.Ptr {
		return ErrNotAStructPtr
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return ErrNotAStructPtr
	}
	return doParse(ref, val)
}

func doParse(ref reflect.Value, val interface{}) error {
	refType := ref.Type()
	for i := 0; i < refType.NumField(); i++ {
		value, err := get(refType.Field(i))
		if err != nil {
			return err
		}
		if value == "" {
			continue
		}
		if err := set(ref.Field(i), value); err != nil {
			return err
		}
	}
	return nil
}

func get(field reflect.StructField) (string, error) {
	name := field.Tag.Get("env")
	// If there's no tag, then there's no expected envirionment variable
	if name == "" {
		return "", nil
	}
	
	// Check to see if we have the "optional" modifier
	bits := strings.Split(name, ",")
	if len(bits) >= 3 {
		return "", errors.New(fmt.Sprintf("Couldn't parse struct tag '%s'; too many ','s (expected at most 1, got %d)", name, len(bits)))
	}
	name = bits[0]
	optional := false
	if len(bits) == 2 {
		if bits[1] != "optional" {
			return "", errors.New(fmt.Sprintf("Couldn't parse struct tag '%s'; expected 'optional' after ',', got '%s'", name, bits[1]))
		}
		optional = true
	}
	
	// Now look in the environment
	value := os.Getenv(name)
	if optional || value != "" {
		return value, nil
	}
	return "", errors.New(fmt.Sprintf("Missing config environment variable '%s'", name))
}

func set(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		bvalue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(bvalue)
	case reflect.Int:
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	default:
		return ErrUnsuportedType
	}
	return nil
}

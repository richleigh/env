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
	var finalErr error
	for i := 0; i < refType.NumField(); i++ {
		value, err := get(refType.Field(i))
		if err != nil && finalErr == nil {
			finalErr = err
		}
		if value == "" {
			continue
		}
		if err := set(ref.Field(i), value); err != nil && finalErr == nil {
			finalErr = err
		}
	}
	return finalErr
}

func get(field reflect.StructField) (string, error) {
	name := field.Tag.Get("env")
	// If there's no tag, then there's no expected envirionment variable
	if name == "" {
		return "", nil
	}

	// Check to see if we have the "optional" or "sensitive" modifiers
	bits := strings.Split(name, ",")
	name = bits[0]
	optional := false
	sensitive := false
	if len(bits) > 1 {
		for _, bit := range bits[1:] {
			if bit == "optional" {
				optional = true
				continue
			}
			if bit == "sensitive" {
				sensitive = true
				continue
			}
			return "", fmt.Errorf("Couldn't parse struct tag '%s'; expected 'optional' or 'sensitive' after ',', got '%s'", name, bit)
		}
	}

	// Now look in the environment
	value := os.Getenv(name)
	if sensitive {
		if err := os.Unsetenv(name); err != nil {
			return "", err
		}
	}
	if optional || value != "" {
		return value, nil
	}
	return "", fmt.Errorf("Missing config environment variable '%s'", name)
}

func set(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		value = strings.Replace(strings.Replace(value, "\\n", "\n", -1), "\"", "", -1)
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

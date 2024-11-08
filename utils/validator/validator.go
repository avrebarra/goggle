package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	global *Validator
)

func Global() *Validator {
	if global == nil {
		global = &Validator{base: validator.New()}
	}
	return global
}

func SetGlobal(in *Validator) {
	global = in
}

func Validate(in interface{}) (err error) {
	return Global().Validate(in)
}

// ***

type Validator struct {
	base *validator.Validate
}

func New(base *validator.Validate) Validator {
	return Validator{base: base}
}

func (v Validator) Validate(in interface{}) (err error) {
	// validate
	if err = v.base.Struct(in); err == nil {
		return nil
	}

	// build validation error
	ve := ValidationError{
		root:         err.(validator.ValidationErrors),
		targtype:     reflect.TypeOf(in).Elem(),
		ErrorEntries: []ValidationErrorEntry{},
	}

	for _, ferr := range err.(validator.ValidationErrors) {
		entry := ValidationErrorEntry{Orig: ferr}
		ve.ErrorEntries = append(ve.ErrorEntries, entry)
	}

	return ve
}

// ***

type ValidationErrorEntry struct {
	Orig validator.FieldError
}

type ValidationError struct {
	root         validator.ValidationErrors
	targtype     reflect.Type
	ErrorEntries []ValidationErrorEntry
}

func (e ValidationError) GetRootError() (out validator.ValidationErrors) {
	return e.root
}

func (e ValidationError) Error() (out string) {
	errstrs := []string{}
	for _, errfield := range e.root {
		fieldName := errfield.Field()
		if f, ok := e.targtype.FieldByName(fieldName); ok {
			if val := f.Tag.Get("alias"); val != "" {
				fieldName = val
			}
		}

		msg := ""
		msg += fieldName
		msg += " must be " + errfield.Tag()

		// print condition parameters, e.g. oneof=red blue -> { red blue }
		if errfield.Param() != "" {
			msg += "{" + errfield.Param() + "}"
		}

		// print actual value
		if errfield.Value() != nil && errfield.Value() != "" {
			msg += fmt.Sprintf(", actual is %v", errfield.Value())
		}

		errstrs = append(errstrs, msg)
	}

	out += strings.Join(errstrs, "; ")

	return
}

func (e ValidationError) Unwrap() error {
	return e.root
}

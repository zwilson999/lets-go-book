package validator

import (
	"strings"
	"unicode/utf8"
)

// validator type which contains map of validation errors for form fields
type Validator struct {
	FieldErrors map[string]string
}

// check if there are any errors and return bool
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// method to add an error message to the FieldErrors map (as long as no entry for that key already exists)
func (v *Validator) AddFieldError(key, message string) {

	// initialize map if it doesn't already exist
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	// add the message for the key
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// adds an error message to our FieldErrors map only if validation check is not 'ok'
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// method to check if a value is not an empty string
func NotBlank(val string) bool {
	return strings.TrimSpace(val) != ""
}

// check if value has characters <= our set threshold of 'n'
func MaxChars(val string, n int) bool {
	return utf8.RuneCountInString(val) <= n
}

// returns true if val is in the list of permitted integers
func PermittedInt(val int, permittedVals ...int) bool {
	for i := range permittedVals {
		if val == permittedVals[i] {
			return true
		}
	}
	return false
}

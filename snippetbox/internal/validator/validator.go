package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// validator type which contains map of validation errors for form fields
type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
}

// check if there are any errors and return bool
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
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

// helper method for adding error messages to the new NonFieldErrors slice
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
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
func PermittedValue[T comparable](val T, permittedVals ...T) bool {
	for i := range permittedVals {
		if val == permittedVals[i] {
			return true
		}
	}
	return false
}

// use the regexp.MustCompile() function to parse a regular expression pattern for sanity checking the format of an email address
// this returns a pointer to a 'compiled' regexp.Regexp type, or panics in the event of an error. Parsing this pattern
// once at startup and storing the compiled *regexp.Regexp in a variable is more performant than re-parsing the pattern each time we need it.
// var EmailRx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$\\/")
var EmailRx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// method returns true if a val contains at least n characters
func MinChars(val string, n int) bool {
	return utf8.RuneCountInString(val) >= n
}

// method returns true if a val matches a provided compiled regex
func Matches(val string, rx *regexp.Regexp) bool {
	return rx.MatchString(val)
}

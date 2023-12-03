package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

// Define a New function to initialize a custom Form struct. Notice that
// this takes the form data as the parameter?

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Implement a Required method to check that specific fields in the form
// data are present and not blank. If any fields fail this check, add the
// appropriate message to the form errors.

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		val := f.Values.Get(field)
		if strings.TrimSpace(val) == "" {
			f.Errors.Add(field, "This field cannot  be blank")
		}
	}
}

// Implement a MaxLength method to check that a specific field in the form
// contains a maximum number of characters. If the check fails then add the
// appropriate message to the form errors.

func (f *Form) MaxLength(field string, d int) {
	val := f.Values.Get(field)

	if utf8.RuneCountInString(val) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d)", d))
	}
}

// Implement a PermittedValues method to check that a specific field in the form
// matches one of a set of specific permitted values. If the check fails
// then add the appropriate message to the form errors.
func (f *Form) PermittedValues(field string, opts ...string) {
	val := f.Values.Get(field)

	for _, opt := range opts {
		if val == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// Implement a Valid method which returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

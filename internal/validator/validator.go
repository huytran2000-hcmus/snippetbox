package validator

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	fieldName    string
	fieldValue   string
	NonFieldErrs []string
	FieldErrs    map[string]string
}

func (v *Validator) IsValid() bool {
	return len(v.FieldErrs) == 0 && len(v.NonFieldErrs) == 0
}

func (v *Validator) CheckField(name string, val string) *Validator {
	v.fieldName = name
	v.fieldValue = val
	return v
}

func (v *Validator) AddFieldError(fieldname string, message string) {
	if v.FieldErrs == nil {
		v.FieldErrs = map[string]string{}
	}

	_, ok := v.FieldErrs[fieldname]
	if !ok {
		v.FieldErrs[fieldname] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrs = append(v.NonFieldErrs, message)
}

func (v *Validator) ToInt(message string) int {
	i, err := strconv.Atoi(v.fieldValue)
	if err != nil {
		v.addFieldError(message)
	}
	return i
}

func (v *Validator) NotBlank(message string) *Validator {
	val := strings.TrimSpace(v.fieldValue)
	if val == "" {
		v.addFieldError(message)
	}

	return v
}

func (v *Validator) LE(message string, n int) *Validator {
	if utf8.RuneCountInString(v.fieldValue) > n {
		v.addFieldError(message)
	}

	return v
}

func (v *Validator) In(message string, permittedArrs ...string) *Validator {
	for _, permitted := range permittedArrs {
		if v.fieldValue == permitted {
			return v
		}
	}

	v.addFieldError(message)
	return v
}

func (v *Validator) GE(message string, n int) *Validator {
	if utf8.RuneCountInString(v.fieldValue) < n {
		v.addFieldError(message)
	}

	return v
}

func (v *Validator) Matches(message string, rx *regexp.Regexp) *Validator {
	ok := rx.MatchString(v.fieldValue)
	if !ok {
		v.addFieldError(message)
	}

	return v
}

func (v *Validator) IsEmail(message string) *Validator {
	return v.Matches(message, EmailRX)
}

func (v *Validator) Equal(message string, val string) *Validator {
	if v.fieldValue != val {
		v.addFieldError(message)
	}

	return v
}

func (v *Validator) addFieldError(message string) {
	v.AddFieldError(v.fieldName, message)
}

func (v *Validator) Value() string {
	return v.fieldValue
}

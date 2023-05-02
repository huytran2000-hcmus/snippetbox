package validator

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldName  string
	FieldValue string
	FieldErrs  map[string]string
}

func (v *Validator) IsValid() bool {
	return len(v.FieldErrs) == 0
}

func (v *Validator) Check(name string, val string) *Validator {
	v.FieldName = name
	v.FieldValue = val
	return v
}

func (v *Validator) ToInt() (int, error) {
	i, err := strconv.Atoi(v.FieldValue)
	return i, err
}

func (v *Validator) NotBlank(message string) *Validator {
	val := strings.TrimSpace(v.FieldValue)
	if val == "" {
		v.addFieldError(message)
	}

	return v
}

func (v *Validator) LE(message string, n int) *Validator {
	if utf8.RuneCountInString(v.FieldValue) > n {
		v.addFieldError(message)
	}

	return v
}

func (v *Validator) In(message string, permittedArrs ...string) *Validator {
	for _, permitted := range permittedArrs {
		if v.FieldValue == permitted {
			return v
		}
	}

	v.addFieldError(message)
	return v
}

func (v *Validator) addFieldError(message string) {
	if v.FieldErrs == nil {
		v.FieldErrs = map[string]string{}
	}

	name := v.FieldName
	_, ok := v.FieldErrs[name]
	if !ok {
		v.FieldErrs[name] = message
	}
}

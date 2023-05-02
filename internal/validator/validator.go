package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	fieldName string
	fieldVal  string
	FieldErrs map[string]string
}

func (v *Validator) IsValid() bool {
	return len(v.FieldErrs) == 0
}

func (v *Validator) Check(name string, val string) *Validator {
	v.fieldName = name
	v.fieldVal = val
	return v
}

func (v *Validator) NotBlank(message string) *Validator {
	val := strings.TrimSpace(v.fieldVal)
	if val == "" {
		v.addFieldError(message)
	}

	return v
}

func (v *Validator) MaxCharacters(message string, n int) *Validator {
	if utf8.RuneCountInString(v.fieldVal) > n {
		v.addFieldError(message)
	}

	return v
}

func (v *Validator) InPermittedArr(message string, permittedArrs ...string) *Validator {
	for _, permitted := range permittedArrs {
		if v.fieldVal == permitted {
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

	name := v.fieldName
	_, ok := v.FieldErrs[name]
	if !ok {
		v.FieldErrs[name] = message
	}
}

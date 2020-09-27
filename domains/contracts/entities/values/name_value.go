package values

import (
	"fmt"
	"github.com/mixmaru/my_contracts/domains/contracts/entities/values/validators"
	"github.com/pkg/errors"
	"strings"
)

const NameMaxLength = 50

// Name値オブジェクト
type NameValue struct {
	value string
}

func NewNameValue(value string) (NameValue, error) {
	validateErrors := NameValidate(value)
	if len(validateErrors) > 0 {
		var msgs []string
		for _, validateError := range validateErrors {
			msgs = append(msgs, validators.ValidErrorText(validateError))
		}
		return NameValue{}, errors.New(fmt.Sprintf("Nameバリデーションエラー。%v", strings.Join(msgs, ", ")))
	}
	return NameValue{
		value: value,
	}, nil
}

func (v *NameValue) Value() string {
	return v.value
}

func NameValidate(name string) (validErrors []int) {
	if validators.IsEmptyString(name) {
		validErrors = append(validErrors, validators.EmptyStringValidError)
	}
	if validators.IsOverLengthString(name, NameMaxLength) {
		validErrors = append(validErrors, validators.OverLengthStringValidError)
	}

	return validErrors
}

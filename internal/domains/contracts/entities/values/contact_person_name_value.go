package values

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values/validators"
	"github.com/pkg/errors"
	"strings"
)

// Name値オブジェクト
type ContactPersonNameValue struct {
	value string
}

const MaxContactPersonNameNum = 50

func NewContactPersonNameValue(value string) (ContactPersonNameValue, error) {
	validateErrors, err := ContactPersonNameValue{}.Validate(value)
	if err != nil {
		return ContactPersonNameValue{}, err
	}

	if len(validateErrors) > 0 {
		var msgs []string
		for _, validateError := range validateErrors {
			msgs = append(msgs, validators.ValidErrorTest(validateError))
		}
		return ContactPersonNameValue{}, errors.New(fmt.Sprintf("Nameバリデーションエラー。%v", strings.Join(msgs, ", ")))
	}
	return ContactPersonNameValue{
		value: value,
	}, nil
}

func (v *ContactPersonNameValue) Value() string {
	return v.value
}

func (v ContactPersonNameValue) Validate(value interface{}) (validErrors []int, err error) {
	name, ok := value.(string)
	if !ok {
		return nil, errors.Errorf("valueをstring型にできませんでした。%t", value)
	}
	return v.contactPersonNameValidate(name), nil
}

func (v *ContactPersonNameValue) contactPersonNameValidate(name string) (validErrors []int) {
	if validators.IsEmptyString(name) {
		validErrors = append(validErrors, validators.EmptyStringValidError)
	}
	if validators.IsOverLengthString(name, MaxContactPersonNameNum) {
		validErrors = append(validErrors, validators.OverLengthStringValidError)
	}

	return validErrors
}

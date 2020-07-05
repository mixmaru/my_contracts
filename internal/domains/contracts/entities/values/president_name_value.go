package values

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values/validators"
	"github.com/pkg/errors"
	"strings"
)

// 社長名値オブジェクト
type PresidentNameValue struct {
	value string
}

const MaxPresidentNameNum = 50

func NewPresidentNameValue(value string) (PresidentNameValue, error) {
	validateErrors, err := PresidentNameValue{}.Validate(value)
	if err != nil {
		return PresidentNameValue{}, err
	}
	if len(validateErrors) > 0 {
		var msgs []string
		for _, validateError := range validateErrors {
			msgs = append(msgs, validators.ValidErrorTest(validateError))
		}
		return PresidentNameValue{}, errors.New(fmt.Sprintf("PresidentNameバリデーションエラー。%v", strings.Join(msgs, ", ")))
	}
	return PresidentNameValue{
		value: value,
	}, nil
}

func (v *PresidentNameValue) Value() string {
	return v.value
}

func (v PresidentNameValue) Validate(value interface{}) (validErrors []int, err error) {
	name, ok := value.(string)
	if !ok {
		return nil, errors.Errorf("valueをstring型にできませんでした。value: %t", value)
	}
	return v.presidentNameValidate(name), nil
}

func (v *PresidentNameValue) presidentNameValidate(name string) (validErrors []int) {
	if validators.IsEmptyString(name) {
		validErrors = append(validErrors, validators.EmptyStringValidError)
	}
	if validators.IsOverLengthString(name, MaxPresidentNameNum) {
		validErrors = append(validErrors, validators.OverLengthStringValidError)
	}

	return validErrors
}

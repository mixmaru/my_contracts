package values

import (
	"fmt"
	"github.com/mixmaru/my_contracts/domains/contracts/entities/values/validators"
	"github.com/pkg/errors"
	"strings"
)

// 社長名値オブジェクト
type CorporationNameValue struct {
	value string
}

const MaxCorporationNameNum = 50

func NewCorporationNameValue(value string) (CorporationNameValue, error) {
	validateErrors, err := CorporationNameValue{}.Validate(value)
	if err != nil {
		return CorporationNameValue{}, err
	}
	if len(validateErrors) > 0 {
		var msgs []string
		for _, validateError := range validateErrors {
			msgs = append(msgs, validators.ValidErrorText(validateError))
		}
		return CorporationNameValue{}, errors.New(fmt.Sprintf("CorporationNameバリデーションエラー。%v", strings.Join(msgs, ", ")))
	}
	return CorporationNameValue{
		value: value,
	}, nil
}

func (v *CorporationNameValue) Value() string {
	return v.value
}

func (v CorporationNameValue) Validate(value interface{}) (validErrors []int, err error) {
	name, ok := value.(string)
	if !ok {
		return nil, errors.Errorf("valueをstring型にできませんでした。value: %t", value)
	}
	return v.corporationNameValidate(name), nil
}

func (v *CorporationNameValue) corporationNameValidate(name string) (validErrors []int) {
	if validators.IsEmptyString(name) {
		validErrors = append(validErrors, validators.EmptyStringValidError)
	}
	if validators.IsOverLengthString(name, MaxPresidentNameNum) {
		validErrors = append(validErrors, validators.OverLengthStringValidError)
	}

	return validErrors
}

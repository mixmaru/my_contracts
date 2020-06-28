package values

import (
	plain_err "errors"
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values/validators"
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
		for _, msg := range validateErrors {
			msgs = append(msgs, msg.Error())
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

func NameValidate(name string) []error {
	var validErrors []error
	if validators.IsEmptyString(name) {
		validErrors = append(validErrors, validators.NewEmptyStringValidError(plain_err.New("空です")))
	}
	if validators.IsOverLengthString(name, NameMaxLength) {
		validErrors = append(validErrors, validators.NewOverLengthStringValidError(plain_err.New(fmt.Sprintf("%v文字より多いです", NameMaxLength))))
	}

	return validErrors
}

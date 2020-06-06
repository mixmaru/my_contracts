package values

import (
	plain_err "errors"
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user/values/validators"
	"github.com/pkg/errors"
	"strings"
)

// Name値オブジェクト
type ContactPersonNameValue struct {
	value string
}

const MaxContactPersonNameNum = 50

func NewContactPersonNameValue(value string) (ContactPersonNameValue, error) {
	validateErrors := ContactPersonNameValidate(value)
	if len(validateErrors) > 0 {
		var msgs []string
		for _, msg := range validateErrors {
			msgs = append(msgs, msg.Error())
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

func ContactPersonNameValidate(name string) []error {
	var validErrors []error
	if validators.IsEmptyString(name) {
		validErrors = append(validErrors, validators.NewEmptyValidError(plain_err.New("空です")))
	}
	if validators.IsOverLengthString(name, MaxContactPersonNameNum) {
		validErrors = append(validErrors, validators.NewOverLengthValidError(plain_err.New(fmt.Sprintf("%v文字より多いです", MaxContactPersonNameNum))))
	}

	return validErrors
}

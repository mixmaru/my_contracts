package values

import (
	plain_err "errors"
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
	validateErrors := PresidentNameValidate(value)
	if len(validateErrors) > 0 {
		var msgs []string
		for _, msg := range validateErrors {
			msgs = append(msgs, msg.Error())
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

func PresidentNameValidate(name string) []error {
	var validErrors []error
	if validators.IsEmptyString(name) {
		validErrors = append(validErrors, validators.NewEmptyValidError(plain_err.New("空です")))
	}
	if validators.IsOverLengthString(name, MaxPresidentNameNum) {
		validErrors = append(validErrors, validators.NewOverLengthValidError(plain_err.New(fmt.Sprintf("%v文字より多いです", MaxPresidentNameNum))))
	}

	return validErrors
}

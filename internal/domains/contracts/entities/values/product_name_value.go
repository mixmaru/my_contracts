package values

import (
	plain_err "errors"
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values/validators"
	"github.com/pkg/errors"
	"strings"
)

const ProductNameMaxLength = 50

// Name値オブジェクト
type ProductNameValue struct {
	value string
}

func NewProductNameValue(value string) (ProductNameValue, error) {
	validateErrors := ProductNameValidate(value)
	if len(validateErrors) > 0 {
		var msgs []string
		for _, msg := range validateErrors {
			msgs = append(msgs, msg.Error())
		}
		return ProductNameValue{}, errors.New(fmt.Sprintf("ProductNameバリデーションエラー。%v", strings.Join(msgs, ", ")))
	}
	return ProductNameValue{
		value: value,
	}, nil
}

func (v *ProductNameValue) Value() string {
	return v.value
}

func ProductNameValidate(name string) []error {
	var validErrors []error
	if validators.IsEmptyString(name) {
		validErrors = append(validErrors, validators.NewEmptyStringValidError(plain_err.New("空です")))
	}
	if validators.IsOverLengthString(name, ProductNameMaxLength) {
		validErrors = append(validErrors, validators.NewOverLengthStringValidError(plain_err.New(fmt.Sprintf("%v文字より多いです", NameMaxLength))))
	}

	return validErrors
}

package values

import (
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
	validateErrors, err := ProductNameValue{}.Validate(value)
	if err != nil {
		return ProductNameValue{}, err
	}

	if len(validateErrors) > 0 {
		var msgs []string
		for _, validateError := range validateErrors {
			msgs = append(msgs, validators.ValidErrorTest(validateError))
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

func (v ProductNameValue) Validate(value interface{}) (validErrors []int, err error) {
	name, ok := value.(string)
	if !ok {
		return nil, errors.New(fmt.Sprintf("valueがstring型にできませんでした。value: %t", value))
	}

	if validators.IsEmptyString(name) {
		validErrors = append(validErrors, validators.EmptyStringValidError)
	}
	if validators.IsOverLengthString(name, ProductNameMaxLength) {
		validErrors = append(validErrors, validators.OverLengthStringValidError)
	}

	return validErrors, nil
}

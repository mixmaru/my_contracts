package values

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/pkg/errors"
	"strings"
)

// product price値オブジェクト
type ProductPriceValue struct {
	value decimal.Decimal
}

func NewProductPriceValue(value string) (ProductPriceValue, error) {
	validateErrors := ProductPriceValidate(value)
	if len(validateErrors) > 0 {
		var msgs []string
		for _, msg := range validateErrors {
			msgs = append(msgs, msg.Error())
		}
		return ProductPriceValue{}, errors.New(fmt.Sprintf("ProductPriceバリデーションエラー。%v", strings.Join(msgs, ", ")))
	}
	priceDecimal, err := decimal.NewFromString(value)
	if err != nil {
		return ProductPriceValue{}, err
	}

	return ProductPriceValue{
		value: priceDecimal,
	}, nil
}

func (v *ProductPriceValue) Value() decimal.Decimal {
	return v.value
}

func ProductPriceValidate(price string) []error {
	var validErrors []error
	//
	//if validators.IsEmptyString(name) {
	//	validErrors = append(validErrors, validators.NewEmptyValidError(plain_err.New("空です")))
	//}
	//if validators.IsOverLengthString(name, NameMaxLength) {
	//	validErrors = append(validErrors, validators.NewOverLengthValidError(plain_err.New(fmt.Sprintf("%v文字より多いです", NameMaxLength))))
	//}

	return validErrors
}

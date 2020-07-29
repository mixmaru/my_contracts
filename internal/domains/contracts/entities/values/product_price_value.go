package values

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values/validators"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/pkg/errors"
	"strings"
)

// product price値オブジェクト
type ProductPriceValue struct {
	value decimal.Decimal
}

func NewProductPriceValue(value string) (ProductPriceValue, error) {
	validateErrors, err := ProductPriceValue{}.Validate(value)
	if err != nil {
		return ProductPriceValue{}, err
	}

	if len(validateErrors) > 0 {
		var msgs []string
		for _, validateError := range validateErrors {
			msgs = append(msgs, validators.ValidErrorText(validateError))
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

func (v ProductPriceValue) Validate(value interface{}) (validErrors []int, err error) {
	price, ok := value.(string)
	if !ok {
		return nil, errors.Errorf("valueをstring型にできませんでした。value: %t", value)
	}
	return v.productPriceValidate(price)
}

func (v *ProductPriceValue) productPriceValidate(price string) (validErrors []int, err error) {
	// 空文字チェック
	if validators.IsEmptyString(price) {
		validErrors = append(validErrors, validators.EmptyStringValidError)
		return validErrors, nil
	}

	// 数値チェック
	if !validators.IsNumericString(price) {
		validErrors = append(validErrors, validators.NumericStringValidError)
		return validErrors, nil
	}

	// 正値チェック
	priceDecimal, err := decimal.NewFromString(price)
	if err != nil {
		return nil, err
	}
	if priceDecimal.IsNegative() {
		validErrors = append(validErrors, validators.NegativeValidError)
	}

	return validErrors, nil
}

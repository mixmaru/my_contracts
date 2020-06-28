package values

import (
	plain_err "errors"
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
	validateErrors, err := ProductPriceValidate(value)
	if err != nil {
		return ProductPriceValue{}, err
	}

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

func ProductPriceValidate(price string) (validErrors []error, err error) {
	// 空文字チェック
	if validators.IsEmptyString(price) {
		validErrors = append(validErrors, validators.NewEmptyStringValidError(plain_err.New("空です")))
		return validErrors, nil
	}

	// 数値チェック
	if !validators.IsNumericString(price) {
		validErrors = append(validErrors, validators.NewNumericStringValidError(plain_err.New("数値ではありません")))
		return validErrors, nil
	}

	// 正値チェック
	priceDecimal, err := decimal.NewFromString(price)
	if err != nil {
		return nil, err
	}
	if priceDecimal.IsNegative() {
		validErrors = append(validErrors, validators.NewNegativeValidError(plain_err.New("マイナスの数値です")))
	}

	return validErrors, nil
}

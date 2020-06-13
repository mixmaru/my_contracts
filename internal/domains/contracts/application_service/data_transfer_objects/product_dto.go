package data_transfer_objects

import (
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
)

type ProductDto struct {
	Name  string
	Price decimal.Decimal
	BaseDto
}

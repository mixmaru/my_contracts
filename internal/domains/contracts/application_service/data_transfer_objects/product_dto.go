package data_transfer_objects

import (
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"time"
)

type ProductDto struct {
	Id        int
	Name      string
	Price     decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
}

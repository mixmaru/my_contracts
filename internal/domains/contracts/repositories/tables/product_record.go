package tables

import (
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
)

type ProductRecord struct {
	Id    int             `db:"id"`
	Name  string          `db:"name"`
	Price decimal.Decimal `db:"price"`
	CreateAtUpdateAt
}

package tables

import (
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"time"
)

type CreateAtUpdateAt struct {
	CreateAt time.Time `db:"created_at"`
	UpdateAt time.Time `db:"updated_at"`
}

type ProductRecord struct {
	Id    int             `db:"id"`
	Name  string          `db:"name"`
	Price decimal.Decimal `db:"price"`
	CreateAtUpdateAt
}

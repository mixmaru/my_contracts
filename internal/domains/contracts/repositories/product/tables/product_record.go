package tables

import (
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"time"
)

type ProductRecord struct {
	Id       int             `db:"id"`
	Name     string          `db:"name"`
	Price    decimal.Decimal `db:"price"`
	CreateAt time.Time       `db:"created_at"`
	UpdateAt time.Time       `db:"updated_at"`
}

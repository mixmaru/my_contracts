package data_mappers

import "github.com/mixmaru/my_contracts/internal/lib/decimal"

type ProductPriceMonthlyMapper struct {
	ProductId int             `db:"product_id"`
	Price     decimal.Decimal `db:"price"`
	CreatedAtUpdatedAtMapper
}

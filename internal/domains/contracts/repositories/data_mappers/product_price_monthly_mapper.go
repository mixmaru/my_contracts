package data_mappers

type ProductPriceMonthlyMapper struct {
	ProductId int `db:"product_id"`
	CreatedAtUpdatedAtMapper
}

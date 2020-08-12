package data_mappers

type ProductMapper struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	CreatedAtUpdatedAtMapper
}

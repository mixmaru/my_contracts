package tables

import "time"

type ProductRecord struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	//Price float64
	CreateAt time.Time `db:"created_at"`
	UpdateAt time.Time `db:"updated_at"`
}

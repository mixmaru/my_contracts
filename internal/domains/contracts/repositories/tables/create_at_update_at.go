package tables

import "time"

type CreateAtUpdateAt struct {
	CreateAt time.Time `db:"created_at"`
	UpdateAt time.Time `db:"updated_at"`
}

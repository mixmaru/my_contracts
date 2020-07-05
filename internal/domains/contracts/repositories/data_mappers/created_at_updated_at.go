package data_mappers

import (
	"gopkg.in/gorp.v2"
	"time"
)

type CreatedAtUpdatedAt struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// insert時に時刻をセットするhook
func (c *CreatedAtUpdatedAt) PreInsert(s gorp.SqlExecutor) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = c.CreatedAt
	return nil
}

// updateに時刻をセットするhook
func (c *CreatedAtUpdatedAt) PreUpdate(s gorp.SqlExecutor) error {
	c.UpdatedAt = time.Now()
	return nil
}

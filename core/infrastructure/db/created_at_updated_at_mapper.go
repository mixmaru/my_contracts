package db

import (
	"gopkg.in/gorp.v2"
	"time"
)

type CreatedAtUpdatedAtMapper struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// insert時に時刻をセットするhook
func (c *CreatedAtUpdatedAtMapper) PreInsert(s gorp.SqlExecutor) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = c.CreatedAt
	return nil
}

// updateに時刻をセットするhook
func (c *CreatedAtUpdatedAtMapper) PreUpdate(s gorp.SqlExecutor) error {
	c.UpdatedAt = time.Now()
	return nil
}

package entities

import "time"

type IBaseEntity interface {
	Id() int
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type BaseEntity struct {
	id        int
	createdAt time.Time
	updatedAt time.Time
}

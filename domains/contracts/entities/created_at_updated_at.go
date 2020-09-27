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

func (b *BaseEntity) Id() int {
	return b.id
}

func (b *BaseEntity) CreatedAt() time.Time {
	return b.createdAt
}

func (b *BaseEntity) UpdatedAt() time.Time {
	return b.updatedAt
}

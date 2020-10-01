package entities

import (
	"time"
)

type RightToUseEntity struct {
	BaseEntity
	validFrom time.Time
	validTo   time.Time
}

func NewRightToUseEntity(validFrom, validTo time.Time) *RightToUseEntity {
	return &RightToUseEntity{
		validFrom: validFrom,
		validTo:   validTo,
	}
}

func (r *RightToUseEntity) ValidFrom() time.Time {
	return r.validFrom
}

func (r *RightToUseEntity) ValidTo() time.Time {
	return r.validTo
}

func NewRightToUseEntityWithData(id int, validFrom, validTo, createdAt, updatedAt time.Time) *RightToUseEntity {
	entity := &RightToUseEntity{}
	entity.LoadData(id, validFrom, validTo, createdAt, updatedAt)
	return entity
}

func (r *RightToUseEntity) LoadData(id int, validFrom, validTo, createdAt, updatedAt time.Time) {
	r.id = id
	r.validFrom = validFrom
	r.validTo = validTo
	r.createdAt = createdAt
	r.updatedAt = updatedAt
}

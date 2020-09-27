package entities

import (
	"time"
)

type RightToUseEntity struct {
	BaseEntity
	contractId int
	validFrom  time.Time
	validTo    time.Time
}

func NewRightToUseEntity(contractId int, validFrom, validTo time.Time) *RightToUseEntity {
	return &RightToUseEntity{
		contractId: contractId,
		validFrom:  validFrom,
		validTo:    validTo,
	}
}

func (r *RightToUseEntity) ContractId() int {
	return r.contractId
}

func (r *RightToUseEntity) ValidFrom() time.Time {
	return r.validFrom
}

func (r *RightToUseEntity) ValidTo() time.Time {
	return r.validTo
}

func NewRightToUseEntityWithData(id, contractId int, validFrom, validTo, createdAt, updatedAt time.Time) *RightToUseEntity {
	entity := &RightToUseEntity{}
	entity.LoadData(id, contractId, validFrom, validTo, createdAt, updatedAt)
	return entity
}

func (r *RightToUseEntity) LoadData(id, contractId int, validFrom, validTo, createdAt, updatedAt time.Time) {
	r.id = id
	r.contractId = contractId
	r.validFrom = validFrom
	r.validTo = validTo
	r.createdAt = createdAt
	r.updatedAt = updatedAt
}

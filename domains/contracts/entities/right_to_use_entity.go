package entities

import (
	"time"
)

type RightToUseEntity struct {
	BaseEntity
	validFrom    time.Time
	validTo      time.Time
	billDetailId int // 請求詳細への関連。0だったら未請求
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

func (r *RightToUseEntity) BillDetailId() int {
	return r.billDetailId
}

func NewRightToUseEntityWithData(id int, validFrom, validTo time.Time, billDetailId int, createdAt, updatedAt time.Time) *RightToUseEntity {
	entity := &RightToUseEntity{}
	entity.LoadData(id, validFrom, validTo, billDetailId, createdAt, updatedAt)
	return entity
}

func (r *RightToUseEntity) LoadData(id int, validFrom, validTo time.Time, billDetailId int, createdAt, updatedAt time.Time) {
	r.id = id
	r.validFrom = validFrom
	r.validTo = validTo
	r.billDetailId = billDetailId
	r.createdAt = createdAt
	r.updatedAt = updatedAt
}

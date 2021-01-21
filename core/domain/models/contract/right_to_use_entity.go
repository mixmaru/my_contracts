package contract

import "time"

type RightToUseEntity struct {
	id           int
	validFrom    time.Time
	validTo      time.Time
	billDetailId int // 請求詳細への関連。0だったら未請求
	createdAt    time.Time
	updatedAt    time.Time
}

func NewRightToUseEntity(validFrom, validTo time.Time) *RightToUseEntity {
	return &RightToUseEntity{
		validFrom: validFrom,
		validTo:   validTo,
	}
}

func (r *RightToUseEntity) Id() int {
	return r.id
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

func (r *RightToUseEntity) CreatedAt() time.Time {
	return r.createdAt
}

func (r *RightToUseEntity) UpdatedAt() time.Time {
	return r.updatedAt
}

// 請求済かどうかを返す
func (r *RightToUseEntity) WasBilling() bool {
	if r.billDetailId == 0 {
		return false
	} else {
		return true
	}
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

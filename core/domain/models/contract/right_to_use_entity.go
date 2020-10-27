package contract

import "time"

type rightToUseEntity struct {
	id           int
	validFrom    time.Time
	validTo      time.Time
	billDetailId int // 請求詳細への関連。0だったら未請求
	createdAt    time.Time
	updatedAt    time.Time
}

func newRightToUseEntity(validFrom, validTo time.Time) *rightToUseEntity {
	return &rightToUseEntity{
		validFrom: validFrom,
		validTo:   validTo,
	}
}

func (r *rightToUseEntity) Id() int {
	return r.id
}

func (r *rightToUseEntity) ValidFrom() time.Time {
	return r.validFrom
}

func (r *rightToUseEntity) ValidTo() time.Time {
	return r.validTo
}

func (r *rightToUseEntity) BillDetailId() int {
	return r.billDetailId
}

func (r *rightToUseEntity) CreatedAt() time.Time {
	return r.createdAt
}

func (r *rightToUseEntity) UpdatedAt() time.Time {
	return r.updatedAt
}

// 請求済かどうかを返す
func (r *rightToUseEntity) WasBilling() bool {
	if r.billDetailId == 0 {
		return false
	} else {
		return true
	}
}

func newRightToUseEntityWithData(id int, validFrom, validTo time.Time, billDetailId int, createdAt, updatedAt time.Time) *rightToUseEntity {
	entity := &rightToUseEntity{}
	entity.LoadData(id, validFrom, validTo, billDetailId, createdAt, updatedAt)
	return entity
}

func (r *rightToUseEntity) LoadData(id int, validFrom, validTo time.Time, billDetailId int, createdAt, updatedAt time.Time) {
	r.id = id
	r.validFrom = validFrom
	r.validTo = validTo
	r.billDetailId = billDetailId
	r.createdAt = createdAt
	r.updatedAt = updatedAt
}

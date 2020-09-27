package entities

import (
	"time"
)

type IUserEntity interface {
}

type UserEntity struct {
	BaseEntity
}

// 保持データをセットし直す
func (u *UserEntity) LoadData(id int, createdAt time.Time, updatedAt time.Time) {
	u.id = id
	u.createdAt = createdAt
	u.updatedAt = updatedAt
}

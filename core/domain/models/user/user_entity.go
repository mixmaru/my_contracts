package entities

import (
	"time"
)

type IUserEntity interface {
}

type UserEntity struct {
	id        int
	createdAt time.Time
	updatedAt time.Time
}

// 保持データをセットし直す
func (u *UserEntity) LoadData(id int, createdAt time.Time, updatedAt time.Time) {
	u.id = id
	u.createdAt = createdAt
	u.updatedAt = updatedAt
}

func (u *UserEntity) Id() int {
	return u.id
}

func (u *UserEntity) CreatedAt() time.Time {
	return u.createdAt
}

func (u *UserEntity) UpdatedAt() time.Time {
	return u.updatedAt
}

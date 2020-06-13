package entities

import (
	"time"
)

type IUserEntity interface {
}

type UserEntity struct {
	BaseEntity
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

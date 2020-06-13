package entities

import (
	"time"
)

type IUserEntity interface {
	Id() int
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type UserEntity struct {
	id        int
	createdAt time.Time
	updatedAt time.Time
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

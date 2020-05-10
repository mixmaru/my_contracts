package user

import (
	"time"
)

type IUser interface {
	Id() int
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type User struct {
	id        int
	createdAt time.Time
	updatedAt time.Time
}

func (u *User) Id() int {
	return u.id
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

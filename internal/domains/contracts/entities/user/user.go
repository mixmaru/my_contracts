package user

import (
	"time"
)

type IUser interface {
	Id() int
	SetId(id int)

	CreatedAt() time.Time
	SetCreatedAt(time time.Time)

	UpdatedAt() time.Time
	SetUpdatedAt(time time.Time)
}

type User struct {
	id        int
	createdAt time.Time
	updatedAt time.Time
}

func (u *User) Id() int {
	return u.id
}

func (u *User) SetId(id int) {
	u.id = id
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) SetCreatedAt(time time.Time) {
	u.createdAt = time
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) SetUpdatedAt(time time.Time) {
	u.updatedAt = time
}

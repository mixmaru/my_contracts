package user

import (
	"time"
)

type IUser interface {
	Id() int
}

type User struct {
	id        int
	createdAt time.Time
	updatedAt time.Time
}

func (u *User) Id() int {
	return u.id
}

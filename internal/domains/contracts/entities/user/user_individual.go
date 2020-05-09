package user

import (
	"time"
)

type UserIndividual struct {
	User
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func (u *UserIndividual) Name() string {
	return u.name
}

func (u *UserIndividual) SetName(name string) {
	u.name = name
}

package user

import (
	"time"
)

type UserIndividual struct {
	*User
	name string
}

func NewUserIndividual() *UserIndividual {
	return &UserIndividual{
		User: &User{},
		name: "",
	}
}

func NewUserIndividualWithData(id int, name string, createdAt time.Time, updatedAt time.Time) *UserIndividual {
	userIndividual := NewUserIndividual()
	userIndividual.id = id
	userIndividual.name = name
	userIndividual.createdAt = createdAt
	userIndividual.updatedAt = updatedAt

	return userIndividual
}

// 保持データをセットし直す
func (u *UserIndividual) LoadData(id int, name string, createdAt time.Time, updatedAt time.Time) {
	u.id = id
	u.name = name
	u.createdAt = createdAt
	u.updatedAt = updatedAt
}

func (u *UserIndividual) Name() string {
	return u.name
}

func (u *UserIndividual) SetName(name string) {
	u.name = name
}

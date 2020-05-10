package user

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user/structures"
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

func LoadUserIndividual(data *structures.UserIndividualView) *UserIndividual {
	userIndividual := NewUserIndividual()
	userIndividual.id = data.Id
	userIndividual.name = data.Name
	userIndividual.createdAt = data.CreatedAt
	userIndividual.updatedAt = data.UpdatedAt

	return userIndividual
}

func (u *UserIndividual) LoadUserIndividual(data *structures.UserIndividualView) {
	u.id = data.Id
	u.name = data.Name
	u.createdAt = data.CreatedAt
	u.updatedAt = data.UpdatedAt
}

func (u *UserIndividual) Name() string {
	return u.name
}

func (u *UserIndividual) SetName(name string) {
	u.name = name
}

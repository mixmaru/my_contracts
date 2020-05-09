package user

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user/structures"
	"time"
)

type UserIndividual struct {
	*User
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func LoadUserIndividual(data *structures.UserIndividualView) *UserIndividual {
	user := &User{}
	user.id = data.Id
	user.createdAt = data.CreatedAt
	user.updatedAt = data.UpdatedAt

	userIndividual := &UserIndividual{}
	userIndividual.User = user
	userIndividual.name = data.Name
	userIndividual.createdAt = data.CreatedAt
	userIndividual.updatedAt = data.UpdatedAt

	return userIndividual
}

func (u *UserIndividual) Name() string {
	return u.name
}

func (u *UserIndividual) SetName(name string) {
	u.name = name
}

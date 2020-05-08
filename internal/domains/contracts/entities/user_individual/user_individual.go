package user_individual

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user_individual/values"
)

type UserIndividual struct {
	user.User
	name      values.Name
	createdAt values.CreatedAt
	updatedAt values.UpdatedAt
}

func (u *UserIndividual) Name() values.Name {
	return u.name
}

func (u *UserIndividual) SetName(name values.Name) {
	u.name = name
}

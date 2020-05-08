package user

import "github.com/mixmaru/my_contracts/internal/domains/contracts/entities/common_values"

type UserIndividual struct {
	User
	name      common_values.Name
	createdAt common_values.CreatedAt
	updatedAt common_values.UpdatedAt
}

func (u *UserIndividual) Name() common_values.Name {
	return u.name
}

func (u *UserIndividual) SetName(name common_values.Name) {
	u.name = name
}

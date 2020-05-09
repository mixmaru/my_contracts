package user_individual

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"time"
)

type UserIndividual struct {
	user.User
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

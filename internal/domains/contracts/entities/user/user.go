package user

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user/values"
)

type IUser interface {
	Id() values.Id
}

type User struct {
	id        values.Id
	createdAt values.CreatedAt
	updatedAt values.UpdatedAt
}

func (u *User) Id() values.Id {
	return u.id
}

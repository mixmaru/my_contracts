package user

import "github.com/mixmaru/my_contracts/internal/domains/contracts/entities/common_values"

type IUser interface {
	Id() common_values.IdInt
}

type User struct {
	id        common_values.IdInt
	createdAt common_values.CreatedAt
	updatedAt common_values.UpdatedAt
}

func (u *User) Id() common_values.IdInt {
	return u.id
}

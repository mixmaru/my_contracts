package user

import "./values"

type IUser interface {
	GetId() values.IdInt
	SetId(id values.IdInt)
}

type User struct {
	id        values.IdInt
	createdAt values.CreatedAt
	updatedAt values.UpdatedAt
}

func (u *User) GetId() values.IdInt {
	return u.id
}

func (u *User) SetId(id values.IdInt) {
	u.id = id
}

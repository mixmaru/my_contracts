package user

import "./values"

type User struct {
	Id        values.IdInt
	createdAt values.CreatedAt
	updatedAt values.UpdatedAt
}

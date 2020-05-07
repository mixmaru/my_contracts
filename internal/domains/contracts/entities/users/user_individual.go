package user

import "./values"

type UserIndividual struct {
	User
	name      Name
	createdAt values.CreatedAt
	updatedAt values.UpdatedAt
}

func (u *UserIndividual) GetName() Name {
	return u.name
}

////// 値オブジェクト定義
type Name struct {
	value string
}

func NewName(name string) Name {
	return Name{
		value: name,
	}
}

func (n *Name) GetValue() string {
	return n.value
}

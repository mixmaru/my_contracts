package user

import "github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user/values"

type UserIndividual struct {
	User
	name      Name
	createdAt values.CreatedAt
	updatedAt values.UpdatedAt
}

func (u *UserIndividual) GetName() Name {
	return u.name
}

func (u *UserIndividual) SetName(name Name) {
	u.name = name
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

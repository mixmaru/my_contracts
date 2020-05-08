package user

//////// エンティティ定義 ///////////
type UserIndividual struct {
	User
	name      Name
	createdAt CreatedAt
	updatedAt UpdatedAt
}

func (u *UserIndividual) Name() Name {
	return u.name
}

func (u *UserIndividual) SetName(name Name) {
	u.name = name
}

//////// 値オブジェクト定義 ///////////
type Name struct {
	value string
}

func NewName(name string) Name {
	return Name{
		value: name,
	}
}

func (n *Name) Value() string {
	return n.value
}

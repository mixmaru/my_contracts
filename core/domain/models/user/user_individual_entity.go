package user

import (
	"time"
)

type UserIndividualEntity struct {
	*UserEntity
	name NameValue
}

func NewUserIndividualEntity(name string) (*UserIndividualEntity, error) {
	nameValue, err := NewNameValue(name)
	if err != nil {
		return nil, err
	}

	return &UserIndividualEntity{
		UserEntity: &UserEntity{},
		name:       nameValue,
	}, nil
}

func NewUserIndividualEntityWithData(id int, name string, createdAt time.Time, updatedAt time.Time) (*UserIndividualEntity, error) {
	userIndividual, err := NewUserIndividualEntity(name)
	if err != nil {
		return nil, err
	}
	userIndividual.id = id
	userIndividual.createdAt = createdAt
	userIndividual.updatedAt = updatedAt

	return userIndividual, nil
}

// 保持データをセットし直す
func (u *UserIndividualEntity) LoadData(id int, name string, createdAt time.Time, updatedAt time.Time) error {
	nameValue, err := NewNameValue(name)
	if err != nil {
		return err
	}

	if u.UserEntity == nil {
		u.UserEntity = &UserEntity{}
	}

	u.id = id
	u.name = nameValue
	u.createdAt = createdAt
	u.updatedAt = updatedAt
	return nil
}

func (u *UserIndividualEntity) Name() string {
	return u.name.Value()
}

func (u *UserIndividualEntity) SetName(name string) error {
	nameValue, err := NewNameValue(name)
	if err != nil {
		return err
	}
	u.name = nameValue
	return nil
}

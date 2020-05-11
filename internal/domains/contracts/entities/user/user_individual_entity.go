package user

import (
	"time"
)

type UserIndividualEntity struct {
	*UserEntity
	name string
}

func NewUserIndividualEntity() *UserIndividualEntity {
	return &UserIndividualEntity{
		UserEntity: &UserEntity{},
		name:       "",
	}
}

func NewUserIndividualEntityWithData(id int, name string, createdAt time.Time, updatedAt time.Time) *UserIndividualEntity {
	userIndividual := NewUserIndividualEntity()
	userIndividual.id = id
	userIndividual.name = name
	userIndividual.createdAt = createdAt
	userIndividual.updatedAt = updatedAt

	return userIndividual
}

// 保持データをセットし直す
func (u *UserIndividualEntity) LoadData(id int, name string, createdAt time.Time, updatedAt time.Time) {
	u.id = id
	u.name = name
	u.createdAt = createdAt
	u.updatedAt = updatedAt
}

func (u *UserIndividualEntity) Name() string {
	return u.name
}

func (u *UserIndividualEntity) SetName(name string) {
	u.name = name
}

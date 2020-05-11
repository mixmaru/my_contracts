package user

import (
	"time"
)

type UserCorporationEntity struct {
	*UserEntity
	contactPersonName string //担当者名
	presidentName     string //社長名
}

func NewUserCorporationEntity() *UserCorporationEntity {
	return &UserCorporationEntity{
		UserEntity: &UserEntity{},
	}
}

func NewUserCorporationEntityWithData(id int, contractPersonName, presidentName string, createdAt, updatedAt time.Time) *UserCorporationEntity {
	user := NewUserCorporationEntity()
	user.id = id
	user.SetContactPersonName(contractPersonName)
	user.SetPresidentName(presidentName)
	user.createdAt = createdAt
	user.updatedAt = updatedAt
	return user
}

func (u *UserCorporationEntity) SetContactPersonName(name string) {
	u.contactPersonName = name
}

func (u *UserCorporationEntity) SetPresidentName(name string) {
	u.presidentName = name
}

func (u *UserCorporationEntity) ContactPersonName() string {
	return u.contactPersonName
}

func (u *UserCorporationEntity) PresidentName() string {
	return u.presidentName
}

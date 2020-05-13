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

// 保持データをセットし直す
func (u *UserCorporationEntity) LoadData(id int, contractPersonName, presidentName string, createdAt, updatedAt time.Time) {
	u.id = id
	u.SetContactPersonName(contractPersonName)
	u.SetPresidentName(presidentName)
	u.createdAt = createdAt
	u.updatedAt = updatedAt
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

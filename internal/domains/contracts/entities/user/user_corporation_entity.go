package user

import (
	"time"
)

type UserCorporationEntity struct {
	UserEntity
	contactPersonName string //担当者名
	presidentName     string //社長名
	createdAt         time.Time
	updatedAt         time.Time
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

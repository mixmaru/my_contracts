package user

import (
	"time"
)

type UserCorporation struct {
	User
	contactPersonName string //担当者名
	presidentName     string //社長名
	createdAt         time.Time
	updatedAt         time.Time
}

func (u *UserCorporation) SetContactPersonName(name string) {
	u.contactPersonName = name
}

func (u *UserCorporation) SetPresidentName(name string) {
	u.presidentName = name
}

func (u *UserCorporation) ContactPersonName() string {
	return u.contactPersonName
}

func (u *UserCorporation) PresidentName() string {
	return u.presidentName
}

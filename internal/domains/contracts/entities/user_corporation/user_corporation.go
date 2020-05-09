package user_corporation

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"time"
)

type UserCorporation struct {
	user.User
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

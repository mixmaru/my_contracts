package user

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user/values"
	"time"
)

type UserCorporationEntity struct {
	*UserEntity
	contactPersonName values.ContactPersonNameValue //担当者名
	presidentName     values.PresidentNameValue     //社長名
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

func (u *UserCorporationEntity) SetContactPersonName(name string) error {
	nameValue, err := values.NewContactPersonNameValue(name)
	if err != nil {
		return err
	}
	u.contactPersonName = nameValue
	return nil
}

func (u *UserCorporationEntity) SetPresidentName(name string) error {
	nameValue, err := values.NewPresidentNameValue(name)
	if err != nil {
		return err
	}
	u.presidentName = nameValue
	return nil
}

func (u *UserCorporationEntity) ContactPersonName() string {
	return u.contactPersonName.Value()
}

func (u *UserCorporationEntity) PresidentName() string {
	return u.presidentName.Value()
}

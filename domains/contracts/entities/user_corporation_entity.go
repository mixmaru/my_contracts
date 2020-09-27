package entities

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities/values"
	"time"
)

type UserCorporationEntity struct {
	*UserEntity
	corporationName   values.CorporationNameValue   //会社名
	contactPersonName values.ContactPersonNameValue //担当者名
	presidentName     values.PresidentNameValue     //社長名
}

func NewUserCorporationEntity(corporationName, contactPersonName, presidentName string) (*UserCorporationEntity, error) {
	corporationNameValue, err := values.NewCorporationNameValue(corporationName)
	if err != nil {
		return nil, err
	}

	contactPersonNameValue, err := values.NewContactPersonNameValue(contactPersonName)
	if err != nil {
		return nil, err
	}

	presidentNameValue, err := values.NewPresidentNameValue(presidentName)
	if err != nil {
		return nil, err
	}

	return &UserCorporationEntity{
		UserEntity:        &UserEntity{},
		corporationName:   corporationNameValue,
		contactPersonName: contactPersonNameValue,
		presidentName:     presidentNameValue,
	}, nil
}

func NewUserCorporationEntityWithData(id int, corporationName, contractPersonName, presidentName string, createdAt, updatedAt time.Time) (*UserCorporationEntity, error) {
	user, err := NewUserCorporationEntity(corporationName, contractPersonName, presidentName)
	if err != nil {
		return nil, err
	}

	user.id = id
	user.createdAt = createdAt
	user.updatedAt = updatedAt
	return user, nil
}

// 保持データをセットし直す
func (u *UserCorporationEntity) LoadData(id int, corporationName, contractPersonName, presidentName string, createdAt, updatedAt time.Time) error {
	if u.UserEntity == nil {
		u.UserEntity = &UserEntity{}
	}
	u.id = id
	err := u.SetCorporationName(corporationName)
	if err != nil {
		return err
	}
	err = u.SetContactPersonName(contractPersonName)
	if err != nil {
		return err
	}
	err = u.SetPresidentName(presidentName)
	if err != nil {
		return err
	}
	u.createdAt = createdAt
	u.updatedAt = updatedAt
	return nil
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

func (u *UserCorporationEntity) SetCorporationName(name string) error {
	nameValue, err := values.NewCorporationNameValue(name)
	if err != nil {
		return err
	}
	u.corporationName = nameValue
	return nil
}

func (u *UserCorporationEntity) ContactPersonName() string {
	return u.contactPersonName.Value()
}

func (u *UserCorporationEntity) PresidentName() string {
	return u.presidentName.Value()
}

func (u *UserCorporationEntity) CorporationName() string {
	return u.corporationName.Value()
}

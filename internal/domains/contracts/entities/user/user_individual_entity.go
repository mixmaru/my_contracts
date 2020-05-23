package user

import (
	"github.com/pkg/errors"
	"time"
	"unicode/utf8"
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

// Name値オブジェクト
type Name struct {
	value string
}

type EmptyValidError struct {
	error
}

func NewName(value string) (Name, []error) {
	validateErrors := nameValidate(value)
	if len(validateErrors) > 0 {
		return Name{}, validateErrors
	}
	return Name{
		value: value,
	}, nil
}

func nameValidate(name string) []error {
	var validErrors []error
	if isEmpty(name) {
		validErrors = append(validErrors, EmptyValidError{errors.New("nameが空です")})
	}

	return validErrors
}

func isEmpty(name string) bool {
	if utf8.RuneCountInString(name) == 0 {
		return true
	} else {
		return false
	}
}

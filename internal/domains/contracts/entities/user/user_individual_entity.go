package user

import (
	"github.com/pkg/errors"
	"time"
	"unicode/utf8"
)

type UserIndividualEntity struct {
	*UserEntity
	name Name
}

func NewUserIndividualEntity(name string) (*UserIndividualEntity, []error) {
	nameValue, errors := NewName(name)
	if len(errors) > 0 {
		return nil, errors
	}

	return &UserIndividualEntity{
		UserEntity: &UserEntity{},
		name:       nameValue,
	}, errors
}

func NewUserIndividualEntityWithData(id int, name string, createdAt time.Time, updatedAt time.Time) (*UserIndividualEntity, []error) {
	userIndividual, errors := NewUserIndividualEntity(name)
	if len(errors) > 0 {
		return nil, errors
	}
	userIndividual.id = id
	userIndividual.createdAt = createdAt
	userIndividual.updatedAt = updatedAt

	return userIndividual, errors
}

// 保持データをセットし直す
func (u *UserIndividualEntity) LoadData(id int, name string, createdAt time.Time, updatedAt time.Time) []error {
	nameValue, errors := NewName(name)
	if len(errors) > 0 {
		return errors
	}
	u.id = id
	u.name = nameValue
	u.createdAt = createdAt
	u.updatedAt = updatedAt
	return errors
}

func (u *UserIndividualEntity) Name() string {
	return u.name.value
}

func (u *UserIndividualEntity) SetName(name string) []error {
	nameValue, errors := NewName(name)
	if len(errors) > 0 {
		return errors
	}
	u.name = nameValue
	return errors
}

// Name値オブジェクト
type Name struct {
	value string
}

// から文字エラー
type EmptyValidError struct {
	error
}

// 文字数オーバーエラー
type OverLengthValidError struct {
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
	if isOverLength(name) {
		validErrors = append(validErrors, OverLengthValidError{errors.New("nameが50文字より多いです")})
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

func isOverLength(name string) bool {
	if utf8.RuneCountInString(name) <= 50 {
		return false
	} else {
		return true
	}
}

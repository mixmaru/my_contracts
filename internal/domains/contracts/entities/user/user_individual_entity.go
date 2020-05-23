package user

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
	"unicode/utf8"
)

type UserIndividualEntity struct {
	*UserEntity
	name Name
}

func NewUserIndividualEntity(name string) (*UserIndividualEntity, error) {
	nameValue, errs := NewName(name)
	if len(errs) > 0 {
		if len(errs) > 0 {
			messages := []string{}
			for _, err := range errs {
				messages = append(messages, err.Error())
			}
			return nil, errors.New(strings.Join(messages, ", "))
		}
	}

	return &UserIndividualEntity{
		UserEntity: &UserEntity{},
		name:       nameValue,
	}, nil
}

func NewUserIndividualEntityWithData(id int, name string, createdAt time.Time, updatedAt time.Time) (*UserIndividualEntity, error) {
	userIndividual, err := NewUserIndividualEntity(name)
	if err != nil {
		return nil, err
	}
	userIndividual.id = id
	userIndividual.createdAt = createdAt
	userIndividual.updatedAt = updatedAt

	return userIndividual, nil
}

// 保持データをセットし直す
func (u *UserIndividualEntity) LoadData(id int, name string, createdAt time.Time, updatedAt time.Time) error {
	nameValue, errs := NewName(name)
	if len(errs) > 0 {
		if len(errs) > 0 {
			messages := []string{}
			for _, err := range errs {
				messages = append(messages, err.Error())
			}
			return errors.New(strings.Join(messages, ", "))
		}
	}
	u.id = id
	u.name = nameValue
	u.createdAt = createdAt
	u.updatedAt = updatedAt
	return nil
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
		validErrors = append(validErrors, OverLengthValidError{errors.New(fmt.Sprintf("nameが50文字より多いです。name: %v", name))})
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

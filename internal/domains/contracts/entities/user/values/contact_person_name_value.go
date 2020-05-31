package values

import (
	plain_err "errors"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// Name値オブジェクト
type ContactPersonNameValue struct {
	value string
}

func NewContactPersonNameValue(value string) (ContactPersonNameValue, error) {
	validateErrors := ContactPersonNameValidate(value)
	if len(validateErrors) > 0 {
		var msgs []string
		for _, msg := range validateErrors {
			msgs = append(msgs, msg.Error())
		}
		return ContactPersonNameValue{}, errors.New(fmt.Sprintf("Nameバリデーションエラー。%v", strings.Join(msgs, ", ")))
	}
	return ContactPersonNameValue{
		value: value,
	}, nil
}

func (v *ContactPersonNameValue) Value() string {
	return v.value
}

func ContactPersonNameValidate(name string) []error {
	var validErrors []error
	if isEmpty(name) {
		validErrors = append(validErrors, EmptyValidError{plain_err.New("nameが空です")})
	}
	if isOverLength(name) {
		validErrors = append(validErrors, OverLengthValidError{plain_err.New(fmt.Sprintf("nameが50文字より多いです。name: %v", name))})
	}

	return validErrors
}

//func isEmpty(name string) bool {
//	if utf8.RuneCountInString(name) == 0 {
//		return true
//	} else {
//		return false
//	}
//}
//
//func isOverLength(name string) bool {
//	if utf8.RuneCountInString(name) <= 50 {
//		return false
//	} else {
//		return true
//	}
//}
//
//// から文字エラー
//type EmptyValidError struct {
//	error
//}
//
//// 文字数オーバーエラー
//type OverLengthValidError struct {
//	error
//}

package validators

import (
	"strconv"
	"unicode/utf8"
)

// 空文字チェック
func IsEmptyString(str string) bool {
	if utf8.RuneCountInString(str) == 0 {
		return true
	} else {
		return false
	}
}

// から文字エラー
type EmptyStringValidError struct {
	error
}

func NewEmptyStringValidError(err error) *EmptyStringValidError {
	return &EmptyStringValidError{error: err}
}

// 文字数オーバーチェック
func IsOverLengthString(str string, maxNum int) bool {
	if utf8.RuneCountInString(str) <= maxNum {
		return false
	} else {
		return true
	}
}

// 文字数オーバーエラー
type OverLengthStringValidError struct {
	error
}

func NewOverLengthStringValidError(error error) *OverLengthStringValidError {
	return &OverLengthStringValidError{error: error}
}

// 数値チェック
func IsNumericString(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return false
	} else {
		return true
	}
}

// 数値チェックエラー
type NumericStringValidError struct {
	error
}

func NewNumericStringValidError(error error) *NumericStringValidError {
	return &NumericStringValidError{error: error}
}

// その他のエラー定義
// マイナス値エラー
type NegativeValidError struct {
	error
}

func NewNegativeValidError(error error) *NegativeValidError {
	return &NegativeValidError{error: error}
}

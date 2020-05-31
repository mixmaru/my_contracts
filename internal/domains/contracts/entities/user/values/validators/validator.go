package validators

import "unicode/utf8"

func IsEmptyString(str string) bool {
	if utf8.RuneCountInString(str) == 0 {
		return true
	} else {
		return false
	}
}

func IsOverLengthString(str string, maxNum int) bool {
	if utf8.RuneCountInString(str) <= maxNum {
		return false
	} else {
		return true
	}
}

// から文字エラー
type EmptyValidError struct {
	error
}

// 文字数オーバーエラー
type OverLengthValidError struct {
	error
}

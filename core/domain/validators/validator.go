package validators

import (
	"strconv"
	"unicode/utf8"
)

const (
	EmptyStringValidError      = iota // 空文字エラー
	OverLengthStringValidError        // 文字数オーバーエラー
	NumericStringValidError           // 数値文字列ではないエラー
	NegativeValidError                // マイナス値エラー
)

var validErrorText = map[int]string{
	EmptyStringValidError:      "空文字エラー",
	OverLengthStringValidError: "文字数オーバーエラー",
	NumericStringValidError:    "数値文字列ではないエラー",
	NegativeValidError:         "マイナス値エラー",
}

func ValidErrorText(errorConst int) string {
	return validErrorText[errorConst]
}

// 空文字チェック
func IsEmptyString(str string) bool {
	if utf8.RuneCountInString(str) == 0 {
		return true
	} else {
		return false
	}
}

// 文字数オーバーチェック
func IsOverLengthString(str string, maxNum int) bool {
	if utf8.RuneCountInString(str) <= maxNum {
		return false
	} else {
		return true
	}
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

package utils

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

const (
	NotDefined  string = "NotDefined"
	Test        string = "Test"
	Development string = "Development"
	Production  string = "Production"
)

func GetExecuteMode() (string, error) {
	executeModeEnv := os.Getenv("MY_CONTRACTS_EXECUTE_MODE")
	switch executeModeEnv {
	case "test":
		return Test, nil
	case "development":
		return Development, nil
	case "production":
		return Production, nil
	case "":
		if isGoTest() {
			return Test, nil
		} else {
			return Development, nil
		}
	default:
		return NotDefined, errors.New(fmt.Sprintf("環境変数MY_CONTRACTS_EXECUTE_MODEが考慮外 MY_CONTRACTS_EXECUTE_MODE: %+v", executeModeEnv))
	}
}

func isGoTest() bool {
	if flag.Lookup("test.v") != nil {
		return true
	} else {
		return false
	}
}

func CreateJstTime(year int, month time.Month, day, hour, min, sec, nsec int) time.Time {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return time.Date(year, month, day, hour, min, sec, nsec, jst)
}

func CreateJstLocation() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

// テスト用に重複しない商品名を作成するメソッド
func CreateUniqProductNameForTest() string {
	unixNano := time.Now().UnixNano()
	suffix := strconv.FormatInt(unixNano, 10)
	name := "商品" + suffix
	return name
}

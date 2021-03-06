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

/*
現在時刻のタイムスタンプの文字列を返す。
テストでユニーク制約によるテスト失敗を避けるためにつかう
*/
func CreateTimestampString() string {
	timestamp := time.Now().UnixNano()
	return strconv.Itoa(int(timestamp))
}

package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
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
		return Development, nil
	default:
		return "", errors.New(fmt.Sprintf("環境変数MY_CONTRACTS_EXECUTE_MODEが考慮外 MY_CONTRACTS_EXECUTE_MODE: %+v", executeModeEnv))
	}
}

package my_logger

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var logger *zap.Logger = nil

func GetLogger() (*zap.Logger, error) {
	if logger != nil {
		// すでにloggerがあればそれを返す
		return logger, nil
	}

	// 動作モード取得
	mode, err := utils.GetExecuteMode()
	if err != nil {
		return nil, errors.Wrap(err, "動作モード取得失敗")
	}

	switch mode {
	case utils.Production:
		logger, err = zap.NewProduction()
	case utils.Development:
	case utils.Test:
		logger, err = zap.NewDevelopment()
	default:
		return nil, errors.New(fmt.Sprintf("動作モードが考慮外。mode: %v", mode))
	}

	if err != nil {
		return nil, errors.Wrap(err, "loggerの作成失敗")
	}
	return logger, nil
}

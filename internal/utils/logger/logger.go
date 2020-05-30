package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var logger *zap.Logger = nil

func init() {
	zap.S().Infow("An info message", "iteration", 1)
}

func GetLogger() (*zap.Logger, error) {
	if logger == nil {
		// todo: 環境によってログ設定を切り替える
		var err error
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, errors.Wrap(err, "loggerの作成失敗")
		}
	}

	return logger, nil
}

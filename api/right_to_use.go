package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/utils/my_logger"
	"net/http"
)

/*
次期使用権の作成バッチ

params:
date string 実行日
*/
//func executeRecur(c echo.Context) error {
//	logger, err := my_logger.GetLogger()
//	if err != nil {
//		return err
//	}
//
//	validErrs := map[string][]string{}
//
//	// 実行日取得
//	executeDate, errMsg := getExecuteDate(c.FormValue("date"))
//	if errMsg != "" {
//		validErrs["date"] = []string{errMsg}
//	}
//
//	if len(validErrs) > 0 {
//		return c.JSON(http.StatusBadRequest, validErrs)
//	}
//
//	contractApp := application_service.NewContractApplicationService()
//	rightToUseDtos, err := contractApp.CreateNextRightToUse(executeDate)
//	if err != nil {
//		logger.Sugar().Errorw("使用権の時期更新に失敗。", "executeDate", executeDate, "err", err)
//		c.Error(err)
//		return err
//	}
//	return c.JSON(http.StatusCreated, rightToUseDtos)
//}

/*
期限切れ使用権権のアーカイブ処理バッチ
実行日以前に使用期限の切れた使用権データをアーカイブする

params:
date string 実行日
*/
func executeRightToUseArchive(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	validErrs := map[string][]string{}
	// 実行日取得
	executeDate, errMsg := getExecuteDate(c.FormValue("date"))
	if errMsg != "" {
		validErrs["date"] = []string{errMsg}
	}

	if len(validErrs) > 0 {
		return c.JSON(http.StatusBadRequest, validErrs)
	}

	contractApp := application_service.NewContractApplicationService()
	rightToUseDtos, err := contractApp.ArchiveExpiredRightToUse(executeDate)
	if err != nil {
		logger.Sugar().Errorw("期限切れ使用権のアーカイブに失敗。", "executeDate", executeDate, "アーカイブ処理された使用権", rightToUseDtos, "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message":             "期限切れ使用権のアーカイブに失敗。",
			"succeed_rightToUses": rightToUseDtos,
		})
	}
	return c.JSON(http.StatusCreated, rightToUseDtos)
}

package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/mixmaru/my_contracts/utils/my_logger"
	"net/http"
	"time"
)

// 請求実行バッチ
// params:
// date string 請求実行日。未指定なら当日指定となる
// curl http://localhost:1323/batches/bills/billing?date=20200601
func executeBilling(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	validErrs := map[string][]string{}

	// 実行日取得
	var executeDate time.Time
	executeDateStr := c.FormValue("date")
	jst := utils.CreateJstLocation()

	if executeDateStr == "" {
		// 日付指定がなければ現在日時を実行日とする
		executeDate = time.Now().In(jst)
	} else {
		// 日付指定があればそれをを実行日とする
		executeDate, err = time.ParseInLocation("20060102", executeDateStr, jst)
		if err != nil {
			// dateに変な値が渡された
			validErrs["date"] = []string{
				"YYYYMMDDの形式ではありません",
			}
		}
	}
	if len(validErrs) > 0 {
		return c.JSON(http.StatusBadRequest, validErrs)
	}

	billApp := application_service.NewBillApplicationService()
	billDtos, err := billApp.ExecuteBilling(executeDate)
	if err != nil {
		logger.Sugar().Errorw("請求実行に失敗。", "executeDate", executeDate, "err", err, "完了したbill", billDtos)
		c.Error(err)
		return err
	}
	return c.JSON(http.StatusCreated, billDtos)
}

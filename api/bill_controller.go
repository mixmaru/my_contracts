package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/core/application/bill/billing"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/mixmaru/my_contracts/utils/my_logger"
	"net/http"
	"time"
)

type BillController struct {
	billingUseCase billing.IBillBillingUseCase
}

func NewBillController(billingUseCase billing.IBillBillingUseCase) *BillController {
	return &BillController{billingUseCase: billingUseCase}
}

// 請求実行バッチ
// params:
// date string 請求実行日。未指定なら当日指定となる
// curl http://localhost:1323/batches/bills/billing?date=20200601
func (b *BillController) Billing(c echo.Context) error {
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

	response, err := b.billingUseCase.Handle(billing.NewBillBillingUseCaseRequest(executeDate))
	if err != nil {
		logger.Sugar().Errorw("請求実行に失敗。", "executeDate", executeDate, "err", err, "完了したbill", response.BillDtos)
		c.Error(err)
		return err
	}
	return c.JSON(http.StatusCreated, response.BillDtos)
}

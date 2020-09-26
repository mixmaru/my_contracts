package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/mixmaru/my_contracts/internal/utils/my_logger"
	"net/http"
	"time"
)

func main() {
	e := newRouter()

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Routerの初期化
func newRouter() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 顧客新規登録
	e.POST("/users/", saveUser)
	// 顧客情報取得
	e.GET("/users/:id", getUser)
	// 商品登録
	e.POST("/products/", saveProduct)
	// 商品情報取得
	e.GET("/products/:id", getProduct)
	// 契約登録
	e.POST("/contracts/", saveContract)
	// 契約情報取得
	e.GET("/contracts/:id", getContract)
	// 請求実行バッチ
	e.POST("/batches/bills/billing", executeBilling)

	return e
}

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
		logger.Sugar().Errorw("請求実行に失敗。", "executeDate", executeDate, "err", err)
		c.Error(err)
		return err
	}
	return c.JSON(http.StatusCreated, billDtos)
}

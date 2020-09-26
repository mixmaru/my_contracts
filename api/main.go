package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/mixmaru/my_contracts/internal/utils/my_logger"
	"net/http"
	"strconv"
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

func saveUser(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	userAppService := application_service.NewUserApplicationService()

	// 顧客タイプで登録処理を分岐
	userType := c.FormValue("type")
	switch userType {
	case "individual":
		name := c.FormValue("name")
		user, validErrs, err := userAppService.RegisterUserIndividual(name)

		if err != nil {
			logger.Sugar().Errorw("個人顧客データ登録に失敗。", "name", name, "err", err)
			c.Error(err)
			return err
		}
		if len(validErrs) > 0 {
			return c.JSON(http.StatusBadRequest, validErrs)
		}
		return c.JSON(http.StatusCreated, user)
	case "corporation":
		corporationName := c.FormValue("corporation_name")
		contactName := c.FormValue("contact_person_name")
		presidentName := c.FormValue("president_name")

		user, validErrs, err := userAppService.RegisterUserCorporation(corporationName, contactName, presidentName)
		if err != nil {
			logger.Sugar().Errorw("法人顧客データ登録に失敗。", "corporationName", corporationName, "contactName", contactName, "presidentName", presidentName, "err", err)
			c.Error(err)
			return err
		}
		if len(validErrs) > 0 {
			return c.JSON(http.StatusBadRequest, validErrs)
		}
		return c.JSON(http.StatusCreated, user)
	default:
		validErrorMessage := map[string][]string{
			"type": []string{
				"typeがindividualでもcorporationでもありません。",
			},
		}
		return c.JSON(http.StatusBadRequest, validErrorMessage)
	}
}

func getUser(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// idに変な値が渡された
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// サービスインスタンス化
	userAppService := application_service.NewUserApplicationService()
	// データ取得
	user, err := userAppService.GetUserById(userId)
	if err != nil {
		logger.Sugar().Errorw("顧客データ取得に失敗。", "userId", userId, "err", err)
		c.Error(err)
		return err
	}

	// データがない
	if user == nil {
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// 返却
	return c.JSON(http.StatusOK, user)
}

// 商品新規登録
// params:
// name string 商品名
// price string 価格
// curl -F "name=A商品" -F "price=10.1" http://localhost:1323/individual_users
func saveProduct(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	// Get name and email
	name := c.FormValue("name")
	price := c.FormValue("price")

	productAppService := application_service.NewProductApplicationService()
	product, validErrs, err := productAppService.Register(name, price)
	if err != nil {
		logger.Sugar().Errorw("商品データ登録に失敗。", "name", name, "price", price, "err", err)
		c.Error(err)
		return err
	}
	if len(validErrs) > 0 {
		return c.JSON(http.StatusBadRequest, validErrs)
	}

	return c.JSON(http.StatusCreated, product)
}

// 商品情報取得
// curl http://localhost:1323/products/1
func getProduct(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// idに変な値が渡された
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// サービスインスタンス化
	productAppService := application_service.NewProductApplicationService()
	// データ取得
	product, err := productAppService.Get(productId)
	if err != nil {
		logger.Sugar().Errorw("商品データ取得に失敗。", "productId", productId, "err", err)
		c.Error(err)
		return err
	}

	// データがない
	if product.Id == 0 {
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// 返却
	return c.JSON(http.StatusOK, product)
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

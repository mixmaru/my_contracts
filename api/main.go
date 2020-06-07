package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/internal/utils/my_logger"
	"net/http"
	"strconv"
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

	// 個人顧客新規登録
	e.POST("/individual_users/", saveIndividualUser)
	// 個人顧客情報取得
	e.GET("/individual_users/:id", getIndividualUser)
	// 個人顧客新規登録
	e.POST("/corporation_users/", saveCorporationUser)
	//e.GET("/users/:id", getUser)
	//e.PUT("/users/:id", updateUser)
	//e.DELETE("/users/:id", deleteUser)

	return e
}

// 個人顧客新規登録
// params:
// name string 個人顧客名
// curl -F "name=個人　太郎" http://localhost:1323/individual_users
func saveIndividualUser(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	// Get name and email
	name := c.FormValue("name")
	userAppService := application_service.NewUserApplicationService()
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
}

// 個人顧客情報取得
// params:
// name string 個人顧客名
// curl http://localhost:1323/individual_users/1
func getIndividualUser(c echo.Context) error {
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
	user, err := userAppService.GetUserIndividual(userId)
	if err != nil {
		logger.Sugar().Errorw("個人顧客データ取得に失敗。", "userId", userId, "err", err)
		c.Error(err)
		return err
	}

	// データがない
	if user.Id == 0 {
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// 返却
	return c.JSON(http.StatusOK, user)
}

// 法人顧客新規登録
// params:
// contact_name string 担当者名
// president_name string 社長名
// curl -F "contact_name=担当　太郎" -F "president_name=社長　太郎" http://localhost:1323/corporation_users/
func saveCorporationUser(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	contactName := c.FormValue("contact_name")
	presidentName := c.FormValue("president_name")

	userAppService := application_service.NewUserApplicationService()
	user, validErrs, err := userAppService.RegisterUserCorporation(contactName, presidentName)
	if err != nil {
		logger.Sugar().Errorw("法人顧客データ登録に失敗。", "contactName", contactName, "presidentName", presidentName, "err", err)
		c.Error(err)
		return err
	}

	if len(validErrs) > 0 {
		return c.JSON(http.StatusBadRequest, validErrs)
	}

	return c.JSON(http.StatusCreated, user)
}

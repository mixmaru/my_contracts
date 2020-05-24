package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service"
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
	e.POST("/individual_users", saveIndividualUser)
	// 個人顧客情報取得
	e.GET("/individual_users/:id", getIndividualUser)
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
	// Get name and email
	name := c.FormValue("name")
	userAppService := application_service.NewUserApplicationService()
	user, validErrs, err := userAppService.RegisterUserIndividual(name)
	if err != nil {
		c.Error(err)
		return err
	}
	if len(validErrs) > 0 {
		validMessages := map[string][]string{
			"name": []string{},
		}
		for _, err := range validErrs {
			validMessages["name"] = append(validMessages["name"], err.Error())
		}
		return c.JSON(http.StatusBadRequest, validMessages)
	}

	return c.JSON(http.StatusCreated, user)
}

// 個人顧客情報取得
// params:
// name string 個人顧客名
// curl http://localhost:1323/individual_users/1
func getIndividualUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err)
		return err
	}

	// サービスインスタンス化
	userAppService := application_service.NewUserApplicationService()
	// データ取得
	user, err := userAppService.GetUserIndividual(userId)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	// 返却
	return c.JSON(http.StatusOK, user)
}

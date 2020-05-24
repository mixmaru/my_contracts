package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service"
	"net/http"
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

	// Routes
	e.POST("/individual_users", saveIndividualUser)
	//e.GET("/users/:id", getUser)
	//e.PUT("/users/:id", updateUser)
	//e.DELETE("/users/:id", deleteUser)

	return e
}

// e.POST("/individual_users", saveUser)
// params:
// name string 個人顧客名
// curl -F "name=個人　太郎" http://localhost:1323/individual_users
func saveIndividualUser(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	userAppService := application_service.NewUserApplicationService()
	user, validErrs, err := userAppService.RegisterUserIndividual(name)
	if err != nil {
		return c.String(http.StatusInternalServerError, "処理に失敗しました。")
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

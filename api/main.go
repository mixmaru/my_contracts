package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mixmaru/my_contracts/core/application/products/create"
	user_create "github.com/mixmaru/my_contracts/core/application/users/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
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

	// controller
	productController := NewProductController(create.NewProductCreateInteractor(db.NewProductRepository()))
	userController := NewUserController(user_create.NewUserIndividualCreateInteractor(db.NewUserRepository()))

	// 顧客新規登録
	e.POST("/users/", userController.Save)
	// 顧客情報取得
	e.GET("/users/:id", getUser)
	// 商品登録
	e.POST("/products/", productController.Crate)
	// 商品情報取得
	e.GET("/products/:id", productController.Get)
	// 契約登録
	e.POST("/contracts/", saveContract)
	// 契約情報取得
	e.GET("/contracts/:id", getContract)
	// 請求実行バッチ
	e.POST("/batches/bills/billing", executeBilling)
	// 使用権継続処理実行バッチ
	e.POST("/batches/right_to_uses/recur", executeRecur)
	// 有効期限切れ使用権のアーカイブ処理バッチ
	e.POST("/batches/right_to_uses/archive", executeRightToUseArchive)

	return e
}

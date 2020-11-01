package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mixmaru/my_contracts/core/application/contracts/archive_expired_right_to_use"
	create2 "github.com/mixmaru/my_contracts/core/application/contracts/create"
	"github.com/mixmaru/my_contracts/core/application/contracts/create_next_right_to_use"
	"github.com/mixmaru/my_contracts/core/application/contracts/get_by_id"
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
	userRep := db.NewUserRepository()
	productRep := db.NewProductRepository()
	contractRep := db.NewContractRepository()
	productController := NewProductController(create.NewProductCreateInteractor(productRep))
	userController := NewUserController(user_create.NewUserIndividualCreateInteractor(userRep))
	contractController := NewContractController(
		create2.NewContractCreateInteractor(userRep, productRep, contractRep),
		get_by_id.NewContractGetByIdInteractor(contractRep, productRep, userRep),
		create_next_right_to_use.NewContractCreateNextRightToUseInteractor(contractRep, productRep),
		archive_expired_right_to_use.NewContractArchiveExpiredRightToUseInteractor(contractRep),
	)

	// 顧客新規登録
	e.POST("/users/", userController.Save)
	// 顧客情報取得
	e.GET("/users/:id", userController.Get)
	// 商品登録
	e.POST("/products/", productController.Crate)
	// 商品情報取得
	e.GET("/products/:id", productController.Get)
	// 契約登録
	e.POST("/contracts/", contractController.Create)
	// 契約情報取得
	e.GET("/contracts/:id", contractController.GetById)
	// 請求実行バッチ
	e.POST("/batches/bills/billing", executeBilling)
	// 使用権継続処理実行バッチ
	e.POST("/batches/right_to_uses/recur", contractController.CreateNextRightToUse)
	// 有効期限切れ使用権のアーカイブ処理バッチ
	e.POST("/batches/right_to_uses/archive", contractController.ArchiveExpiredRightToUse)

	return e
}

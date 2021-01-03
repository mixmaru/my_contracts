package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mixmaru/my_contracts/core/application/bill/billing"
	"github.com/mixmaru/my_contracts/core/application/contracts/archive_expired_right_to_use"
	create2 "github.com/mixmaru/my_contracts/core/application/contracts/create"
	"github.com/mixmaru/my_contracts/core/application/contracts/create_next_right_to_use"
	"github.com/mixmaru/my_contracts/core/application/contracts/get_by_id"
	create3 "github.com/mixmaru/my_contracts/core/application/customer_property_type/create"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type/get_all"
	get_by_id2 "github.com/mixmaru/my_contracts/core/application/customer_property_type/get_by_id"
	create4 "github.com/mixmaru/my_contracts/core/application/customer_type/create"
	get_by_id3 "github.com/mixmaru/my_contracts/core/application/customer_type/get_by_id"
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

	userRep := db.NewUserRepository()
	productRep := db.NewProductRepository()
	contractRep := db.NewContractRepository()
	billRep := db.NewBillRepository()
	customerPropertyTypeRep := db.NewCustomerPropertyTypeRepository()
	customerTypeRep := db.NewCustomerTypeRepository()

	// controller
	productController := NewProductController(create.NewProductCreateInteractor(productRep))
	userController := NewUserController(user_create.NewUserIndividualCreateInteractor(userRep))
	contractController := NewContractController(
		create2.NewContractCreateInteractor(userRep, productRep, contractRep),
		get_by_id.NewContractGetByIdInteractor(contractRep, productRep, userRep),
		create_next_right_to_use.NewContractCreateNextRightToUseInteractor(contractRep, productRep),
		archive_expired_right_to_use.NewContractArchiveExpiredRightToUseInteractor(contractRep),
	)
	billController := NewBillController(billing.NewBillBillingInteractor(productRep, contractRep, billRep))
	customerPropertyTypeController := NewCustomerPropertyTypeController(
		create3.NewCustomerPropertyTypeCreateInteractor(customerPropertyTypeRep),
		get_by_id2.NewCustomerPropertyTypeGetByIdInteractor(customerPropertyTypeRep),
		get_all.NewCustomerPropertyTypeGetAllInteractor(customerPropertyTypeRep),
	)
	customerTypeController := NewCustomerTypeController(
		create4.NewCustomerTypeCreateInteractor(customerTypeRep, customerPropertyTypeRep),
		get_by_id3.NewCustomerTypeGetByIdInteractor(customerTypeRep, customerPropertyTypeRep),
	)
	customerController := NewCustomerController()

	// 顧客新規登録
	e.POST("/users/", userController.Save)
	// 顧客情報取得
	e.GET("/users/:id", userController.Get)
	// カスタマータイプ新規登録
	e.POST("/customer_types/", customerTypeController.Create)
	// カスタマータイプ取得
	e.GET("/customer_types/:id", customerTypeController.GetById)
	// カスタマープロパティタイプ新規登録
	e.POST("/customer_property_types/", customerPropertyTypeController.Create)
	// カスタマープロパティタイプ取得
	e.GET("/customer_property_types/:id", customerPropertyTypeController.GetById)
	// カスタマープロパティタイプ全取得
	e.GET("/customer_property_types/", customerPropertyTypeController.GetAll)
	// カスタマー新規登録
	e.POST("/customer/", customerController.Create)
	// 商品登録
	e.POST("/products/", productController.Crate)
	// 商品情報取得
	e.GET("/products/:id", productController.Get)
	// 契約登録
	e.POST("/contracts/", contractController.Create)
	// 契約情報取得
	e.GET("/contracts/:id", contractController.GetById)
	// 請求実行バッチ
	e.POST("/batches/bills/billing", billController.Billing)
	// 使用権継続処理実行バッチ
	e.POST("/batches/right_to_uses/recur", contractController.CreateNextRightToUse)
	// 有効期限切れ使用権のアーカイブ処理バッチ
	e.POST("/batches/right_to_uses/archive", contractController.ArchiveExpiredRightToUse)

	return e
}

package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
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

	// 個人顧客新規登録
	e.POST("/individual_users/", saveIndividualUser)
	// 個人顧客情報取得
	e.GET("/individual_users/:id", getIndividualUser)
	// 個人顧客新規登録
	e.POST("/corporation_users/", saveCorporationUser)
	// 法人顧客情報取得
	e.GET("/corporation_users/:id", getCorporationUser)
	// 商品登録
	e.POST("/products/", saveProduct)
	// 商品情報取得
	e.GET("/products/:id", getProduct)
	// 契約登録
	e.POST("/contracts/", saveContract)
	// 契約情報取得
	e.GET("/contracts/:id", getContract)

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
// contact_person_name string 担当者名
// president_name string 社長名
// curl -F "contact_person_name=担当　太郎" -F "president_name=社長　太郎" http://localhost:1323/corporation_users/
func saveCorporationUser(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	contactName := c.FormValue("contact_person_name")
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

// 法人顧客情報取得
// params:
// contact_person_name string 担当者名
// president_name string 社長名
// curl http://localhost:1323/individual_users/1
func getCorporationUser(c echo.Context) error {
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
	user, err := userAppService.GetUserCorporation(userId)
	if err != nil {
		logger.Sugar().Errorw("法人顧客データ取得に失敗。", "userId", userId, "err", err)
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

// 契約新規登録
// params:
// user_id string
// product_id string
// curl -F "user_id=1" -F "product_id=2" http://localhost:1323/contracts
func saveContract(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	// Get name and email
	userId, err := strconv.Atoi(c.FormValue("user_id"))
	if err != nil {
		// user_idに変な値が渡された
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}
	productId, err := strconv.Atoi(c.FormValue("product_id"))
	if err != nil {
		// product_idに変な値が渡された
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	app := application_service.NewContractApplicationService()
	contract, validErrs, err := app.Register(userId, productId)
	if err != nil {
		logger.Sugar().Errorw("契約データ登録に失敗。", "userId", userId, "productId", productId, "err", err)
		c.Error(err)
		return err
	}
	if len(validErrs) > 0 {
		return c.JSON(http.StatusBadRequest, validErrs)
	}

	return c.JSON(http.StatusCreated, contract)
}

// 商品情報取得
// curl http://localhost:1323/contracts/1
func getContract(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	contractId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// idに変な値が渡された
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// サービスインスタンス化
	//productAppService := application_service.NewProductApplicationService()
	contractAppService := application_service.NewContractApplicationService()
	// データ取得
	contract, product, user, err := contractAppService.GetById(contractId)
	if err != nil {
		logger.Sugar().Errorw("商品データ取得に失敗。", "contractId", contractId, "err", err)
		c.Error(err)
		return err
	}

	// データがない
	if contract.Id == 0 {
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// 返却データを用意
	switch user.(type) {
	case data_transfer_objects.UserIndividualDto:
		retContract := newContractDataForUserIndividual(contract.Id, product, user.(data_transfer_objects.UserIndividualDto), contract.CreatedAt, contract.UpdatedAt)
		return c.JSON(http.StatusOK, retContract)
	case data_transfer_objects.UserCorporationDto:
		retContract := newContractDataForUserCorporation(contract.Id, product, user.(data_transfer_objects.UserCorporationDto), contract.CreatedAt, contract.UpdatedAt)
		return c.JSON(http.StatusOK, retContract)
	default:
		logger.Sugar().Errorw("商品データ取得に失敗。userDtoが想定の型ではない。", "user", user, "err", err)
		c.Error(err)
		return err
	}
}

type contractDataForUserCorporation struct {
	contractData
	User data_transfer_objects.UserCorporationDto
}

func newContractDataForUserCorporation(id int, product data_transfer_objects.ProductDto, user data_transfer_objects.UserCorporationDto, createdAt time.Time, updatedAt time.Time) contractDataForUserCorporation {
	c := contractDataForUserCorporation{}
	c.Id = id
	c.User = user
	c.Product = product
	c.CreatedAt = createdAt
	c.UpdatedAt = updatedAt
	return c
}

type contractDataForUserIndividual struct {
	contractData
	User data_transfer_objects.UserIndividualDto
}

func newContractDataForUserIndividual(id int, product data_transfer_objects.ProductDto, user data_transfer_objects.UserIndividualDto, createdAt time.Time, updatedAt time.Time) contractDataForUserIndividual {
	c := contractDataForUserIndividual{}
	c.Id = id
	c.User = user
	c.Product = product
	c.CreatedAt = createdAt
	c.UpdatedAt = updatedAt
	return c
}

type contractData struct {
	Id        int
	Product   data_transfer_objects.ProductDto
	CreatedAt time.Time
	UpdatedAt time.Time
}

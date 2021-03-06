package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/core/application/products/create"
	"github.com/mixmaru/my_contracts/core/application/products/get"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils/my_logger"
	"net/http"
	"strconv"
)

type ProductController struct {
	createUseCase create.IProductCreateUseCase
}

func NewProductController(createUseCase create.IProductCreateUseCase) *ProductController {
	return &ProductController{createUseCase: createUseCase}
}

// 商品新規登録
// params:
// name string 商品名
// price string 価格
// curl -F "name=A商品" -F "price=10.1" http://localhost:1323/individual_users
func (con *ProductController) Crate(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	request := create.NewProductCreateUseCaseRequest(c.FormValue("name"), c.FormValue("price"))
	response, err := con.createUseCase.Handle(request)
	if err != nil {
		logger.Sugar().Errorw("商品データ登録に失敗。", "name", request.Name, "price", request.Price, "err", err)
		c.Error(err)
		return err
	}
	if len(response.ValidationError) > 0 {
		return c.JSON(http.StatusBadRequest, response.ValidationError)
	}

	return c.JSON(http.StatusCreated, response.ProductDto)
}

// 商品情報取得
// curl http://localhost:1323/products/1
func (con *ProductController) Get(c echo.Context) error {
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
	interactor := get.NewProductGetInteractor(db.NewProductRepository())
	// データ取得
	response, err := interactor.Handle(get.NewProductGetUseCaseRequest(productId))
	if err != nil {
		logger.Sugar().Errorw("商品データ取得に失敗。", "productId", productId, "err", err)
		c.Error(err)
		return err
	}

	// データがない
	if response.ProductDto.Id == 0 {
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// 返却
	return c.JSON(http.StatusOK, response.ProductDto)
}

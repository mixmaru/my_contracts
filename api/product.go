package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/internal/utils/my_logger"
	"net/http"
	"strconv"
)

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

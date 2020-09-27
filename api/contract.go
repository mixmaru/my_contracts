package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/utils/my_logger"
	"net/http"
	"strconv"
	"time"
)

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

	validErrs := map[string][]string{}
	userId, err := strconv.Atoi(c.FormValue("user_id"))
	if err != nil {
		// user_idに変な値が渡された
		validErrs["user_id"] = []string{
			"数値ではありません",
		}
	}
	productId, err := strconv.Atoi(c.FormValue("product_id"))
	if err != nil {
		// product_idに変な値が渡された
		validErrs["product_id"] = []string{
			"数値ではありません",
		}
	}
	if len(validErrs) > 0 {
		return c.JSON(http.StatusBadRequest, validErrs)
	}

	app := application_service.NewContractApplicationService()
	contract, validErrs, err := app.Register(userId, productId, time.Now())
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
		retContract := newContractDataForUserIndividual(contract, product, user.(data_transfer_objects.UserIndividualDto), contract.CreatedAt, contract.UpdatedAt)
		return c.JSON(http.StatusOK, retContract)
	case data_transfer_objects.UserCorporationDto:
		retContract := newContractDataForUserCorporation(contract, product, user.(data_transfer_objects.UserCorporationDto), contract.CreatedAt, contract.UpdatedAt)
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

func newContractDataForUserCorporation(contract data_transfer_objects.ContractDto, product data_transfer_objects.ProductDto, user data_transfer_objects.UserCorporationDto, createdAt time.Time, updatedAt time.Time) contractDataForUserCorporation {
	c := contractDataForUserCorporation{}
	c.Id = contract.Id
	c.ContractDate = contract.ContractDate
	c.BillingStartDate = contract.BillingStartDate
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

func newContractDataForUserIndividual(contract data_transfer_objects.ContractDto, product data_transfer_objects.ProductDto, user data_transfer_objects.UserIndividualDto, createdAt time.Time, updatedAt time.Time) contractDataForUserIndividual {
	c := contractDataForUserIndividual{}
	c.Id = contract.Id
	c.ContractDate = contract.ContractDate
	c.BillingStartDate = contract.BillingStartDate
	c.User = user
	c.Product = product
	c.CreatedAt = createdAt
	c.UpdatedAt = updatedAt
	return c
}

type contractData struct {
	Id               int
	Product          data_transfer_objects.ProductDto
	ContractDate     time.Time
	BillingStartDate time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

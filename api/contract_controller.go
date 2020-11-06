package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/core/application/contracts"
	"github.com/mixmaru/my_contracts/core/application/contracts/archive_expired_right_to_use"
	"github.com/mixmaru/my_contracts/core/application/contracts/create"
	"github.com/mixmaru/my_contracts/core/application/contracts/create_next_right_to_use"
	"github.com/mixmaru/my_contracts/core/application/contracts/get_by_id"
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/application/users"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/mixmaru/my_contracts/utils/my_logger"
	"net/http"
	"strconv"
	"time"
)

type ContractController struct {
	createUseCase                   create.IContractCreateUseCase
	getByIdUseCase                  get_by_id.IContractGetByIdUseCase
	crateNextRightToUseUseCase      create_next_right_to_use.IContractCreateNextRightToUseUseCase
	archiveExpiredRightToUseUseCase archive_expired_right_to_use.IContractArchiveExpiredRightToUseUseCase
}

func NewContractController(createUseCase create.IContractCreateUseCase, getByIdUseCase get_by_id.IContractGetByIdUseCase, crateNextRightToUseUseCase create_next_right_to_use.IContractCreateNextRightToUseUseCase, archiveExpiredRightToUseUseCase archive_expired_right_to_use.IContractArchiveExpiredRightToUseUseCase) *ContractController {
	return &ContractController{createUseCase: createUseCase, getByIdUseCase: getByIdUseCase, crateNextRightToUseUseCase: crateNextRightToUseUseCase, archiveExpiredRightToUseUseCase: archiveExpiredRightToUseUseCase}
}

// 契約新規登録
// params:
// user_id string
// product_id string
// curl -F "user_id=1" -F "product_id=2" http://localhost:1323/contracts
func (cont *ContractController) Create(c echo.Context) error {
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

	response, err := cont.createUseCase.Handle(create.NewContractCreateUseCaseRequest(userId, productId, time.Now()))
	if err != nil {
		logger.Sugar().Errorw("契約データ登録に失敗。", "userId", userId, "productId", productId, "err", err)
		c.Error(err)
		return err
	}
	if len(response.ValidationErrors) > 0 {
		return c.JSON(http.StatusBadRequest, response.ValidationErrors)
	}

	return c.JSON(http.StatusCreated, response.ContractDto)
}

// 商品情報取得
// curl http://localhost:1323/contracts/1
func (cont *ContractController) GetById(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	contractId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// idに変な値が渡された
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// データ取得
	response, err := cont.getByIdUseCase.Handle(get_by_id.NewContractGetByIdUseCaseRequest(contractId))
	if err != nil {
		logger.Sugar().Errorw("商品データ取得に失敗。", "contractId", contractId, "err", err)
		c.Error(err)
		return err
	}

	// データがない
	if response.ContractDto.Id == 0 {
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// 返却データを用意
	switch response.UserDto.(type) {
	case users.UserIndividualDto:
		retContract := newContractDataForUserIndividual(response.ContractDto, response.ProductDto, response.UserDto.(users.UserIndividualDto), response.ContractDto.CreatedAt, response.ContractDto.UpdatedAt)
		return c.JSON(http.StatusOK, retContract)
	case users.UserCorporationDto:
		retContract := newContractDataForUserCorporation(response.ContractDto, response.ProductDto, response.UserDto.(users.UserCorporationDto), response.ContractDto.CreatedAt, response.ContractDto.UpdatedAt)
		return c.JSON(http.StatusOK, retContract)
	default:
		logger.Sugar().Errorw("商品データ取得に失敗。userDtoが想定の型ではない。", "user", response.UserDto, "err", err)
		c.Error(err)
		return err
	}

}

type contractDataForUserCorporation struct {
	contractData
	User users.UserCorporationDto
}

func newContractDataForUserCorporation(contract contracts.ContractDto, product products.ProductDto, user users.UserCorporationDto, createdAt time.Time, updatedAt time.Time) contractDataForUserCorporation {
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
	User users.UserIndividualDto
}

func newContractDataForUserIndividual(contract contracts.ContractDto, product products.ProductDto, user users.UserIndividualDto, createdAt time.Time, updatedAt time.Time) contractDataForUserIndividual {
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
	Product          products.ProductDto
	ContractDate     time.Time
	BillingStartDate time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

/*
次期使用権の作成バッチ

params:
date string 実行日
*/
func (cont *ContractController) CreateNextRightToUse(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	validErrs := map[string][]string{}

	// 実行日取得
	executeDate, errMsg := getExecuteDate(c.FormValue("date"))
	if errMsg != "" {
		validErrs["date"] = []string{errMsg}
	}

	if len(validErrs) > 0 {
		return c.JSON(http.StatusBadRequest, validErrs)
	}

	response, err := cont.crateNextRightToUseUseCase.Handle(create_next_right_to_use.NewContractCreateNextRightToUseUseCaseRequest(executeDate))
	if err != nil {
		logger.Sugar().Errorw("使用権の時期更新に失敗。", "executeDate", executeDate, "err", err)
		c.Error(err)
		return err
	}
	return c.JSON(http.StatusCreated, response.NextTermContracts)
}

func getExecuteDate(dateStr string) (executeDate time.Time, errMsg string) {
	jst := utils.CreateJstLocation()

	if dateStr == "" {
		// 日付指定がなければ現在日時を実行日とする
		executeDate = time.Now().In(jst)
	} else {
		// 日付指定があればそれをを実行日とする
		var err error
		executeDate, err = time.ParseInLocation("20060102", dateStr, jst)
		if err != nil {
			// dateに変な値が渡された
			return time.Time{}, "YYYYMMDDの形式ではありません"
		}
	}
	return executeDate, ""
}

/*
期限切れ使用権権のアーカイブ処理バッチ
実行日以前に使用期限の切れた使用権データをアーカイブする

params:
date string 実行日
*/
func (cont *ContractController) ArchiveExpiredRightToUse(c echo.Context) error {

	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	validErrs := map[string][]string{}
	// 実行日取得
	executeDate, errMsg := getExecuteDate(c.FormValue("date"))
	if errMsg != "" {
		validErrs["date"] = []string{errMsg}
	}

	if len(validErrs) > 0 {
		return c.JSON(http.StatusBadRequest, validErrs)
	}

	response, err := cont.archiveExpiredRightToUseUseCase.Handle(archive_expired_right_to_use.NewContractArchiveExpiredRightToUseUseCaseRequest(executeDate))
	if err != nil {
		logger.Sugar().Errorw("期限切れ使用権のアーカイブに失敗。", "executeDate", executeDate, "アーカイブ処理された使用権", response, "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message":             "期限切れ使用権のアーカイブに失敗。",
			"succeed_rightToUses": response.ArchivedRightToUses,
		})
	}
	return c.JSON(http.StatusCreated, response.ArchivedRightToUses)
}

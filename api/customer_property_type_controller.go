package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type/create"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type/get_by_id"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type/get_by_ids"
	"github.com/mixmaru/my_contracts/utils/my_logger"
	"net/http"
	"strconv"
)

type CustomerPropertyTypeController struct {
	createCustomerPropertyTypeUseCase   create.ICustomerPropertyTypeCreateUseCase
	getByIdsCustomerPropertyTypeUseCase get_by_ids.ICustomerPropertyTypeGetByIdsUseCase
	getByIdCustomerPropertyTypeUseCase  get_by_id.ICustomerPropertyTypeGetByIdUseCase
}

func NewCustomerPropertyTypeController(
	createCustomerPropertyTypeUseCase create.ICustomerPropertyTypeCreateUseCase,
	getByIdsCustomerPropertyTypeUseCase get_by_ids.ICustomerPropertyTypeGetByIdsUseCase,
	getByIdCustomerPropertyTypeUseCase get_by_id.ICustomerPropertyTypeGetByIdUseCase,
) *CustomerPropertyTypeController {
	return &CustomerPropertyTypeController{
		createCustomerPropertyTypeUseCase:   createCustomerPropertyTypeUseCase,
		getByIdsCustomerPropertyTypeUseCase: getByIdsCustomerPropertyTypeUseCase,
		getByIdCustomerPropertyTypeUseCase:  getByIdCustomerPropertyTypeUseCase,
	}
}

// カスタマープロパティタイプ新規登録
// params:
// name string カスタマープロパティ名（性別、好きな食べ物、住所、、とか）
// type string カスタマープロパティの型（string or numeric）
func (cont *CustomerPropertyTypeController) Create(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	name := c.FormValue("name")
	propertyType := c.FormValue("type")
	request := create.NewCustomerPropertyTypeCreateUseCaseRequest(name, propertyType)
	response, err := cont.createCustomerPropertyTypeUseCase.Handle(request)
	if err != nil {
		logger.Sugar().Errorw("カスタマープロパティタイプデータ登録に失敗。", "name", name, "type", propertyType, "err", err)
		c.Error(err)
		return err
	}

	if len(response.ValidationError) > 0 {
		return c.JSON(http.StatusBadRequest, response.ValidationError)
	}

	return c.JSON(http.StatusCreated, response.CustomerPropertyTypeDto)
}

// カスタマープロパティタイプ取得
// params:
// id int カスタマープロパティタイプID
func (cont *CustomerPropertyTypeController) GetById(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	aa := c.Param("id")
	println(aa)
	validErrs := map[string][]string{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// idに変な値が渡された
		validErrs["id"] = []string{
			"数値ではありません",
		}
	}
	if len(validErrs) > 0 {
		return c.JSON(http.StatusBadRequest, validErrs)
	}

	request := get_by_id.NewCustomerPropertyTypeGetByIdUseCaseRequest(id)
	response, err := cont.getByIdCustomerPropertyTypeUseCase.Handle(request)
	if err != nil {
		logger.Sugar().Errorw("カスタマープロパティタイプデータ取得に失敗。", "id", id, "err", err)
		c.Error(err)
		return err
	}

	if len(response.ValidationError) > 0 {
		return c.JSON(http.StatusBadRequest, response.ValidationError)
	}
	if response.CustomerPropertyTypeDto.Id == 0 {
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	return c.JSON(http.StatusOK, response.CustomerPropertyTypeDto)
}

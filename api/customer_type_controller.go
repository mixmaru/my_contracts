package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/core/application/customer_type/create"
	"github.com/mixmaru/my_contracts/core/application/customer_type/get_by_id"
	"github.com/mixmaru/my_contracts/utils/my_logger"
	"net/http"
	"strconv"
)

type CustomerTypeController struct {
	createUseCase  create.ICustomerTypeCreateUseCase
	getByIdUseCase get_by_id.ICustomerTypeGetByIdUseCase
}

func NewCustomerTypeController(
	createUseCase create.ICustomerTypeCreateUseCase,
	getByIdUseCase get_by_id.ICustomerTypeGetByIdUseCase,
) *CustomerTypeController {
	return &CustomerTypeController{createUseCase: createUseCase, getByIdUseCase: getByIdUseCase}
}

// カスタマータイプ新規登録
// params:
// name string カスタマータイプ名（お得意様、個人、法人、役所、NPO、、とか）
// property_type_ids string カスタマープロパティのid（"1, 5, 3, 10")
func (cont *CustomerTypeController) Create(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	params, err := c.FormParams()
	if err != nil {
		logger.Sugar().Errorw("パラメータ取得に失敗。", "echo", c, "err", err)
		c.Error(err)
		return err
	}

	validErrors := map[string][]string{}
	name, ok := params["name"]
	if !ok {
		validErrors["name"] = []string{"入力されていません"}
	}

	propertyIds, ok := params["customer_property_ids"]
	if !ok {
		validErrors["customer_property_ids"] = []string{"入力されていません"}
	}
	if len(validErrors) > 0 {
		return c.JSON(http.StatusBadRequest, validErrors)
	}

	propertyTypeIds := make([]int, 0, len(params["customer_property_ids"]))
	for _, idStr := range propertyIds {
		idStr, err := strconv.Atoi(idStr)
		if err != nil {
			validErrors["customer_property_ids"] = []string{"数値ではありません"}
			return c.JSON(http.StatusBadRequest, validErrors)
		}
		propertyTypeIds = append(propertyTypeIds, idStr)
	}

	request := create.NewCustomerTypeCreateUseCaseRequest(name[0], propertyTypeIds)
	response, err := cont.createUseCase.Handle(request)

	if err != nil {
		logger.Sugar().Errorw("カスタマータイプデータ登録に失敗。", "request", request, "err", err)
		c.Error(err)
		return err
	}

	if len(response.ValidationError) > 0 {
		return c.JSON(http.StatusBadRequest, response.ValidationError)
	}

	return c.JSON(http.StatusCreated, response.CustomerTypeDto)
	//func (cont *CustomerPropertyTypeController) Create(c echo.Context) error {
	//	logger, err := my_logger.GetLogger()
	//	if err != nil {
	//		return err
	//	}
	//
	//	name := c.FormValue("name")
	//	propertyType := c.FormValue("type")
	//	request := create.NewCustomerPropertyTypeCreateUseCaseRequest(name, propertyType)
	//	response, err := cont.createCustomerPropertyTypeUseCase.Handle(request)
	//	if err != nil {
	//		logger.Sugar().Errorw("カスタマープロパティタイプデータ登録に失敗。", "name", name, "type", propertyType, "err", err)
	//		c.Error(err)
	//		return err
	//	}
	//
	//	if len(response.ValidationError) > 0 {
	//		return c.JSON(http.StatusBadRequest, response.ValidationError)
	//	}
	//
	//	return c.JSON(http.StatusCreated, response.CustomerPropertyTypeDto)
	//}
}

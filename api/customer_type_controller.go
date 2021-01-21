package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/core/application/customer_type/create"
	"github.com/mixmaru/my_contracts/core/application/customer_type/get_by_id"
	"github.com/mixmaru/my_contracts/utils/my_logger"
	"net/http"
	"net/url"
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

	name, propertyTypeIds, validErrors := getParamsAndValidation(params)
	if len(validErrors) > 0 {
		return c.JSON(http.StatusBadRequest, validErrors)
	}

	request := create.NewCustomerTypeCreateUseCaseRequest(name, propertyTypeIds)
	response, err := cont.createUseCase.Handle(request)

	if err != nil {
		logger.Sugar().Errorw("カスタマータイプデータ登録に失敗。", "request", request, "err", err)
		c.Error(err)
		return err
	}

	if len(response.ValidationErrors) > 0 {
		return c.JSON(http.StatusBadRequest, response.ValidationErrors)
	}

	return c.JSON(http.StatusCreated, response.CustomerTypeDto)
}

func getParamsAndValidation(params url.Values) (name string, propertyTypeIds []int, validErrors map[string][]string) {
	validErrors = map[string][]string{}
	nameParam, ok := params["name"]
	if !ok {
		validErrors["name"] = []string{"入力されていません"}
	}

	propertyIdsParam, ok := params["customer_property_ids"]
	if !ok {
		validErrors["customer_property_ids"] = []string{"入力されていません"}
	}
	if len(validErrors) > 0 {
		return "", nil, validErrors
	}

	propertyTypeIds = make([]int, 0, len(propertyIdsParam))
	for _, idStr := range propertyIdsParam {
		idStr, err := strconv.Atoi(idStr)
		if err != nil {
			validErrors["customer_property_ids"] = []string{"数値ではありません"}
			return "", nil, validErrors
		}
		propertyTypeIds = append(propertyTypeIds, idStr)
	}

	return nameParam[0], propertyTypeIds, validErrors
}

// カスタマータイプ取得
// params:
// id int カスタマータイプId
func (cont *CustomerTypeController) GetById(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// 数値出なかった場合not found
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	request := get_by_id.NewCustomerTypeGetByIdUseCaseRequest(id)
	response, err := cont.getByIdUseCase.Handle(request)
	if err != nil {
		logger.Sugar().Errorw("カスタマータイプデータ取得に失敗。", "request", request, "err", err)
		c.Error(err)
		return err
	}

	if len(response.ValidationError) > 0 {
		return c.JSON(http.StatusBadRequest, response.ValidationError)
	}
	if response.CustomerTypeDto.Id == 0 {
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	return c.JSON(http.StatusOK, response.CustomerTypeDto)
}

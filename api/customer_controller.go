package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/core/application/customer/create"
	"github.com/mixmaru/my_contracts/core/application/customer/get_by_id"
	"github.com/mixmaru/my_contracts/utils/my_logger"
	"net/http"
	"strconv"
)

type CustomerController struct {
	createUseCase  create.ICustomerCreateUseCase
	getByIdUseCase get_by_id.ICustomerGetByIdUseCase
}

func NewCustomerController(createUseCase create.ICustomerCreateUseCase, getByIdUseCase get_by_id.ICustomerGetByIdUseCase) *CustomerController {
	return &CustomerController{createUseCase: createUseCase, getByIdUseCase: getByIdUseCase}
}

// カスタマー新規登録
// params:
// name string カスタマー名
// customer_type_id カスタマータイプID
// properties {"1": "男", "2": 20}
func (cont *CustomerController) Create(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	params := new(customerCreateParams)
	if err := c.Bind(params); err != nil {
		logger.Sugar().Errorw("パラメータ取得に失敗。", "echo", c, "err", err)
		c.Error(err)
		return err
	}

	request := create.NewCustomerCreateUseCaseRequest(params.Name, params.CustomerTypeId, params.Properties)
	response, err := cont.createUseCase.Handle(request)
	if err != nil {
		logger.Sugar().Errorw("カスタマー新規登録に失敗。", "request", request, "err", err)
		c.Error(err)
		return err
	}
	if len(response.ValidationErrors) > 0 {
		return c.JSON(http.StatusBadRequest, response.ValidationErrors)
	}
	return c.JSON(http.StatusCreated, response.CustomerDto)

	// カスタマータイプからプロパティ名をリスト取得する
	//propertyNames :=
	//プロパティ名をパラメータから取得する
	//とれなければ空文字にする
	//いんたラクターの生成
	//request構造体をつくる
	//name string
	//typeId int
	//properties map[string]string
	//response, validErrors, err := interactor.Handle()
	// バリデーションチェック
	// dtoを返却する

	//params, err := c.FormParams()
	//if err != nil {
	//	logger.Sugar().Errorw("パラメータ取得に失敗。", "echo", c, "err", err)
	//	c.Error(err)
	//	return err
	//}
	//
	//name, propertyTypeIds, validErrors := getParamsAndValidation(params)
	//if len(validErrors) > 0 {
	//	return c.JSON(http.StatusBadRequest, validErrors)
	//}
	//
	//request := create.NewCustomerTypeCreateUseCaseRequest(name, propertyTypeIds)
	//response, err := cont.createUseCase.Handle(request)
	//
	//if err != nil {
	//	logger.Sugar().Errorw("カスタマータイプデータ登録に失敗。", "request", request, "err", err)
	//	c.Error(err)
	//	return err
	//}
	//
	//if len(response.ValidationErrors) > 0 {
	//	return c.JSON(http.StatusBadRequest, response.ValidationErrors)
	//}
	//
	//return c.JSON(http.StatusCreated, response.CustomerTypeDto)
}

type customerCreateParams struct {
	Name           string              `json:"name"`
	CustomerTypeId int                 `json:"customer_type_id"`
	Properties     map[int]interface{} `json:"properties"`
}

func (cont *CustomerController) GetById(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	customerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// idに変な値が渡された
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// データ取得
	response, err := cont.getByIdUseCase.Handle(get_by_id.NewCustomerGetByIdUseCaseRequest(customerId))
	if err != nil {
		logger.Sugar().Errorw("カスタマーデータ取得に失敗。", "customerId", customerId, "err", err)
		c.Error(err)
		return err
	}

	// データがない
	if response.CustomerDto.Id == 0 {
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// 返却データを用意
	return c.JSON(http.StatusOK, response.CustomerDto)
}

//
//func getParamsAndValidation(params url.Values) (name string, propertyTypeIds []int, validErrors map[string][]string) {
//	validErrors = map[string][]string{}
//	nameParam, ok := params["name"]
//	if !ok {
//		validErrors["name"] = []string{"入力されていません"}
//	}
//
//	propertyIdsParam, ok := params["customer_property_ids"]
//	if !ok {
//		validErrors["customer_property_ids"] = []string{"入力されていません"}
//	}
//	if len(validErrors) > 0 {
//		return "", nil, validErrors
//	}
//
//	propertyTypeIds = make([]int, 0, len(propertyIdsParam))
//	for _, idStr := range propertyIdsParam {
//		idStr, err := strconv.Atoi(idStr)
//		if err != nil {
//			validErrors["customer_property_ids"] = []string{"数値ではありません"}
//			return "", nil, validErrors
//		}
//		propertyTypeIds = append(propertyTypeIds, idStr)
//	}
//
//	return nameParam[0], propertyTypeIds, validErrors
//}
//
//// カスタマータイプ取得
//// params:
//// id int カスタマータイプId
//func (cont *CustomerTypeController) GetById(c echo.Context) error {
//	logger, err := my_logger.GetLogger()
//	if err != nil {
//		return err
//	}
//
//	idParam := c.Param("id")
//	id, err := strconv.Atoi(idParam)
//	if err != nil {
//		// 数値出なかった場合not found
//		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
//	}
//
//	request := get_by_id.NewCustomerTypeGetByIdUseCaseRequest(id)
//	response, err := cont.getByIdUseCase.Handle(request)
//	if err != nil {
//		logger.Sugar().Errorw("カスタマータイプデータ取得に失敗。", "request", request, "err", err)
//		c.Error(err)
//		return err
//	}
//
//	if len(response.ValidationErrors) > 0 {
//		return c.JSON(http.StatusBadRequest, response.ValidationErrors)
//	}
//	if response.CustomerTypeDto.Id == 0 {
//		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
//	}
//
//	return c.JSON(http.StatusOK, response.CustomerTypeDto)
//}

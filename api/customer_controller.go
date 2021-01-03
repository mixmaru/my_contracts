package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomerController struct {
	//createUseCase create.ICustomerCreateUseCase
}

func NewCustomerController() *CustomerController {
	return &CustomerController{}
}

// カスタマー新規登録
// params:
// name string カスタマー名
// customer_type_id カスタマータイプID
// properties {"性別": "男", "年齢": 20}
func (cont *CustomerController) Create(c echo.Context) error {
	//logger, err := my_logger.GetLogger()
	//if err != nil {
	//	return err
	//}
	params := new(customerCreateParams)
	if err := c.Bind(params); err != nil {
		// todo: return 5xx error and write log
		return nil
	}
	return c.JSON(http.StatusOK, params)

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
	//if len(response.ValidationError) > 0 {
	//	return c.JSON(http.StatusBadRequest, response.ValidationError)
	//}
	//
	//return c.JSON(http.StatusCreated, response.CustomerTypeDto)
	return nil
}

type customerCreateParams struct {
	Name           string                 `json:"name"`
	CustomerTypeId int                    `json:"customer_type_id"`
	Properties     map[string]interface{} `json:"properties"`
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
//	if len(response.ValidationError) > 0 {
//		return c.JSON(http.StatusBadRequest, response.ValidationError)
//	}
//	if response.CustomerTypeDto.Id == 0 {
//		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
//	}
//
//	return c.JSON(http.StatusOK, response.CustomerTypeDto)
//}

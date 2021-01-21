package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mixmaru/my_contracts/core/application/users/create"
	"github.com/mixmaru/my_contracts/core/application/users/get"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils/my_logger"
	"net/http"
	"strconv"
)

type UserController struct {
	createIndividualUseCase create.IUserIndividualCreateUseCase
}

func NewUserController(createIndividualUseCase create.IUserIndividualCreateUseCase) *UserController {
	return &UserController{createIndividualUseCase: createIndividualUseCase}
}

func (u *UserController) Save(c echo.Context) error {
	logger, err := my_logger.GetLogger()
	if err != nil {
		return err
	}

	// 顧客タイプで登録処理を分岐
	userType := c.FormValue("type")
	switch userType {
	case "individual":
		name := c.FormValue("name")
		interactor := create.NewUserIndividualCreateInteractor(db.NewUserRepository())
		response, err := interactor.Handle(create.NewUserIndividualCreateUseCaseRequest(name))
		if err != nil {
			logger.Sugar().Errorw("個人顧客データ登録に失敗。", "name", name, "err", err)
			c.Error(err)
			return err
		}
		if len(response.ValidationErrors) > 0 {
			return c.JSON(http.StatusBadRequest, response.ValidationErrors)
		}
		return c.JSON(http.StatusCreated, response.UserDto)
	case "corporation":
		corporationName := c.FormValue("corporation_name")
		contactName := c.FormValue("contact_person_name")
		presidentName := c.FormValue("president_name")

		interactor := create.NewUserCorporationCreateInteractor(db.NewUserRepository())
		response, err := interactor.Handle(create.NewUserCorporationCreateUseCaseRequest(corporationName, contactName, presidentName))
		if err != nil {
			logger.Sugar().Errorw("法人顧客データ登録に失敗。", "corporationName", corporationName, "contactName", contactName, "presidentName", presidentName, "err", err)
			c.Error(err)
			return err
		}
		if len(response.ValidationErrors) > 0 {
			return c.JSON(http.StatusBadRequest, response.ValidationErrors)
		}
		return c.JSON(http.StatusCreated, response.UserDto)
	default:
		validErrorMessage := map[string][]string{
			"type": []string{
				"typeがindividualでもcorporationでもありません。",
			},
		}
		return c.JSON(http.StatusBadRequest, validErrorMessage)
	}
}

func (u *UserController) Get(c echo.Context) error {
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
	interactor := get.NewUserGetInteractor(db.NewUserRepository())
	// データ取得
	response, err := interactor.Handle(get.NewUserGetUseCaseRequest(userId))
	if err != nil {
		logger.Sugar().Errorw("顧客データ取得に失敗。", "userId", userId, "err", err)
		c.Error(err)
		return err
	}

	// データがない
	if response.UserDto == nil {
		return c.JSON(http.StatusNotFound, echo.ErrNotFound)
	}

	// 返却
	return c.JSON(http.StatusOK, response.UserDto)
}

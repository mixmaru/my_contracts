package create

import (
	"fmt"
	"github.com/mixmaru/my_contracts/core/application/users"
	entities "github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/mixmaru/my_contracts/core/domain/validators"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/domains/contracts/entities/values"
	"github.com/pkg/errors"
)

////// 個人顧客
type UserIndividualCreateInteractor struct {
	userRepository users.IUserRepository
}

func NewUserIndividualCreateInteractor(userRepository users.IUserRepository) *UserIndividualCreateInteractor {
	return &UserIndividualCreateInteractor{userRepository: userRepository}
}

func (u *UserIndividualCreateInteractor) Handle(request *UserIndividualCreateUseCaseRequest) (*UserIndividualCreateUseCaseResponse, error) {
	// 入力値バリデーション
	validationErrors, err := userIndividualValidation(request.Name)
	if err != nil {
		return NewUserIndividualCreateUseCaseResponse(users.UserIndividualDto{}, nil), err
	}
	if len(validationErrors) > 0 {
		return NewUserIndividualCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), nil
	}

	// エンティティ作成
	userEntity, err := entities.NewUserIndividualEntity(request.Name)
	if err != nil {
		return NewUserIndividualCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), err
	}

	// トランザクション開始
	dbMap, err := db.GetConnection()
	if err != nil {
		return NewUserIndividualCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), err
	}
	defer dbMap.Db.Close()
	tran, err := dbMap.Begin()
	if err != nil {
		return NewUserIndividualCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), errors.Wrapf(err, "トランザクション開始失敗")
	}

	// リポジトリ使って保存
	savedId, err := u.userRepository.SaveUserIndividual(userEntity, tran)
	if err != nil {
		return NewUserIndividualCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), err
	}

	// 保存データ再読込
	savedUserEntity, err := u.userRepository.GetUserIndividualById(savedId, tran)
	if err != nil {
		return NewUserIndividualCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), err
	}

	// コミット
	err = tran.Commit()
	if err != nil {
		return NewUserIndividualCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), errors.Wrapf(err, "コミット失敗")
	}
	userDto := users.NewUserIndividualDtoFromEntity(savedUserEntity)
	return NewUserIndividualCreateUseCaseResponse(userDto, validationErrors), nil
}

func userIndividualValidation(name string) (validationErrors map[string][]string, err error) {
	validationErrors = map[string][]string{}

	// 個人顧客名バリデーション
	nameValidErrors := values.NameValidate(name)
	if len(nameValidErrors) > 0 {
		validationErrors["name"] = []string{}
	}
	for _, validErr := range nameValidErrors {
		var errorMessage string
		switch validErr {
		case validators.EmptyStringValidError:
			errorMessage = "空です"
		case validators.OverLengthStringValidError:
			errorMessage = fmt.Sprintf("%v文字より多いです", values.NameMaxLength)
		default:
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。name: %v, validErrorText: %v", name, validators.ValidErrorText(validErr)))
		}
		validationErrors["name"] = append(validationErrors["name"], errorMessage)
	}

	return validationErrors, nil
}

////// 法人顧客
type UserCorporationCreateInteractor struct {
	userRepository users.IUserRepository
}

func NewUserCorporationCreateInteractor(userRepository users.IUserRepository) *UserCorporationCreateInteractor {
	return &UserCorporationCreateInteractor{userRepository: userRepository}
}

// 法人顧客を新規登録する
// 成功時、登録した法人顧客情報を返却する
func (u *UserCorporationCreateInteractor) Handle(request *UserCorporationCreateUseCaseRequest) (*UserCorporationCreateUseCaseResponse, error) {
	response := &UserCorporationCreateUseCaseResponse{}
	// 入力値バリデーション
	validationErrors, err := registerUserCorporationValidation(request)
	if err != nil {
		return response, err
	}
	if len(validationErrors) > 0 {
		response.ValidationErrors = validationErrors
		return response, nil
	}

	// リポジトリ登録用にデータ作成
	entity, err := entities.NewUserCorporationEntity(request.CorporationName, request.ContactPersonName, request.PresidentName)
	if err != nil {
		return response, err
	}

	// トランザクション開始
	dbMap, err := db.GetConnection()
	if err != nil {
		return response, err
	}
	defer dbMap.Db.Close()
	tran, err := dbMap.Begin()
	if err != nil {
		return response, errors.Wrap(err, "トランザクション開始失敗")
	}

	// リポジトリつかって保存実行
	savedId, err := u.userRepository.SaveUserCorporation(entity, tran)
	if err != nil {
		return response, errors.WithMessagef(err, "法人顧客データ登録失敗。entity: %v", entity)
	}

	// 再読込
	registeredUser, err := u.userRepository.GetUserCorporationById(savedId, tran)
	if err != nil {
		return response, errors.WithMessagef(err, "保存法人顧客データ再読込失敗。savedId: %v", savedId)
	}

	err = tran.Commit()
	if err != nil {
		return response, errors.WithMessagef(err, "法人顧客データCommit失敗。request: %v", request)
	}

	// 登録データを取得してdtoにつめる
	response.UserDto = users.NewUserCorporationDtoFromEntity(registeredUser)
	return response, nil
}

func registerUserCorporationValidation(request *UserCorporationCreateUseCaseRequest) (validationErrors map[string][]string, err error) {
	validationErrors = map[string][]string{}

	// 会社名バリデーション
	corporationNameValidErrors, err := corporationNameValidation(request.CorporationName)
	if err != nil {
		return nil, err
	}
	if len(corporationNameValidErrors) > 0 {
		validationErrors["corporation_name"] = corporationNameValidErrors
	}

	// 担当者名バリデーション
	contactPersonNameValidErrors, err := contactPersonValidation(request.ContactPersonName)
	if err != nil {
		return nil, err
	}
	if len(contactPersonNameValidErrors) > 0 {
		validationErrors["contact_person_name"] = contactPersonNameValidErrors
	}

	// 社長名バリデーション
	presidentNameValidErrors, err := presidentNameValidation(request.PresidentName)
	if err != nil {
		return nil, err
	}
	if len(presidentNameValidErrors) > 0 {
		validationErrors["president_name"] = presidentNameValidErrors
	}

	return validationErrors, nil
}

func corporationNameValidation(corporationName string) (validErrors []string, err error) {
	corporationNameValidErrors, err := values.CorporationNameValue{}.Validate(corporationName)
	if err != nil {
		return nil, err
	}
	for _, validError := range corporationNameValidErrors {
		var errorMessage string
		switch validError {
		case validators.EmptyStringValidError:
			errorMessage = "空です"
		case validators.OverLengthStringValidError:
			errorMessage = fmt.Sprintf("%v文字より多いです", values.MaxCorporationNameNum)
		default:
			return nil, errors.New(fmt.Sprintf("想定外エラー。corporation_name: %v, validErrorText: %v", corporationName, validators.ValidErrorText(validError)))
		}
		validErrors = append(validErrors, errorMessage)
	}
	return validErrors, nil
}

func contactPersonValidation(contactPersonName string) (validErrors []string, err error) {
	contactPersonNameValidErrors, err := values.ContactPersonNameValue{}.Validate(contactPersonName)
	if err != nil {
		return nil, err
	}

	for _, validError := range contactPersonNameValidErrors {
		var errorMessage string
		switch validError {
		case validators.EmptyStringValidError:
			errorMessage = "空です"
		case validators.OverLengthStringValidError:
			errorMessage = fmt.Sprintf("%v文字より多いです", values.MaxContactPersonNameNum)
		default:
			return nil, errors.New(fmt.Sprintf("想定外エラー。contact_person_name: %v, validErrorText: %v", contactPersonName, validators.ValidErrorText(validError)))
		}
		validErrors = append(validErrors, errorMessage)
	}
	return validErrors, nil
}

func presidentNameValidation(presidentName string) (validErrors []string, err error) {
	presidentNameValidErrors, err := values.PresidentNameValue{}.Validate(presidentName)
	if err != nil {
		return nil, err
	}

	for _, validError := range presidentNameValidErrors {
		var errorMessage string
		switch validError {
		case validators.EmptyStringValidError:
			errorMessage = "空です"
		case validators.OverLengthStringValidError:
			errorMessage = fmt.Sprintf("%v文字より多いです", values.MaxPresidentNameNum)
		default:
			return nil, errors.New(fmt.Sprintf("想定外エラー。president_name: %v, validErrorText: %v", presidentName, validators.ValidErrorText(validError)))
		}
		validErrors = append(validErrors, errorMessage)
	}
	return validErrors, nil
}

package application_service

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values/validators"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/pkg/errors"
)

type UserApplicationService struct {
	userRepository interfaces.IUserRepository
}

// 個人顧客を新規登録する
// 成功時、登録した個人顧客情報を返却する
func (s *UserApplicationService) RegisterUserIndividual(name string) (userIndividualDto data_transfer_objects.UserIndividualDto, validationErrors map[string][]string, err error) {
	// 入力値バリデーション
	validationErrors, err = userIndividualValidation(name)
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, nil, err
	}
	if len(validationErrors) > 0 {
		return data_transfer_objects.UserIndividualDto{}, validationErrors, nil
	}

	// エンティティ作成
	userEntity, err := entities.NewUserIndividualEntity(name)
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, nil, err
	}

	// トランザクション開始
	dbMap, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, nil, err
	}
	defer dbMap.Db.Close()
	tran, err := dbMap.Begin()
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, nil, errors.WithStack(err)
	}

	// リポジトリ使って保存
	savedId, err := s.userRepository.SaveUserIndividual(userEntity, tran)
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, nil, err
	}

	// 保存データ再読込
	savedUserEntity, err := s.userRepository.GetUserIndividualById(savedId, tran)
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, nil, err
	}

	// コミット
	err = tran.Commit()
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, nil, err
	}
	userDto := data_transfer_objects.NewUserIndividualDtoFromEntity(savedUserEntity)
	return userDto, nil, nil
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

// 法人顧客を新規登録する
// 成功時、登録した法人顧客情報を返却する
func (s *UserApplicationService) RegisterUserCorporation(corporationName, contactPersonName string, presidentName string) (userCorporationDto data_transfer_objects.UserCorporationDto, validationErrors map[string][]string, err error) {
	// 入力値バリデーション
	validationErrors, err = registerUserCorporationValidation(corporationName, contactPersonName, presidentName)
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, nil, err
	}
	if len(validationErrors) > 0 {
		return data_transfer_objects.UserCorporationDto{}, validationErrors, nil
	}

	// リポジトリ登録用にデータ作成
	entity, err := entities.NewUserCorporationEntity(corporationName, contactPersonName, presidentName)
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, validationErrors, err
	}

	// トランザクション開始
	dbMap, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, nil, err
	}
	defer dbMap.Db.Close()
	tran, err := dbMap.Begin()
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, nil, errors.Wrap(err, "トランザクション開始失敗")
	}

	// リポジトリつかって保存実行
	savedId, err := s.userRepository.SaveUserCorporation(entity, tran)
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, nil, errors.WithMessagef(err, "法人顧客データ登録失敗。contractPersonName: %v, presidentName: %v", contactPersonName, presidentName)
	}

	// 再読込
	registeredUser, err := s.userRepository.GetUserCorporationById(savedId, tran)
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, nil, errors.WithMessagef(err, "保存法人顧客データ再読込。contractPersonName: %v, presidentName: %v", contactPersonName, presidentName)
	}

	err = tran.Commit()
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, nil, errors.WithMessagef(err, "法人顧客データCommit失敗。contractPersonName: %v, presidentName: %v", contactPersonName, presidentName)
	}

	// 登録データを取得してdtoにつめる
	userDto := data_transfer_objects.NewUserCorporationDtoFromEntity(registeredUser)
	return userDto, nil, nil
}

func registerUserCorporationValidation(corporationName, contactPersonName string, presidentName string) (validationErrors map[string][]string, err error) {
	validationErrors = map[string][]string{}

	// 会社名バリデーション
	corporationNameValidErrors, err := values.CorporationNameValue{}.Validate(contactPersonName)
	if err != nil {
		return nil, err
	}
	if len(corporationNameValidErrors) > 0 {
		validationErrors["corporation_name"] = []string{}
	}
	for _, validError := range corporationNameValidErrors {
		var errorMessage string
		switch validError {
		case validators.EmptyStringValidError:
			errorMessage = "空です"
		case validators.OverLengthStringValidError:
			errorMessage = fmt.Sprintf("%v文字より多いです", values.MaxCorporationNameNum)
		default:
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。corporation_name: %v, validErrorText: %v", corporationName, validators.ValidErrorText(validError)))
		}
		validationErrors["corporation_name"] = append(validationErrors["corporation_name"], errorMessage)
	}

	// 担当者名バリデーション
	contactPersonNameValidErrors, err := values.ContactPersonNameValue{}.Validate(contactPersonName)
	if err != nil {
		return nil, err
	}

	if len(contactPersonNameValidErrors) > 0 {
		validationErrors["contact_person_name"] = []string{}
	}
	for _, validError := range contactPersonNameValidErrors {
		var errorMessage string
		switch validError {
		case validators.EmptyStringValidError:
			errorMessage = "空です"
		case validators.OverLengthStringValidError:
			errorMessage = fmt.Sprintf("%v文字より多いです", values.MaxContactPersonNameNum)
		default:
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。contact_person_name: %v, validErrorText: %v", contactPersonName, validators.ValidErrorText(validError)))
		}
		validationErrors["contact_person_name"] = append(validationErrors["contact_person_name"], errorMessage)
	}

	// 社長名バリデーション
	presidentNameValidErrors, err := values.PresidentNameValue{}.Validate(presidentName)
	if err != nil {
		return nil, err
	}

	if len(presidentNameValidErrors) > 0 {
		validationErrors["president_name"] = []string{}
	}
	for _, validError := range presidentNameValidErrors {
		//validationErrors["president_name"] = append(validationErrors["president_name"], validError.Error())
		var errorMessage string
		switch validError {
		case validators.EmptyStringValidError:
			errorMessage = "空です"
		case validators.OverLengthStringValidError:
			errorMessage = fmt.Sprintf("%v文字より多いです", values.MaxPresidentNameNum)
		default:
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。contact_person_name: %v, validErrorText: %v", contactPersonName, validators.ValidErrorText(validError)))
		}
		validationErrors["president_name"] = append(validationErrors["president_name"], errorMessage)
	}

	return validationErrors, nil
}

func (s *UserApplicationService) GetUserById(userId int) (usrDto interface{}, err error) {
	conn, err := db_connection.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	gotUser, err := s.userRepository.GetUserById(userId, conn)
	if err != nil {
		return nil, err
	}

	if gotUser == nil {
		// データがない場合
		return nil, nil
	}

	switch gotUser.(type) {
	case *entities.UserIndividualEntity:
		userDto := data_transfer_objects.NewUserIndividualDtoFromEntity(gotUser.(*entities.UserIndividualEntity))
		return userDto, nil
	case *entities.UserCorporationEntity:
		userDto := data_transfer_objects.NewUserCorporationDtoFromEntity(gotUser.(*entities.UserCorporationEntity))
		return userDto, nil
	default:
		return nil, errors.Errorf("考慮していないtypeが来た。type: %t, userId: %v", gotUser, userId)
	}
}

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
	userDto := createUserIndividualDtoFromEntity(savedUserEntity)
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
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。name: %v, validErrorText: %v", name, validators.ValidErrorTest(validErr)))
		}
		validationErrors["name"] = append(validationErrors["name"], errorMessage)
	}

	return validationErrors, nil
}

// 個人顧客情報を取得して返却する
func (s *UserApplicationService) GetUserIndividual(userId int) (data_transfer_objects.UserIndividualDto, error) {
	conn, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, err
	}
	defer conn.Db.Close()

	user, err := s.userRepository.GetUserIndividualById(userId, conn)
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, err
	}
	if user == nil {
		// データがない場合、空データ構造体を返す
		return data_transfer_objects.UserIndividualDto{}, nil
	} else {
		userDto := createUserIndividualDtoFromEntity(user)
		return userDto, nil
	}
}

func createUserIndividualDtoFromEntity(entity *entities.UserIndividualEntity) data_transfer_objects.UserIndividualDto {
	return data_transfer_objects.UserIndividualDto{
		Name: entity.Name(),
		BaseDto: data_transfer_objects.BaseDto{
			Id:        entity.Id(),
			CreatedAt: entity.CreatedAt(),
			UpdatedAt: entity.UpdatedAt(),
		},
	}
}

func createUserCorporationDtoFromEntity(entity *entities.UserCorporationEntity) data_transfer_objects.UserCorporationDto {
	return data_transfer_objects.UserCorporationDto{
		ContactPersonName: entity.ContactPersonName(),
		PresidentName:     entity.PresidentName(),
		BaseDto: data_transfer_objects.BaseDto{
			Id:        entity.Id(),
			CreatedAt: entity.CreatedAt(),
			UpdatedAt: entity.UpdatedAt(),
		},
	}
}

// 法人顧客を新規登録する
// 成功時、登録した法人顧客情報を返却する
func (s *UserApplicationService) RegisterUserCorporation(contactPersonName string, presidentName string) (userCorporationDto data_transfer_objects.UserCorporationDto, validationErrors map[string][]string, err error) {
	// 入力値バリデーション
	validationErrors, err = registerUserCorporationValidation(contactPersonName, presidentName)
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, nil, err
	}
	if len(validationErrors) > 0 {
		return data_transfer_objects.UserCorporationDto{}, validationErrors, nil
	}

	// リポジトリ登録用にデータ作成
	entity, err := entities.NewUserCorporationEntity(contactPersonName, presidentName)
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
	userDto := createUserCorporationDtoFromEntity(registeredUser)
	return userDto, nil, nil
}

func registerUserCorporationValidation(contactPersonName string, presidentName string) (validationErrors map[string][]string, err error) {
	validationErrors = map[string][]string{}

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
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。contact_person_name: %v, validErrorText: %v", contactPersonName, validators.ValidErrorTest(validError)))
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
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。contact_person_name: %v, validErrorText: %v", contactPersonName, validators.ValidErrorTest(validError)))
		}
		validationErrors["president_name"] = append(validationErrors["president_name"], errorMessage)
	}

	return validationErrors, nil
}

// 法人顧客情報を取得して返却する
func (s *UserApplicationService) GetUserCorporation(userId int) (data_transfer_objects.UserCorporationDto, error) {
	conn, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, err
	}
	defer conn.Db.Close()

	gotUser, err := s.userRepository.GetUserCorporationById(userId, conn)
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, err
	}

	if gotUser == nil {
		// データがない場合、空データ構造体を返す
		return data_transfer_objects.UserCorporationDto{}, nil
	} else {
		userDto := createUserCorporationDtoFromEntity(gotUser)
		return userDto, nil
	}
}

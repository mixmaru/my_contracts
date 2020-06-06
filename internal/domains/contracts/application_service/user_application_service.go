package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user/values"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/pkg/errors"
)

type UserApplicationService struct {
	userRepository interfaces.IUserRepository
}

// 個人顧客を新規登録する
// 成功時、登録した個人顧客情報を返却する
func (s *UserApplicationService) RegisterUserIndividual(name string) (data_transfer_objects.UserIndividualDto, ValidationErrors, error) {
	// 入力値バリデーション
	nameValidErrors := values.NameValidate(name)
	if len(nameValidErrors) > 0 {
		validErrors := ValidationErrors{
			"name": nameValidErrors,
		}
		return data_transfer_objects.UserIndividualDto{}, validErrors, nil
	}

	// エンティティ作成
	userEntity, err := user.NewUserIndividualEntity(name)
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, ValidationErrors{}, err
	}

	// トランザクション開始
	dbMap, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, ValidationErrors{}, err
	}
	defer dbMap.Db.Close()
	tran, err := dbMap.Begin()
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, ValidationErrors{}, errors.WithStack(err)
	}

	// リポジトリ使って保存
	userEntity, err = s.userRepository.SaveUserIndividual(userEntity, tran)
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, ValidationErrors{}, err
	}

	// コミット
	err = tran.Commit()
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, ValidationErrors{}, err
	}
	userDto := createUserIndividualDtoFromEntity(userEntity)
	return userDto, ValidationErrors{}, nil
}

// 個人顧客情報を取得して返却する
func (s *UserApplicationService) GetUserIndividual(userId int) (data_transfer_objects.UserIndividualDto, error) {
	user, err := s.userRepository.GetUserIndividualById(userId, nil)
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

func createUserIndividualDtoFromEntity(entity *user.UserIndividualEntity) data_transfer_objects.UserIndividualDto {
	return data_transfer_objects.UserIndividualDto{
		Id:        entity.Id(),
		Name:      entity.Name(),
		CreatedAt: entity.CreatedAt(),
		UpdatedAt: entity.UpdatedAt(),
	}
}

func createUserCorporationDtoFromEntity(entity *user.UserCorporationEntity) data_transfer_objects.UserCorporationDto {
	return data_transfer_objects.UserCorporationDto{
		Id:                entity.Id(),
		ContactPersonName: entity.ContactPersonName(),
		PresidentName:     entity.PresidentName(),
		CreatedAt:         entity.CreatedAt(),
		UpdatedAt:         entity.UpdatedAt(),
	}
}

// 法人顧客を新規登録する
// 成功時、登録した法人顧客情報を返却する
func (s *UserApplicationService) RegisterUserCorporation(contactPersonName string, presidentName string) (data_transfer_objects.UserCorporationDto, ValidationErrors, error) {
	// 入力値バリデーション
	validationErrors := ValidationErrors{}
	contactPersonNameValidErrors := values.ContactPersonNameValidate(contactPersonName)
	if len(contactPersonNameValidErrors) > 0 {
		validationErrors["contactPersonName"] = contactPersonNameValidErrors
	}
	presidentNameValidErrors := values.PresidentNameValidate(presidentName)
	if len(presidentNameValidErrors) > 0 {
		validationErrors["presidentName"] = presidentNameValidErrors
	}
	if len(validationErrors) > 0 {
		return data_transfer_objects.UserCorporationDto{}, validationErrors, nil
	}

	// リポジトリ登録用にデータ作成
	entity := user.NewUserCorporationEntity()
	entity.SetContactPersonName(contactPersonName)
	entity.SetPresidentName(presidentName)

	// トランザクション開始
	dbMap, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, ValidationErrors{}, err
	}
	defer dbMap.Db.Close()
	tran, err := dbMap.Begin()
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, ValidationErrors{}, errors.Wrap(err, "トランザクション開始失敗")
	}

	// リポジトリつかって保存実行
	registeredUser, err := s.userRepository.SaveUserCorporation(entity, tran)
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, ValidationErrors{}, errors.WithMessagef(err, "法人顧客データ登録失敗。contractPersonName: %v, presidentName: %v", contactPersonName, presidentName)
	}
	err = tran.Commit()
	if err != nil {
		return data_transfer_objects.UserCorporationDto{}, ValidationErrors{}, errors.WithMessagef(err, "法人顧客データCommit失敗。contractPersonName: %v, presidentName: %v", contactPersonName, presidentName)
	}

	// 登録データを取得してdtoにつめる
	userDto := createUserCorporationDtoFromEntity(registeredUser)
	return userDto, ValidationErrors{}, nil
}

type ValidationErrors = map[string][]error

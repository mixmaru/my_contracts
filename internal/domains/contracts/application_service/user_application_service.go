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
func (s *UserApplicationService) RegisterUserIndividual(name string) (data_transfer_objects.UserIndividualDto, ValidationError, error) {
	// 入力値バリデーション
	validErrors := values.NameValidate(name)
	if len(validErrors) > 0 {
		return data_transfer_objects.UserIndividualDto{}, validErrors, nil
	}

	// エンティティ作成
	userEntity, err := user.NewUserIndividualEntity(name)
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, []error{}, err
	}

	// トランザクション開始
	dbMap, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, []error{}, err
	}
	defer dbMap.Db.Close()
	tran, err := dbMap.Begin()
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, []error{}, errors.WithStack(err)
	}

	// リポジトリ使って保存
	userEntity, err = s.userRepository.SaveUserIndividual(userEntity, tran)
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, []error{}, err
	}

	// コミット
	err = tran.Commit()
	if err != nil {
		return data_transfer_objects.UserIndividualDto{}, []error{}, err
	}
	userDto := createUserDtoFromEntity(userEntity)
	return userDto, []error{}, nil
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
		userDto := createUserDtoFromEntity(user)
		return userDto, nil
	}
}

func createUserDtoFromEntity(entity *user.UserIndividualEntity) data_transfer_objects.UserIndividualDto {
	return data_transfer_objects.UserIndividualDto{
		Id:        entity.Id(),
		Name:      entity.Name(),
		CreatedAt: entity.CreatedAt(),
		UpdatedAt: entity.UpdatedAt(),
	}
}

// 個人顧客を新規登録する
// 成功時、登録した個人顧客情報を返却する
func (s *UserApplicationService) RegisterUserCorporation(name string) (data_transfer_objects.UserCorporationDto, ValidationError, error) {
	return data_transfer_objects.UserCorporationDto{}, nil, nil
}

type ValidationError = []error

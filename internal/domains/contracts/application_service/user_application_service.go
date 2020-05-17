package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/pkg/errors"
)

type UserApplicationService struct {
	userRepository interfaces.IUserRepository
}

func NewUserApplicationService(userRepository interfaces.IUserRepository) *UserApplicationService {
	return &UserApplicationService{
		userRepository: userRepository,
	}
}

// 個人顧客を新規登録する
// 成功時、userIDを返却する
func (s *UserApplicationService) RegisterUserIndividual(name string) (int, error) {
	// エンティティ作成
	userEntity := user.NewUserIndividualEntity()
	userEntity.SetName(name)

	// トランザクション開始
	dbMap, err := db_connection.GetConnection()
	if err != nil {
		return 0, err
	}
	defer dbMap.Db.Close()
	tran, err := dbMap.Begin()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	// リポジトリ使って保存
	userEntity, err = s.userRepository.SaveUserIndividual(userEntity, tran)
	if err != nil {
		return 0, err
	}

	// コミット
	err = tran.Commit()
	if err != nil {
		return 0, err
	}
	return userEntity.Id(), nil
}

// 個人顧客情報を取得して返却する
func (s *UserApplicationService) GetUserIndividual(userId int) (data_transfer_objects.UserIndividualDto, error) {
	userDto := data_transfer_objects.UserIndividualDto{}
	user, err := s.userRepository.GetUserIndividualById(userId, nil)
	if err != nil {
		return userDto, err
	}

	// データ詰め直し
	userDto.Id = user.Id()
	userDto.Name = user.Name()
	userDto.CreatedAt = user.CreatedAt()
	userDto.UpdatedAt = user.UpdatedAt()

	return userDto, nil
}

package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
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

	// リポジトリ使って保存
	err = s.userRepository.SaveUserIndividual(userEntity, tran)
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

// 個人顧客情報を取得する
// 成功時、userIDを返却する
func (s *UserApplicationService) GetUserIndividual(userId int) (*user.UserIndividualEntity, error) {
	user, err := s.userRepository.GetUserIndividualById(userId, nil)
	if err != nil {
		return nil, err
	}
	return user, nil
}

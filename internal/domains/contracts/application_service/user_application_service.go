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
	s.userRepository.SaveUserIndividual(userEntity, tran)

	// コミット
	err = tran.Commit()
	if err != nil {
		return 0, err
	}
	return userEntity.Id(), nil
}

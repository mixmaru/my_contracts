package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
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

	// リポジトリ使って保存
	//repository := NewRepo
	return 0, nil
}

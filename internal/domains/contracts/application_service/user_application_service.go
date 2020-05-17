package application_service

import "github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"

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
	return 0, nil
}

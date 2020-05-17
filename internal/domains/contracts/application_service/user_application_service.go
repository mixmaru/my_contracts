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

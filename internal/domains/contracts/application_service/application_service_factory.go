package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories"
)

// UserApplicationService
func NewUserApplicationService() *UserApplicationService {
	return &UserApplicationService{userRepository: repositories.NewUserRepository()}
}

func NewUserApplicationServiceWithMock(userRepository interfaces.IUserRepository) *UserApplicationService {
	return &UserApplicationService{userRepository: userRepository}
}

// ProductApplicationService
func NewProductApplicationService() *ProductApplicationService {
	return &ProductApplicationService{productRepository: repositories.NewProductRepository()}
}

func NewProductApplicationServiceWithMock(prodcutRepository interfaces.IProductRepository) *ProductApplicationService {
	return &ProductApplicationService{productRepository: prodcutRepository}
}

// ContractApplicationService
func NewContractApplicationService() *ContractApplicationService {
	return &ContractApplicationService{
		ContractRepository:   repositories.NewContractRepository(),
		UserRepository:       repositories.NewUserRepository(),
		ProductRepository:    repositories.NewProductRepository(),
		RightToUseRepository: repositories.NewRightToUseRepository(),
	}
}

func NewContractApplicationServiceWithMock(contractRepository interfaces.IContractRepository) *ContractApplicationService {
	return &ContractApplicationService{ContractRepository: contractRepository}
}

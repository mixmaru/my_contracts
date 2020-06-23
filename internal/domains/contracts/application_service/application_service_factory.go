package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories"
)

// UserApplicationService
func NewUserApplicationService() *UserApplicationService {
	return &UserApplicationService{userRepository: &repositories.UserRepository{}}
}

func NewUserApplicationServiceWithMock(userRepository interfaces.IUserRepository) *UserApplicationService {
	return &UserApplicationService{userRepository: userRepository}
}

// ProductApplicationService
func NewProductApplicationService() *ProductApplicationService {
	return &ProductApplicationService{productRepository: &repositories.ProductRepository{}}
}

func NewProductApplicationServiceWithMock(prodcutRepository interfaces.IProductRepository) *ProductApplicationService {
	return &ProductApplicationService{productRepository: prodcutRepository}
}

package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/product"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user"
)

// UserApplicationService
func NewUserApplicationService() *UserApplicationService {
	return &UserApplicationService{userRepository: &user.UserRepository{}}
}

func NewUserApplicationServiceWithMock(userRepository interfaces.IUserRepository) *UserApplicationService {
	return &UserApplicationService{userRepository: userRepository}
}

// ProductApplicationService
func NewProductApplicationService() *ProductApplicationService {
	return &ProductApplicationService{productRepository: &product.ProductRepository{}}
}

func NewProductApplicationServiceWithMock(prodcutRepository interfaces.IProductRepository) *ProductApplicationService {
	return &ProductApplicationService{productRepository: prodcutRepository}
}

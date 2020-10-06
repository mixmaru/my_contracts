package application_service

import (
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories"
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
		contractRepository: repositories.NewContractRepository(),
		userRepository:     repositories.NewUserRepository(),
		productRepository:  repositories.NewProductRepository(),
	}
}

func NewContractApplicationServiceWithMock(contractRepository interfaces.IContractRepository) *ContractApplicationService {
	return &ContractApplicationService{contractRepository: contractRepository}
}

// BillApplicationService
func NewBillApplicationService() *BillApplicationService {
	return &BillApplicationService{
		productRepository:  repositories.NewProductRepository(),
		contractRepository: repositories.NewContractRepository(),
		billRepository:     repositories.NewBillRepository(),
	}
}

func NewBillApplicationServiceWithMock(
	productRepository interfaces.IProductRepository,
	contractRepository interfaces.IContractRepository,
	billRepository interfaces.IBillRepository,
) *BillApplicationService {
	return &BillApplicationService{
		productRepository:  productRepository,
		contractRepository: contractRepository,
		billRepository:     billRepository,
	}
}

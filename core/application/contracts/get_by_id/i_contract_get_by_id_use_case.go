package get_by_id

import (
	"github.com/mixmaru/my_contracts/core/application/contracts"
	"github.com/mixmaru/my_contracts/core/application/products"
)

type IContractGetByIdUseCase interface {
	Handle(request *ContractGetByIdUseCaseRequest) (*ContractGetByIdUseCaseResponse, error)
}

type ContractGetByIdUseCaseRequest struct {
	ContractId int
}

func NewContractGetByIdUseCaseRequest(contractId int) *ContractGetByIdUseCaseRequest {
	return &ContractGetByIdUseCaseRequest{ContractId: contractId}
}

type ContractGetByIdUseCaseResponse struct {
	ContractDto contracts.ContractDto
	ProductDto  products.ProductDto
	UserDto     interface{}
}

func NewContractGetByIdUseCaseResponse(contractDto contracts.ContractDto, productDto products.ProductDto, userDto interface{}) *ContractGetByIdUseCaseResponse {
	return &ContractGetByIdUseCaseResponse{ContractDto: contractDto, ProductDto: productDto, UserDto: userDto}
}

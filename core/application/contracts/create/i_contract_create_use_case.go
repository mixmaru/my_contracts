package create

import (
	"github.com/mixmaru/my_contracts/core/application/contracts"
	"time"
)

type IContractCreateUseCase interface {
	Handle(request *ContractCreateUseCaseRequest) (*ContractCreateUseCaseResponse, error)
}

type ContractCreateUseCaseRequest struct {
	UserId           int
	ProductId        int
	ContractDateTime time.Time //契約日
}

func NewContractCreateUseCaseRequest(userId int, productId int, contractDateTime time.Time) *ContractCreateUseCaseRequest {
	return &ContractCreateUseCaseRequest{UserId: userId, ProductId: productId, ContractDateTime: contractDateTime}
}

type ContractCreateUseCaseResponse struct {
	ContractDto      contracts.ContractDto
	ValidationErrors map[string][]string
}

func NewContractCreateUseCaseResponse(contractDto contracts.ContractDto, validationErrors map[string][]string) *ContractCreateUseCaseResponse {
	return &ContractCreateUseCaseResponse{ContractDto: contractDto, ValidationErrors: validationErrors}
}

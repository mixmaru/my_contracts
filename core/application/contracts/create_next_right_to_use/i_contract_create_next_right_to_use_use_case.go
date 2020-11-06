package create_next_right_to_use

import (
	"github.com/mixmaru/my_contracts/core/application/contracts"
	"time"
)

type IContractCreateNextRightToUseUseCase interface {
	Handle(request *ContractCreateNextRightToUseUseCaseRequest) (*ContractCreateNextRightToUseUseCaseResponse, error)
}

type ContractCreateNextRightToUseUseCaseRequest struct {
	ExecuteDate time.Time
}

func NewContractCreateNextRightToUseUseCaseRequest(executeDate time.Time) *ContractCreateNextRightToUseUseCaseRequest {
	return &ContractCreateNextRightToUseUseCaseRequest{ExecuteDate: executeDate}
}

type ContractCreateNextRightToUseUseCaseResponse struct {
	NextTermContracts []contracts.ContractDto
}

func NewContractCreateNextRightToUseUseCaseResponse(nextTermContracts []contracts.ContractDto) *ContractCreateNextRightToUseUseCaseResponse {
	return &ContractCreateNextRightToUseUseCaseResponse{NextTermContracts: nextTermContracts}
}

package archive_expired_right_to_use

import (
	"github.com/mixmaru/my_contracts/core/application/contracts"
	"time"
)

type IContractArchiveExpiredRightToUseUseCase interface {
	Handle(request *ContractArchiveExpiredRightToUseUseCaseRequest) (*ContractArchiveExpiredRightToUseUseCaseResponse, error)
}

type ContractArchiveExpiredRightToUseUseCaseRequest struct {
	BaseDate time.Time
}

func NewContractArchiveExpiredRightToUseUseCaseRequest(baseDate time.Time) *ContractArchiveExpiredRightToUseUseCaseRequest {
	return &ContractArchiveExpiredRightToUseUseCaseRequest{BaseDate: baseDate}
}

type ContractArchiveExpiredRightToUseUseCaseResponse struct {
	ArchivedRightToUses []contracts.RightToUseDto
}

func NewContractArchiveExpiredRightToUseUseCaseResponse(archivedRightToUses []contracts.RightToUseDto) *ContractArchiveExpiredRightToUseUseCaseResponse {
	return &ContractArchiveExpiredRightToUseUseCaseResponse{ArchivedRightToUses: archivedRightToUses}
}

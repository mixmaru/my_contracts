package billing

import (
	"github.com/mixmaru/my_contracts/core/application/bill"
	"time"
)

type IBillBillingUseCase interface {
	Handle(request *BillBillingUseCaseRequest) (*BillBillingUseCaseResponse, error)
}

type BillBillingUseCaseRequest struct {
	ExecuteDate time.Time
}

func NewBillBillingUseCaseRequest(executeDate time.Time) *BillBillingUseCaseRequest {
	return &BillBillingUseCaseRequest{ExecuteDate: executeDate}
}

type BillBillingUseCaseResponse struct {
	BillDtos []bill.BillDto
}

func NewBillBillingUseCaseResponse(billDtos []bill.BillDto) *BillBillingUseCaseResponse {
	return &BillBillingUseCaseResponse{BillDtos: billDtos}
}

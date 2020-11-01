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

type BillBillingUseCaseResponse struct {
	BillDto bill.BillDto
}

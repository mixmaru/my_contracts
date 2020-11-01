package bill

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"time"
)

type BillDto struct {
	Id                      int
	BillingDate             time.Time
	UserId                  int
	PaymentConfirmed        bool
	PaymentConfirmedAt      time.Time
	TotalAmountExcludingTax string
	BillDetails             []BillDetailDto
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type BillDetailDto struct {
	Id            int
	RightToUseId  int
	BillingAmount string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewBillDtoFromEntity(entity *entities.BillAggregation) (BillDto, error) {
	paymentConfirmedAt, isNil, err := entity.PaymentConfirmedAt()
	if err != nil {
		return BillDto{}, err
	}
	var paymentConfirmed bool
	if isNil {
		paymentConfirmed = false
	} else {
		paymentConfirmed = true
	}
	amountExcludingTax := entity.TotalAmountExcludingTax()

	dto := BillDto{}
	dto.Id = entity.Id()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	dto.BillingDate = entity.BillingDate()
	dto.UserId = entity.UserId()
	dto.PaymentConfirmed = paymentConfirmed
	dto.PaymentConfirmedAt = paymentConfirmedAt
	dto.TotalAmountExcludingTax = amountExcludingTax.String()

	// detail
	details := entity.BillDetails()
	dto.BillDetails = make([]BillDetailDto, 0, len(details))
	for _, detail := range details {
		billingAmount := detail.BillingAmount()

		detailDto := BillDetailDto{}
		detailDto.Id = detail.Id()
		detailDto.CreatedAt = detail.CreatedAt()
		detailDto.UpdatedAt = detail.UpdatedAt()
		detailDto.RightToUseId = detail.RightToUseId()
		detailDto.BillingAmount = billingAmount.String()
		dto.BillDetails = append(dto.BillDetails, detailDto)
	}
	return dto, nil
}

package data_transfer_objects

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"time"
)

type ContractDto struct {
	UserId           int
	ProductId        int
	ContractDate     time.Time
	BillingStartDate time.Time
	RightToUseDtos   []RightToUseDto
	BaseDto
}

func NewContractDtoFromEntity(entity *entities.ContractEntity) ContractDto {
	dto := ContractDto{}
	dto.Id = entity.Id()
	dto.ProductId = entity.ProductId()
	dto.UserId = entity.UserId()
	dto.ContractDate = entity.ContractDate()
	dto.BillingStartDate = entity.BillingStartDate()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	return dto
}

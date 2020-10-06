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
	// rightToUseを作成する
	dto.RightToUseDtos = make([]RightToUseDto, 0, len(entity.RightToUses()))
	for _, rightToUseEntity := range entity.RightToUses() {
		dto.RightToUseDtos = append(dto.RightToUseDtos, NewRightToUseDtoFromEntity(rightToUseEntity))
	}
	return dto
}

type RightToUseDto struct {
	ValidFrom time.Time
	ValidTo   time.Time
	BaseDto
}

func NewRightToUseDtoFromEntity(entity *entities.RightToUseEntity) RightToUseDto {
	dto := RightToUseDto{}
	dto.Id = entity.Id()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	dto.ValidFrom = entity.ValidFrom()
	dto.ValidTo = entity.ValidTo()
	return dto
}

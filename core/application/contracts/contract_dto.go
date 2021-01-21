package contracts

import (
	"github.com/mixmaru/my_contracts/core/domain/models/contract"
	"time"
)

type ContractDto struct {
	Id               int
	UserId           int
	ProductId        int
	ContractDate     time.Time
	BillingStartDate time.Time
	RightToUseDtos   []RightToUseDto
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewContractDtoFromEntity(entity *contract.ContractEntity) ContractDto {
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
	Id        int
	ValidFrom time.Time
	ValidTo   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewRightToUseDtoFromEntity(entity *contract.RightToUseEntity) RightToUseDto {
	dto := RightToUseDto{}
	dto.Id = entity.Id()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	dto.ValidFrom = entity.ValidFrom()
	dto.ValidTo = entity.ValidTo()
	return dto
}

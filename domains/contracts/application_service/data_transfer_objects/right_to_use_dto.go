package data_transfer_objects

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"time"
)

type RightToUseDto struct {
	ContractId int
	ValidFrom  time.Time
	ValidTo    time.Time
	BaseDto
}

func NewRightToUseDtoFromEntity(contractId int, entity *entities.RightToUseEntity) RightToUseDto {
	dto := RightToUseDto{}
	dto.Id = entity.Id()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	dto.ContractId = contractId
	dto.ValidFrom = entity.ValidFrom()
	dto.ValidTo = entity.ValidTo()
	return dto
}

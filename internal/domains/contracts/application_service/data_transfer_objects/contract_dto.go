package data_transfer_objects

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
)

type ContractDto struct {
	UserId    int
	ProductId int
	BaseDto
}

func NewContractDtoFromEntity(entity *entities.ContractEntity) ContractDto {
	dto := ContractDto{}
	dto.Id = entity.Id()
	dto.ProductId = entity.ProductId()
	dto.UserId = entity.UserId()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	return dto
}

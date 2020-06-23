package data_transfer_objects

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
)

type ProductDto struct {
	Name  string
	Price decimal.Decimal
	BaseDto
}

func NewProductDtoFromEntity(entity *entities.ProductEntity) ProductDto {
	dto := ProductDto{}
	dto.Id = entity.Id()
	dto.Name = entity.Name()
	dto.Price = entity.Price()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	return dto
}

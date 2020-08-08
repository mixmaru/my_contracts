package data_transfer_objects

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
)

type ProductDto struct {
	Name  string
	Price string
	BaseDto
}

func NewProductDtoFromEntity(entity *entities.ProductEntity) ProductDto {
	price := entity.MonthlyPrice()

	dto := ProductDto{}
	dto.Id = entity.Id()
	dto.Name = entity.Name()
	dto.Price = price.String()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	return dto
}

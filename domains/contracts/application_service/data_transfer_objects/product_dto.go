package data_transfer_objects

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/pkg/errors"
)

type ProductDto struct {
	Name  string
	Price string
	BaseDto
}

func NewProductDtoFromEntity(entity *entities.ProductEntity) ProductDto {
	price, exist := entity.MonthlyPrice()
	if !exist {
		errors.Errorf("月額料金が取得できなかった。%v", entity)
	}

	dto := ProductDto{}
	dto.Id = entity.Id()
	dto.Name = entity.Name()
	dto.Price = price.String()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	return dto
}

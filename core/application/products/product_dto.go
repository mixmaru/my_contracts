package products

import (
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/pkg/errors"
	"time"
)

type ProductDto struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Price     string
}

func NewProductDtoFromEntity(entity *product.ProductEntity) ProductDto {
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

package create

import (
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/pkg/errors"
	"time"
)

type IProductCreateUseCase interface {
	Handle(request ProductCreateUseCaseRequest) (ProductCreateUseCaseResponse, error)
}

type ProductCreateUseCaseRequest struct {
	Name  string
	Price string
}

func NewProductCreateUseCaseRequest(name string, price string) *ProductCreateUseCaseRequest {
	return &ProductCreateUseCaseRequest{
		Name:  name,
		Price: price,
	}
}

type ProductCreateUseCaseResponse struct {
	ProductDto      ProductDto
	ValidationError map[string][]string
}

func NewProductCreateUseCaseResponse(productDto ProductDto, validError map[string][]string) *ProductCreateUseCaseResponse {
	return &ProductCreateUseCaseResponse{
		ProductDto:      productDto,
		ValidationError: validError,
	}
}

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

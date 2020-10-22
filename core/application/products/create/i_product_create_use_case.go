package create

import (
	"github.com/mixmaru/my_contracts/core/application/products/dto"
)

type IProductCreateUseCase interface {
	Handle(request *ProductCreateUseCaseRequest) (*ProductCreateUseCaseResponse, error)
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
	ProductDto      dto.ProductDto
	ValidationError map[string][]string
}

func NewProductCreateUseCaseResponse(productDto dto.ProductDto, validError map[string][]string) *ProductCreateUseCaseResponse {
	return &ProductCreateUseCaseResponse{
		ProductDto:      productDto,
		ValidationError: validError,
	}
}

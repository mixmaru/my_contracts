package create

import (
	"github.com/mixmaru/my_contracts/core/application/products"
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
	ProductDto      products.ProductDto
	ValidationError map[string][]string
}

func NewProductCreateUseCaseResponse(productDto products.ProductDto, validError map[string][]string) *ProductCreateUseCaseResponse {
	return &ProductCreateUseCaseResponse{
		ProductDto:      productDto,
		ValidationError: validError,
	}
}

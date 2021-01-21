package get

import (
	"github.com/mixmaru/my_contracts/core/application/products"
)

type IProductGetUesCase interface {
	Handle(request *ProductGetUseCaseRequest) (*ProductGetUseCaseResponse, error)
}

type ProductGetUseCaseRequest struct {
	Id int
}

func NewProductGetUseCaseRequest(id int) *ProductGetUseCaseRequest {
	return &ProductGetUseCaseRequest{
		Id: id,
	}
}

type ProductGetUseCaseResponse struct {
	ProductDto products.ProductDto
}

func NewProductGetUseCaseResponse(productDto products.ProductDto) *ProductGetUseCaseResponse {
	return &ProductGetUseCaseResponse{
		ProductDto: productDto,
	}
}

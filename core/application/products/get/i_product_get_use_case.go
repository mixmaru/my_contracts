package get

import "github.com/mixmaru/my_contracts/core/application/products/dto"

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
	ProductDto dto.ProductDto
}

func NewProductGetUseCaseResponse(productDto dto.ProductDto) *ProductGetUseCaseResponse {
	return &ProductGetUseCaseResponse{
		ProductDto: productDto,
	}
}

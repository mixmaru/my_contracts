package get

import (
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
)

type ProductGetInteractor struct {
	productRepository products.IProductRepository
}

func NewProductGetInteractor(productRepository products.IProductRepository) *ProductGetInteractor {
	return &ProductGetInteractor{
		productRepository: productRepository,
	}
}

func (p *ProductGetInteractor) Handle(request *ProductGetUseCaseRequest) (*ProductGetUseCaseResponse, error) {
	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	// リポジトリつかってデータ取得
	entity, err := p.productRepository.GetById(request.Id, conn)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		// データがない
		return NewProductGetUseCaseResponse(products.ProductDto{}), nil
	}

	// dtoにつめる
	dto := products.NewProductDtoFromEntity(entity)

	// 返却
	return NewProductGetUseCaseResponse(dto), nil
}

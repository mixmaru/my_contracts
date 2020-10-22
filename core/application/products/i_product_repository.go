package products

import (
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"gopkg.in/gorp.v2"
)

type IProductRepository interface {
	Save(productEntity *product.ProductEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetById(id int, executor gorp.SqlExecutor) (*product.ProductEntity, error)
	GetByName(name string, executor gorp.SqlExecutor) ([]*product.ProductEntity, error)
	GetByRightToUseId(rightToUseId int, executor gorp.SqlExecutor) (*product.ProductEntity, error)
}

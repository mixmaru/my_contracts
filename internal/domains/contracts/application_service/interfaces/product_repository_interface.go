package interfaces

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"gopkg.in/gorp.v2"
)

type IProductRepository interface {
	Save(productEntity *entities.ProductEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetById(id int, executor gorp.SqlExecutor) (*entities.ProductEntity, error)
	GetByName(name string, executor gorp.SqlExecutor) (*entities.ProductEntity, error)
}

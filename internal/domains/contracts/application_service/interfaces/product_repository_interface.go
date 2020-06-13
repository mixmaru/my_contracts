package interfaces

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"gopkg.in/gorp.v2"
)

type IProductRepository interface {
	Save(productEntity *entities.ProductEntity, transaction *gorp.Transaction) (*entities.ProductEntity, error)
	GetById(id int, transaction *gorp.Transaction) (*entities.ProductEntity, error)
}

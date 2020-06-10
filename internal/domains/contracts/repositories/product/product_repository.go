package product

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"gopkg.in/gorp.v2"
)

type ProductRepository struct {
}

// 商品エンティティを保存する
func (r *ProductRepository) Save(productEntity *product.ProductEntiry, transaction *gorp.Transaction) (*produt.ProductEntirty, error) {

}

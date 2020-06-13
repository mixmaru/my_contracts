package product

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductRepository_Save(t *testing.T) {
	r := ProductRepository{}
	productEntity := entities.NewProductEntity("商品名", decimal.NewFromFloat(1000))
	_, err := r.Save(productEntity, nil)
	assert.NoError(t, err)
	//assert.NotEqual(t, 0, productEntity.Id())
}

package product

import (
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

// UserIndividualのインスタンス化をテスト
func TestProductEntity_New(t *testing.T) {
	// インスタンス化
	productEntity := New("name", decimal.NewFromFloat(1000))

	// テスト
	assert.Equal(t, "name", productEntity.Name())
	price := productEntity.Price()
	assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
}

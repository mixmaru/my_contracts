package entities

import (
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
)

type ProductEntity struct {
	name  string
	price decimal.Decimal
	BaseEntity
}

func NewProductEntity(name string, price decimal.Decimal) *ProductEntity {
	return &ProductEntity{
		name:  name,
		price: price,
	}
}

func (p *ProductEntity) Name() string {
	return p.name
}

func (p *ProductEntity) Price() decimal.Decimal {
	return p.price
}

package entities

import (
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"time"
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

func NewProductEntityWithData(id int, name string, price decimal.Decimal, createdAt, updatedAt time.Time) *ProductEntity {
	productEntity := ProductEntity{}
	productEntity.LoadData(id, name, price, createdAt, updatedAt)
	return &productEntity
}

func (p *ProductEntity) Name() string {
	return p.name
}

func (p *ProductEntity) Price() decimal.Decimal {
	return p.price
}

// 保持データをセットし直す
func (p *ProductEntity) LoadData(id int, name string, price decimal.Decimal, createdAt time.Time, updatedAt time.Time) {
	p.id = id
	p.name = name
	p.price = price
	p.createdAt = createdAt
	p.updatedAt = updatedAt
}

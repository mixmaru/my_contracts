package entities

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"time"
)

type ProductEntity struct {
	name  values.ProductNameValue
	price decimal.Decimal
	BaseEntity
}

func NewProductEntity(name string, price decimal.Decimal) (*ProductEntity, error) {
	nameValue, err := values.NewProductNameValue(name)
	if err != nil {
		return nil, err
	}
	return &ProductEntity{
		name:  nameValue,
		price: price,
	}, nil
}

func NewProductEntityWithData(id int, name string, price decimal.Decimal, createdAt, updatedAt time.Time) (*ProductEntity, error) {
	productEntity := ProductEntity{}
	err := productEntity.LoadData(id, name, price, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}
	return &productEntity, nil
}

func (p *ProductEntity) Name() string {
	return p.name.Value()
}

func (p *ProductEntity) Price() decimal.Decimal {
	return p.price
}

// 保持データをセットし直す
func (p *ProductEntity) LoadData(id int, name string, price decimal.Decimal, createdAt time.Time, updatedAt time.Time) error {
	nameValue, err := values.NewProductNameValue(name)
	if err != nil {
		return err
	}

	p.id = id
	p.name = nameValue
	p.price = price
	p.createdAt = createdAt
	p.updatedAt = updatedAt
	return nil
}

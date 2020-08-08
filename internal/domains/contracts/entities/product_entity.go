package entities

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/pkg/errors"
	"time"
)

type ProductEntity struct {
	name            values.ProductNameValue
	priceYearly     *ProductPriceYearlyEntity     // 年間契約の場合の情報を保持させる
	priceMonthly    *ProductPriceMonthlyEntity    // 月契約の場合の情報を保持させる
	priceLump       *ProductPriceLumpEntity       // 一括購入の場合の情報を保持させる
	priceCustomTerm *ProductPriceCustomTermEntity // 任意期間契約の倍の情報を保持させる
	BaseEntity
}

type ProductPriceYearlyEntity struct {
	price values.ProductPriceValue
}

type ProductPriceMonthlyEntity struct {
	price values.ProductPriceValue
}

type ProductPriceLumpEntity struct {
	price values.ProductPriceValue
}

type ProductPriceCustomTermEntity struct {
	price values.ProductPriceValue
	term  int // 契約更新サイクル日数
}

func NewProductEntity(name string, price string) (*ProductEntity, error) {
	nameValue, err := values.NewProductNameValue(name)
	if err != nil {
		return nil, err
	}
	priceValue, err := values.NewProductPriceValue(price)
	if err != nil {
		return nil, err
	}

	return &ProductEntity{
		name:         nameValue,
		priceMonthly: &ProductPriceMonthlyEntity{price: priceValue},
	}, nil
}

func NewProductEntityWithData(id int, name string, price string, createdAt, updatedAt time.Time) (*ProductEntity, error) {
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

func (p *ProductEntity) MonthlyPrice() (decimal.Decimal, error) {
	if p.priceMonthly != nil {
		return p.priceMonthly.price.Value(), nil
	} else {
		return decimal.Decimal{}, errors.Errorf("月契約価格が存在しない。productId: %v", p.Id())
	}
}

// 保持データをセットし直す
func (p *ProductEntity) LoadData(id int, name string, price string, createdAt time.Time, updatedAt time.Time) error {
	nameValue, err := values.NewProductNameValue(name)
	if err != nil {
		return err
	}
	priceValue, err := values.NewProductPriceValue(price)
	if err != nil {
		return err
	}
	if p.priceMonthly == nil {
		p.priceMonthly = &ProductPriceMonthlyEntity{}
	}

	p.id = id
	p.name = nameValue
	p.priceMonthly.price = priceValue
	p.createdAt = createdAt
	p.updatedAt = updatedAt
	return nil
}

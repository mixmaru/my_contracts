package product

import (
	"time"

	"github.com/mixmaru/my_contracts/lib/decimal"
	"github.com/pkg/errors"
)

type ProductEntity struct {
	id              int
	name            ProductNameValue
	priceYearly     *ProductPriceYearlyEntity     // 年間契約の場合の情報を保持させる
	priceMonthly    *ProductPriceMonthlyEntity    // 月契約の場合の情報を保持させる
	priceLump       *ProductPriceLumpEntity       // 一括購入の場合の情報を保持させる
	priceCustomTerm *ProductPriceCustomTermEntity // 任意期間契約の倍の情報を保持させる
	createdAt       time.Time
	updatedAt       time.Time
}

func NewProductEntity(name string, price string) (*ProductEntity, error) {
	nameValue, err := NewProductNameValue(name)
	if err != nil {
		return nil, err
	}
	priceValue, err := NewProductPriceValue(price)
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

func (p *ProductEntity) Id() int {
	return p.id
}

func (p *ProductEntity) CreatedAt() time.Time {
	return p.createdAt
}

func (p *ProductEntity) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *ProductEntity) Name() string {
	return p.name.Value()
}

/*
月契約があれば月額金額を返す。
なければboolでfalseが返る
*/
func (p *ProductEntity) MonthlyPrice() (decimal.Decimal, bool) {
	if p.priceMonthly != nil {
		return p.priceMonthly.price.Value(), true
	} else {
		return decimal.Decimal{}, false
	}
}

// 保持データをセットし直す
func (p *ProductEntity) LoadData(id int, name string, price string, createdAt time.Time, updatedAt time.Time) error {
	nameValue, err := NewProductNameValue(name)
	if err != nil {
		return err
	}
	priceValue, err := NewProductPriceValue(price)
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

const (
	TermMonthly string = "monthly"
	TermYearly  string = "yearly"
	TermCustom  string = "custom"
	TermLump    string = "lump" // 一括購入
)

func (p *ProductEntity) GetTermType() (termType string, err error) {
	if p.priceMonthly != nil {
		return TermMonthly, nil
	} else if p.priceYearly != nil {
		return TermYearly, nil
	} else if p.priceCustomTerm != nil {
		return TermCustom, nil
	} else if p.priceLump != nil {
		return TermLump, nil
	}
	return "", errors.Errorf("考慮外。productEntity: %+v", p)
}

/*
カスタム期間商品の場合、カスタム期間を返す。
カスタム期間商品ではない場合、エラーを返す
*/
func (p *ProductEntity) GetCustomTerm() (int, error) {
	if p.priceCustomTerm == nil {
		return 0, errors.Errorf("カスタム機関商品ではありません。productEntity: %+v", p)
	}

	return p.priceCustomTerm.term, nil
}

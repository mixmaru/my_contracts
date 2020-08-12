package domain_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"time"
)

type BillingCalculatorDomainService struct {
	ProductRepository interfaces.IProductRepository
}

func NewBillingCalculatorDomainService(productRepository interfaces.IProductRepository) *BillingCalculatorDomainService {
	return &BillingCalculatorDomainService{ProductRepository: productRepository}
}

func (b *BillingCalculatorDomainService) BillingAmount(contract *entities.ContractEntity, targetDate time.Time, executor gorp.SqlExecutor) (decimal.Decimal, error) {
	// 商品データを取得する
	product, err := b.ProductRepository.GetById(contract.ProductId(), executor)
	if err != nil {
		return decimal.Decimal{}, err
	}
	// 計算して返す
	price, exist := product.MonthlyPrice()
	if !exist {
		return decimal.Decimal{}, errors.Errorf("月額料金が取得できなかった。product: %+v", product)
	}
	return price, nil
}

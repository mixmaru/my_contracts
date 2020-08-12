package domain_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"math"
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
	// 月額を取得
	price, exist := product.MonthlyPrice()
	if !exist {
		return decimal.Decimal{}, errors.Errorf("月額料金が取得できなかった。product: %+v", product)
	}
	billingStartDate := contract.BillingStartDate()
	// 課金開始日から次の課金開始日の間の日数を取得する
	thisTermDateNum := b.thisTermDateNum(billingStartDate)
	// 課金開始日からtargetDateの間の日数を計算する
	DateNumForTargetDate := int(math.Ceil(targetDate.Sub(billingStartDate).Hours() / 24))
	// 金額を計算する（月額 / 月間日数 * 対象日までの日数）
	div := price.Div(decimal.NewFromInt(int64(thisTermDateNum)))
	retPrice := div.Mul(decimal.NewFromInt(int64(DateNumForTargetDate)))
	return retPrice, nil
}

func (b *BillingCalculatorDomainService) thisTermDateNum(lastBillingStartDate time.Time) int {
	duration, _ := time.ParseDuration("744h")
	return int(math.Ceil(duration.Hours() / 24))
}

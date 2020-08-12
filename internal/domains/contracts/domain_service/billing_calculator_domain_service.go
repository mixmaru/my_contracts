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

	// 規定日数に足りない場合のみ日割り計算をする（満額時も日割りロジックで算出しようとすると小数点以下の都合で定価がおかしくなる）
	if DateNumForTargetDate < thisTermDateNum {
		// 日割り金額を計算する（月額 / 月間日数 * 対象日までの日数）
		div := price.Div(decimal.NewFromInt(int64(thisTermDateNum)))
		retPrice := div.Mul(decimal.NewFromInt(int64(DateNumForTargetDate)))
		return retPrice.Truncate(0), nil
	} else {
		// 満額
		return price, nil
	}
}

func (b *BillingCalculatorDomainService) thisTermDateNum(lastBillingStartDate time.Time) int {
	// 次の課金開始日を算出
	nextBillingStartDate := lastBillingStartDate.AddDate(0, 1, 0)
	// 日数算出
	oneMonthDuration := nextBillingStartDate.Sub(lastBillingStartDate)
	return int(math.Ceil(oneMonthDuration.Hours() / 24))
}

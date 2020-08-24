package domain_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"math"
	"time"
)

type BillingCalculatorDomainService struct {
	productRepository    interfaces.IProductRepository
	contractRepository   interfaces.IContractRepository
	rightToUseRepository interfaces.IRightToUseRepository
}

func NewBillingCalculatorDomainService(productRepository interfaces.IProductRepository, contractRepository interfaces.IContractRepository, rightToUseRepository interfaces.IRightToUseRepository) *BillingCalculatorDomainService {
	return &BillingCalculatorDomainService{
		productRepository:    productRepository,
		contractRepository:   contractRepository,
		rightToUseRepository: rightToUseRepository,
	}
}

func (b *BillingCalculatorDomainService) BillingAmount(rightToUseId int, executor gorp.SqlExecutor) (decimal.Decimal, error) {
	////// 満額請求金額を取得する
	// 使用権データ取得
	rightToUse, err := b.rightToUseRepository.GetById(rightToUseId, executor)
	if err != nil {
		return decimal.Decimal{}, errors.WithMessagef(err, "使用権データの取得失敗。rightToUseId: %v", rightToUseId)
	}
	if rightToUse == nil {
		return decimal.Decimal{}, errors.Errorf("使用権データが存在しない。rightToUseId: %v", rightToUseId)
	}
	// 商品データ、契約データ取得
	contract, product, _, err := b.contractRepository.GetById(rightToUse.ContractId(), executor)
	if err != nil {
		return decimal.Decimal{}, errors.WithMessagef(err, "商品データ、契約データの取得失敗。rightToUse: %v", rightToUse)
	}

	////// 満額期間日数を取得計算する（今は月払い固定）
	// 使用権の開始日から1ヶ月後の同日までに存在する日数を算出
	subDuration := rightToUse.ValidFrom().AddDate(0, 1, 0).Sub(rightToUse.ValidFrom())
	fullBillingDateNum := subDuration.Hours() / 24

	////// 使用権の課金対象期間を算出する
	// 課金開始日を決定
	var billingStartDate time.Time
	if contract.BillingStartDate().After(rightToUse.ValidFrom()) {
		billingStartDate = contract.BillingStartDate()
	} else {
		billingStartDate = rightToUse.ValidFrom()
	}

	// 課金対象日数を算出。時間を24で割って、余ったら切り上げ
	billHours := rightToUse.ValidTo().Sub(billingStartDate).Hours()
	billDate := math.Ceil(billHours / float64(24))

	// 定価
	price, ok := product.MonthlyPrice()
	if !ok {
		return decimal.Decimal{}, errors.Errorf("月額金額が設定されてない。product: %+v", product)
	}

	////// 満額期間でないなら日割り計算する
	if billDate == fullBillingDateNum {
		// 満了なので定価
		return price, nil
	} else {
		// 日割り （1日あたり金額 * 課金日数）
		fullBillingDateNumDecimal := decimal.NewFromFloat(fullBillingDateNum) // decimal化
		billDateDecimal := decimal.NewFromFloat(billDate)                     // decimal化
		byDayPrice := price.Div(fullBillingDateNumDecimal)
		dailyRatePrice := byDayPrice.Mul(billDateDecimal)
		return dailyRatePrice.Truncate(0), nil
	}
	//// 商品データを取得する
	//product, err := b.productRepository.GetById(contract.ProductId(), executor)
	//if err != nil {
	//	return decimal.Decimal{}, err
	//}
	//// 月額を取得
	//price, exist := product.MonthlyPrice()
	//if !exist {
	//	return decimal.Decimal{}, errors.Errorf("月額料金が取得できなかった。product: %+v", product)
	//}
	//billingStartDate := contract.LastBillingStartDate(targetDate)
	//// 課金開始日から次の課金開始日の間の日数を取得する
	//thisTermFullDateNum := b.billingTermFullDateNum(billingStartDate)
	//// 課金開始日からtargetDateの間の日数を計算する
	//DateNumForTargetDate := int(math.Ceil(targetDate.Sub(billingStartDate).Hours() / 24))
	//
	//return b.prorate(price, thisTermFullDateNum, DateNumForTargetDate), nil
}

/*
（1ヶ月契約の）与えられた日を課金開始日として、次の課金開始日までの日数を返す

	billingStartDate	課金開始日
*/
func (b *BillingCalculatorDomainService) billingTermFullDateNum(billingStartDate time.Time) int {
	// 次の課金開始日を算出
	nextBillingStartDate := billingStartDate.AddDate(0, 1, 0)
	// 日数算出
	oneMonthDuration := nextBillingStartDate.Sub(billingStartDate)
	return int(math.Ceil(oneMonthDuration.Hours() / 24))
}

/*
日割り計算

	basePrice 	定価
	fullDateNum 基底日数
	usedDateNum 使用日
*/
func (b *BillingCalculatorDomainService) prorate(basePrice decimal.Decimal, fullDateNum int, usedDateNum int) decimal.Decimal {
	// 規定日数に足りない場合のみ日割り計算をする（満額時も日割りロジックで算出しようとすると小数点以下の都合で定価がおかしくなる）
	if usedDateNum >= fullDateNum {
		// 満額
		return basePrice
	}

	// 日割り金額を計算する（月額 / 月間日数 * 対象日までの日数）
	div := basePrice.Div(decimal.NewFromInt(int64(fullDateNum)))
	retPrice := div.Mul(decimal.NewFromInt(int64(usedDateNum)))
	return retPrice.Truncate(0)
}

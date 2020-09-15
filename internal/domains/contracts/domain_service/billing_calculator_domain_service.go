package domain_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"math"
	"time"
)

type BillingCalculatorDomainService struct {
	productRepository    interfaces.IProductRepository
	contractRepository   interfaces.IContractRepository
	rightToUseRepository interfaces.IRightToUseRepository
	billRepository       interfaces.IBillRepository
}

func NewBillingCalculatorDomainService(productRepository interfaces.IProductRepository, contractRepository interfaces.IContractRepository, rightToUseRepository interfaces.IRightToUseRepository, billRepository interfaces.IBillRepository) *BillingCalculatorDomainService {
	return &BillingCalculatorDomainService{
		productRepository:    productRepository,
		contractRepository:   contractRepository,
		rightToUseRepository: rightToUseRepository,
		billRepository:       billRepository,
	}
}

// 渡した指定日を実行日として請求の実行をする
//
// ※もしデータ量がもの凄く多かったら、長期間ロックがかかるかも。それであれば、1件1件取得 => commitを長時間続けたほうがいいのかもしれない。
// その場合、長時間処理時間がかかるので、別スレッドとかで非同期に動かすことを検討する
func (b *BillingCalculatorDomainService) ExecuteBilling(executeDate time.Time, executor gorp.SqlExecutor) error {
	// 対象使用権を取得する
	rightToUses, err := b.rightToUseRepository.GetBillingTargetByBillingDate(executeDate, executor)
	if err != nil {
		return err
	}

	billAggs, err := b.createBillAggligationsFromRightToUseEntities(executeDate, rightToUses, executor)

	for _, billAgg := range billAggs {
		_, err := b.billRepository.Create(billAgg, executor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BillingCalculatorDomainService) createBillAggligationsFromRightToUseEntities(executeDate time.Time, rightToUses []*entities.RightToUseEntity, executor gorp.SqlExecutor) ([]*entities.BillAggregation, error) {
	var retBillAggs []*entities.BillAggregation

	prevUserId := 0
	var billAgg *entities.BillAggregation
	orderNum := 0
	for _, rightToUse := range rightToUses {

		// userIdを取得
		contract, _, _, err := b.contractRepository.GetById(rightToUse.ContractId(), executor)
		if err != nil {
			return nil, err
		}
		userId := contract.UserId()

		// userId毎にbillAggを作成する
		if userId != prevUserId {
			orderNum = 0
			prevUserId = userId
			billAgg = entities.NewBillingAggregation(executeDate, userId)
			retBillAggs = append(retBillAggs, billAgg)
		}
		orderNum += 1
		amount, err := b.BillingAmount(rightToUse.Id(), executor)
		if err != nil {
			return nil, err
		}
		err = billAgg.AddBillDetail(entities.NewBillingDetailEntity(orderNum, rightToUse.Id(), amount))
		if err != nil {
			return nil, err
		}
	}

	return retBillAggs, nil
}

func (b *BillingCalculatorDomainService) BillingAmount(rightToUseId int, executor gorp.SqlExecutor) (decimal.Decimal, error) {
	////// 満額請求金額を取得する
	// 必要データ取得
	rightToUse, contract, product, _, err := b.getEntitiesByRightToUseId(rightToUseId, executor)
	if err != nil {
		return decimal.Decimal{}, err
	}

	////// 満額期間日数と理療機関を取得計算する（今は月払い固定）
	usageDateNum, fullBillingDateNum := b.getUsageDate(rightToUse.ValidFrom(), rightToUse.ValidTo(), contract.BillingStartDate())

	// 定価取得
	price, ok := product.MonthlyPrice()
	if !ok {
		return decimal.Decimal{}, errors.Errorf("月額金額が設定されてない。product: %+v", product)
	}

	////// 満額期間でないなら日割り計算する
	if usageDateNum == fullBillingDateNum {
		// 満了なので定価
		return price, nil
	} else {
		// 日割り
		return b.prorate(price, fullBillingDateNum, usageDateNum), nil
	}
}

// 請求金額計算に必要なデータを取得する
func (b *BillingCalculatorDomainService) getEntitiesByRightToUseId(rightToUseId int, executor gorp.SqlExecutor) (
	*entities.RightToUseEntity,
	*entities.ContractEntity,
	*entities.ProductEntity,
	interface{},
	error,
) {
	// 使用権データ取得
	rightToUse, err := b.rightToUseRepository.GetById(rightToUseId, executor)
	if err != nil {
		return nil, nil, nil, nil, errors.WithMessagef(err, "使用権データの取得失敗。rightToUseId: %v", rightToUseId)
	}
	if rightToUse == nil {
		return nil, nil, nil, nil, errors.Errorf("使用権データが存在しない。rightToUseId: %v", rightToUseId)
	}
	// 商品データ、契約データ取得
	contract, product, user, err := b.contractRepository.GetById(rightToUse.ContractId(), executor)
	if err != nil {
		return nil, nil, nil, nil, errors.WithMessagef(err, "商品データ、契約データの取得失敗。rightToUse: %v", rightToUse)
	}
	return rightToUse, contract, product, user, nil
}

/*
利用権の開始・終了日と課金開始日から、課金対象日数と満了日数を返す

args
	validFrom			利用開始日
	validTo				利用終了日
	billingStartDate	課金開始日
*/
func (b *BillingCalculatorDomainService) getUsageDate(validFrom, validTo, billingStartDate time.Time) (usageDateNum, fullBillingDateNum int) {
	// jstで扱うようにする（そうしないと日割り計算の基準日がおかしくなる）
	jst := utils.CreateJstLocation()
	from := validFrom.In(jst)
	to := validTo.In(jst)
	// 使用権の開始日から1ヶ月後の同日までに存在する日数を算出
	subDuration := from.AddDate(0, 1, 0).Sub(from)
	fullBillingDateNum = int(subDuration.Hours() / 24)

	////// 使用権の課金対象期間を算出する
	// 課金開始日を決定
	var realBillingStartDate time.Time
	if billingStartDate.After(from) {
		realBillingStartDate = billingStartDate
	} else {
		realBillingStartDate = from
	}

	// 課金対象日数を算出。時間を24で割って、余ったら切り上げ
	billHours := to.Sub(realBillingStartDate).Hours()
	usageDateNum = int(math.Ceil(billHours / float64(24)))

	return usageDateNum, fullBillingDateNum
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

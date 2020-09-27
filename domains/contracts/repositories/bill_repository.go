package repositories

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/data_mappers"
	"github.com/mixmaru/my_contracts/lib/decimal"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"time"
)

type BillRepository struct {
	*BaseRepository
}

func NewBillRepository() *BillRepository {
	return &BillRepository{
		&BaseRepository{},
	}
}

func (b *BillRepository) Create(billAggregation *entities.BillAggregation, executor gorp.SqlExecutor) (savedId int, err error) {
	// マッパーを用意
	billMap := data_mappers.BillMapper{}
	billMap.BillingDate = billAggregation.BillingDate()
	billMap.UserId = billAggregation.UserId()

	// 保存実行
	err = executor.Insert(&billMap)
	if err != nil {
		return 0, errors.Wrapf(err, "billレコードの保存に失敗しました。billMap: %+v", billMap)
	}

	// マッパー用意
	billDetails := billAggregation.BillDetails()
	billDetailMaps := make([]interface{}, 0, len(billDetails))
	orderNum := 0
	for _, detail := range billDetails {
		orderNum += 1
		detailMap := data_mappers.BillDetailMapper{}
		detailMap.BillingAmount = detail.BillingAmount()
		detailMap.RightToUseId = detail.RightToUseId()
		detailMap.OrderNum = orderNum
		detailMap.BillId = billMap.Id
		billDetailMaps = append(billDetailMaps, &detailMap)
	}

	// 保存実行
	err = executor.Insert(billDetailMaps...)
	if err != nil {
		return 0, errors.Wrapf(err, "bill_detailsレコードの保存に失敗しました。billDetailMaps: %+v", billDetailMaps)
	}
	return billMap.Id, nil
}

func (b *BillRepository) GetById(id int, executor gorp.SqlExecutor) (aggregation *entities.BillAggregation, err error) {
	// クエリ作成
	query := createGetByIdQuery()
	// マッパー用意
	var mappers []*BillAndBillDetailsMapper

	// データ取得
	_, err = executor.Select(&mappers, query, id)
	if err != nil {
		return nil, errors.Wrapf(err, "billデータ取得失敗。id: %v", id)
	}

	// 返却データに組み立てる
	detailEntities := make([]*entities.BillDetailEntity, 0, len(mappers))
	for _, detail := range mappers {
		entity := entities.NewBillingDetailsEntityWithData(
			detail.DetailId,
			detail.DetailRightToUseId,
			detail.DetailBillingAmount,
			detail.DetailCreatedAt,
			detail.DetailUpdatedAt,
		)
		detailEntities = append(detailEntities, entity)
	}
	billAgg := entities.NewBillingAggregationWithData(
		mappers[0].Id,
		mappers[0].BillingDate,
		mappers[0].UserId,
		mappers[0].PaymentConfirmedAt,
		detailEntities,
		mappers[0].CreatedAt,
		mappers[0].UpdatedAt,
	)

	// 返却
	return billAgg, nil
}

func (b *BillRepository) GetByUserId(userId int, executor gorp.SqlExecutor) (aggregation []*entities.BillAggregation, err error) {
	query := createGetByUserIdQuery()
	// マッパー用意
	var mappers []*BillAndBillDetailsMapper

	// データ取得
	_, err = executor.Select(&mappers, query, userId)
	if err != nil {
		return nil, errors.Wrapf(err, "billデータ取得失敗。query: %v ,userId: %v", query, userId)
	}

	// 返却データに組み立てる
	retBillAggs, err := createBillAggsFromMappers(mappers)
	if err != nil {
		return nil, err
	}

	return retBillAggs, nil
}

type BillAndBillDetailsMapper struct {
	data_mappers.BillMapper
	DetailId            int             `db:"detail_id"`
	DetailOrderNum      int             `db:"detail_order_num"`
	DetailRightToUseId  int             `db:"detail_right_to_use_id"`
	DetailBillingAmount decimal.Decimal `db:"detail_billing_amount"`
	DetailCreatedAt     time.Time       `db:"detail_created_at"`
	DetailUpdatedAt     time.Time       `db:"detail_updated_at"`
}

func createGetByIdQuery() string {
	baseQuery := createGetBaseQuery()
	query := fmt.Sprintf(baseQuery, "bills.id = $1")
	return query
}

func createGetByUserIdQuery() string {
	baseQuery := createGetBaseQuery()
	query := fmt.Sprintf(baseQuery, "bills.user_id = $1")
	return query
}

func createGetBaseQuery() string {
	baseQuery := `
SELECT
       bills.id AS id,
       bills.billing_date AS billing_date,
       bills.user_id AS user_id,
       bills.payment_confirmed_at AS payment_confirmed_at,
       bills.created_at AS created_at,
       bills.updated_at AS updated_at,
       bd.id AS detail_id,
       bd.order_num AS detail_order_num,
       bd.right_to_use_id AS detail_right_to_use_id,
       bd.billing_amount AS detail_billing_amount,
       bd.created_at AS detail_created_at,
       bd.updated_at AS detail_updated_at
FROM bills
INNER JOIN bill_details bd on bills.id = bd.bill_id
WHERE %v
ORDER BY bd.order_num
;
`
	return baseQuery
}

// 返却データに組み立てる
func createBillAggsFromMappers(mappers []*BillAndBillDetailsMapper) ([]*entities.BillAggregation, error) {
	var retBillAggs []*entities.BillAggregation
	prevId := 0
	var billAgg *entities.BillAggregation
	for _, record := range mappers {
		if record.Id != prevId {
			prevId = record.Id
			// 前回ループで作ったbillAggがあればretBillAggsに追加する
			if billAgg != nil {
				retBillAggs = append(retBillAggs, billAgg)
			}
			// 新しいbillAggを作成する
			billAgg = entities.NewBillingAggregationWithData(
				record.Id,
				record.BillingDate,
				record.UserId,
				record.PaymentConfirmedAt,
				[]*entities.BillDetailEntity{},
				record.CreatedAt,
				record.UpdatedAt,
			)
		}
		// detailsを作ってbillAggに追加する
		detail := entities.NewBillingDetailsEntityWithData(
			record.DetailId,
			record.DetailRightToUseId,
			record.DetailBillingAmount,
			record.DetailCreatedAt,
			record.DetailUpdatedAt,
		)
		err := billAgg.AddBillDetail(detail)
		if err != nil {
			return nil, err
		}
	}
	// 前回ループで作ったbillAggがあればretBillAggsに追加する
	if billAgg != nil {
		retBillAggs = append(retBillAggs, billAgg)
	}

	return retBillAggs, nil
}

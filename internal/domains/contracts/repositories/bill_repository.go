package repositories

import (
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
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
	for _, detail := range billDetails {
		detailMap := data_mappers.BillDetailMapper{}
		detailMap.BillingAmount = detail.BillingAmount()
		detailMap.RightToUseId = detail.RightToUseId()
		detailMap.OrderNum = detail.OrderNum()
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
	query := `
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
WHERE bills.id = $1
ORDER BY bd.order_num
;
`
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
			detail.DetailOrderNum,
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

type BillAndBillDetailsMapper struct {
	data_mappers.BillMapper
	DetailId            int             `db:"detail_id"`
	DetailOrderNum      int             `db:"detail_order_num"`
	DetailRightToUseId  int             `db:"detail_right_to_use_id"`
	DetailBillingAmount decimal.Decimal `db:"detail_billing_amount"`
	DetailCreatedAt     time.Time       `db:"detail_created_at"`
	DetailUpdatedAt     time.Time       `db:"detail_updated_at"`
}

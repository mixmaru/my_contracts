package db

import (
	"database/sql"
	"fmt"
	"github.com/mixmaru/my_contracts/core/domain/models/bill"
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

func (b *BillRepository) Create(billEntity *bill.BillEntity, executor gorp.SqlExecutor) (savedId int, err error) {
	// マッパーを用意
	billMap := BillMapper{}
	billMap.BillingDate = billEntity.BillingDate()
	billMap.UserId = billEntity.UserId()

	// 保存実行
	err = executor.Insert(&billMap)
	if err != nil {
		return 0, errors.Wrapf(err, "billレコードの保存に失敗しました。billMap: %+v", billMap)
	}

	// マッパー用意
	billDetails := billEntity.BillDetails()
	billDetailMaps := make([]interface{}, 0, len(billDetails))
	orderNum := 0
	for _, detail := range billDetails {
		orderNum += 1
		detailMap := BillDetailMapper{}
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

func (b *BillRepository) GetById(id int, executor gorp.SqlExecutor) (billEntity *bill.BillEntity, err error) {
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
	detailEntities := make([]*bill.BillDetailEntity, 0, len(mappers))
	for _, detail := range mappers {
		entity := bill.NewBillDetailsEntityWithData(
			detail.DetailId,
			detail.DetailRightToUseId,
			detail.DetailBillingAmount,
			detail.DetailCreatedAt,
			detail.DetailUpdatedAt,
		)
		detailEntities = append(detailEntities, entity)
	}
	billAgg := bill.NewBillEntityWithData(
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

func (b *BillRepository) GetByUserId(userId int, executor gorp.SqlExecutor) (billEntities []*bill.BillEntity, err error) {
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
	BillMapper
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
func createBillAggsFromMappers(mappers []*BillAndBillDetailsMapper) ([]*bill.BillEntity, error) {
	var retBill []*bill.BillEntity
	prevId := 0
	var billEntity *bill.BillEntity
	for _, record := range mappers {
		if record.Id != prevId {
			prevId = record.Id
			// 前回ループで作ったbillAggがあればretBillAggsに追加する
			if billEntity != nil {
				retBill = append(retBill, billEntity)
			}
			// 新しいbillAggを作成する
			billEntity = bill.NewBillEntityWithData(
				record.Id,
				record.BillingDate,
				record.UserId,
				record.PaymentConfirmedAt,
				[]*bill.BillDetailEntity{},
				record.CreatedAt,
				record.UpdatedAt,
			)
		}
		// detailsを作ってbillAggに追加する
		detail := bill.NewBillDetailsEntityWithData(
			record.DetailId,
			record.DetailRightToUseId,
			record.DetailBillingAmount,
			record.DetailCreatedAt,
			record.DetailUpdatedAt,
		)
		err := billEntity.AddBillDetail(detail)
		if err != nil {
			return nil, err
		}
	}
	// 前回ループで作ったbillAggがあればretBillAggsに追加する
	if billEntity != nil {
		retBill = append(retBill, billEntity)
	}

	return retBill, nil
}

type BillMapper struct {
	Id                 int          `db:"id"`
	BillingDate        time.Time    `db:"billing_date"`
	UserId             int          `db:"user_id"`
	PaymentConfirmedAt sql.NullTime `db:"payment_confirmed_at"`
	CreatedAtUpdatedAtMapper
}

type BillDetailMapper struct {
	Id            int             `db:"id"`
	BillId        int             `db:"bill_id"`
	OrderNum      int             `db:"order_num"`
	RightToUseId  int             `db:"right_to_use_id"`
	BillingAmount decimal.Decimal `db:"billing_amount"`
	CreatedAtUpdatedAtMapper
}

func (b *BillDetailMapper) SetDataToEntity(entity interface{}) error {
	value, ok := entity.(*bill.BillDetailEntity)
	if !ok {
		return errors.Errorf("*entities.BillDetailEntityではないものが渡された。entity: %T", entity)
	}
	value.LoadData(
		b.Id,
		b.RightToUseId,
		b.BillingAmount,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return nil
}

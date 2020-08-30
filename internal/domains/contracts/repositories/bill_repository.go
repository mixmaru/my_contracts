package repositories

import (
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
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

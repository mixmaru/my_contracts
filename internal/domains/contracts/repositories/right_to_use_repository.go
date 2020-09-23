package repositories

import (
	"database/sql"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"time"
)

type RightToUseRepository struct {
	*BaseRepository
}

func NewRightToUseRepository() *RightToUseRepository {
	return &RightToUseRepository{
		&BaseRepository{},
	}
}

func (r *RightToUseRepository) Create(rightToUseEntity *entities.RightToUseEntity, executor gorp.SqlExecutor) (savedId int, err error) {
	// データマッパーを用意する
	mapper := data_mappers.NewRightToUseMapperFromEntity(rightToUseEntity)

	// データ登録実行する
	err = executor.Insert(&mapper)
	if err != nil {
		return 0, errors.Wrapf(err, "データ登録失敗。mapper: %+v", mapper)
	}

	return mapper.Id, nil
}

func (r *RightToUseRepository) GetById(id int, executor gorp.SqlExecutor) (*entities.RightToUseEntity, error) {
	// データマッパー用意
	mapper := &data_mappers.RightToUseMapper{}
	// query用意
	query := `
SELECT 
	id,
	contract_id,
	valid_from,
	valid_to,
	created_at,
	updated_at
FROM right_to_use
WHERE id = $1;
`
	// 取得実行
	err := executor.SelectOne(mapper, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, errors.Wrapf(err, "使用権データの取得失敗。id: %v, query: %v", id, query)
		}
	}

	// entityに詰める
	entity := entities.NewRightToUseEntityWithData(mapper.Id, mapper.ContractId, mapper.ValidFrom, mapper.ValidTo, mapper.CreatedAt, mapper.UpdatedAt)
	// 返却
	return entity, nil
}

/*
渡した日（請求実行日）以前の請求すべき請求（まだ請求実行されていない）がある使用権データをすべて返す

例）請求実行日が6/1の場合
契約の課金開始日が6/1の使用権（期間：6/1 ~ 6/30）=> 請求すべき。
契約の課金開始日が6/2の使用権（期間：6/1 ~ 6/30）=> 請求すべきでない。
契約の課金開始日が6/2の使用権（期間：7/1 ~ 7/31）=> 請求すべきでない。
*/
func (r *RightToUseRepository) GetBillingTargetByBillingDate(billingDate time.Time, executor gorp.SqlExecutor) ([]*entities.RightToUseEntity, error) {
	query := `
SELECT 
	rtu.id,
	rtu.contract_id,
	rtu.valid_from,
	rtu.valid_to,
	rtu.created_at,
	rtu.updated_at
FROM right_to_use rtu
LEFT OUTER JOIN bill_details bd ON rtu.id = bd.right_to_use_id
INNER JOIN contracts c ON c.id = rtu.contract_id
WHERE bd.id IS NULL
AND valid_from <= $1
AND c.billing_start_date <= $1
ORDER BY c.user_id, rtu.id
;
`
	var mappers []*data_mappers.RightToUseMapper
	var _, err = executor.Select(&mappers, query, billingDate)
	if err != nil {
		return nil, errors.Wrapf(err, "請求対象使用権の取得に失敗しました。query: %v, billingDate: %v", query, billingDate)
	}

	// mapperからentityを作る
	retEntities := make([]*entities.RightToUseEntity, 0, len(mappers))
	for _, mapper := range mappers {
		entity := entities.NewRightToUseEntityWithData(mapper.Id, mapper.ContractId, mapper.ValidFrom, mapper.ValidTo, mapper.CreatedAt, mapper.UpdatedAt)
		retEntities = append(retEntities, entity)
	}

	return retEntities, nil
}

/*
渡した日（実行日）から5日以内に終了し、かつ、まだ次の期間の使用権データが存在しない使用権を全て返す

例）実行日が6/1の場合
使用権の終了日が6/1の使用権=> 返る
使用権の終了日が6/6の使用権=> 返る
使用権の終了日が6/7の使用権=> 返らない
使用権の終了日が6/1だが、次（6/2 ~ 7/1の期間）の使用権が存在する=> 返らない
*/
func (r *RightToUseRepository) GetRecurTargets(executeDate time.Time, executor gorp.SqlExecutor) ([]*entities.RightToUseEntity, error) {
	from := executeDate
	to := executeDate.AddDate(0, 0, 5)

	query := `
WITH tmp_t AS (
    SELECT *, row_number() over (partition by contract_id order by valid_to DESC) AS num FROM right_to_use
)
SELECT 
	id,
	contract_id,
	valid_from,
	valid_to,
	created_at,
	updated_at
FROM tmp_t
WHERE num = 1
AND $1 <= tmp_t.valid_to
AND tmp_t.valid_to < $2
;`

	var mappers []*data_mappers.RightToUseMapper
	var _, err = executor.Select(&mappers, query, from, to)
	if err != nil {
		return nil, errors.Wrapf(err, "継続処理対象使用権の取得に失敗しました。query: %v, from: %v, to: %v", query, from, to)
	}

	// mapperからentityを作る
	retEntities := make([]*entities.RightToUseEntity, 0, len(mappers))
	for _, mapper := range mappers {
		entity := entities.NewRightToUseEntityWithData(mapper.Id, mapper.ContractId, mapper.ValidFrom, mapper.ValidTo, mapper.CreatedAt, mapper.UpdatedAt)
		retEntities = append(retEntities, entity)
	}

	return retEntities, nil
}

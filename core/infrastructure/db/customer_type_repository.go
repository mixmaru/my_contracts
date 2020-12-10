package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type CustomerTypeRepository struct{}

func NewCustomerTypeRepository() *CustomerTypeRepository {
	return &CustomerTypeRepository{}
}

// カスタマータイプを新規保存する
func (r *CustomerTypeRepository) Create(customerTypeEntity *customer.CustomerTypeEntity, executor gorp.SqlExecutor) (savedId int, err error) {
	////// customer_typeの保存
	// mappperに詰める
	customerTypeMapper := CustomerTypeMapper{
		Name: customerTypeEntity.Name(),
	}
	// 保存実行
	if err := executor.Insert(&customerTypeMapper); err != nil {
		return 0, errors.Wrapf(err, "customer_typeテーブルへの保存に失敗しました。%v", customerTypeEntity)
	}

	////// costomer_propertiesの保存
	// mappperに詰める
	customerPropertyTypeEntities := customerTypeEntity.CustomerPropertyTypes()
	customerPropertyMappers := make([]interface{}, 0, len(customerPropertyTypeEntities))
	for _, entity := range customerPropertyTypeEntities {
		tmpType, err := propertyTypeStringToInt(entity.ParamType())
		if err != nil {
			return 0, err
		}
		mapper := CustomerPropertyMapper{
			Name: entity.Name(),
			Type: tmpType,
		}
		customerPropertyMappers = append(customerPropertyMappers, &mapper)
	}
	// 保存実行
	if err := executor.Insert(customerPropertyMappers...); err != nil {
		return 0, errors.Wrapf(err, "customer_propertyテーブルへの保存に失敗しました。%v", customerTypeEntity)
	}

	////// customer_types_customer_propertiesの保存
	// mappperに詰める
	relations := make([]interface{}, 0, len(customerPropertyMappers))
	for index, propertyMapperInterface := range customerPropertyMappers {
		propertyMapper, ok := propertyMapperInterface.(*CustomerPropertyMapper)
		if !ok {
			return 0, errors.Wrapf(err, "*CustomerPropertyMapperへのキャストに失敗しました。%v", propertyMapperInterface)
		}
		ralationMapper := CustomerTypeCustomerPropertyMapper{
			CustomerTypeId:     customerTypeMapper.Id,
			CustomerPropertyId: propertyMapper.Id,
			Order:              index + 1,
		}
		relations = append(relations, &ralationMapper)
	}
	// 保存実行
	if err := executor.Insert(relations...); err != nil {
		return 0, errors.Wrapf(err, "customer_types_customer_propertiesテーブルへの保存に失敗しました。%v", customerTypeEntity)
	}

	return customerTypeMapper.Id, nil
}

const (
	PROPERTY_TYPE_STRING  = 0
	PROPERTY_TYPE_NUMERIC = 1
)

func propertyTypeStringToInt(strType string) (int, error) {
	switch strType {
	case customer.PROPERTY_TYPE_STRING:
		return PROPERTY_TYPE_STRING, nil
	case customer.PROPERTY_TYPE_NUMERIC:
		return PROPERTY_TYPE_NUMERIC, nil
	default:
		return -1, errors.Errorf("想定外の値が渡されました。strType: %v", strType)
	}
}

type CustomerTypeMapper struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	CreatedAtUpdatedAtMapper
}

type CustomerPropertyMapper struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	Type int    `db:"type""`
	CreatedAtUpdatedAtMapper
}

type CustomerTypeCustomerPropertyMapper struct {
	CustomerTypeId     int `db:"customer_type_id"`
	CustomerPropertyId int `db:"customer_property_id"`
	Order              int `db:"order"`
	CreatedAtUpdatedAtMapper
}

//func (r *ContractRepository) Create(contractEntity *contract.ContractEntity, executor gorp.SqlExecutor) (savedId int, err error) {
//	////// contractの保存
//	// data_mapperオブジェクトに詰め替え
//	contractMapper := ContractMapper{
//		UserId:                   contractEntity.UserId(),
//		ProductId:                contractEntity.ProductId(),
//		ContractDate:             contractEntity.ContractDate(),
//		BillingStartDate:         contractEntity.BillingStartDate(),
//		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{},
//	}
//	// 新規保存実行
//	err = executor.Insert(&contractMapper)
//	if err != nil {
//		return 0, errors.Wrapf(err, "contractsテーブルへの保存に失敗しました。%v", contractEntity)
//	}
//
//	// 使用権の保存
//	rightToUses := contractEntity.RightToUses()
//	for _, rightToUseEntity := range rightToUses {
//		////// rightToUseの保存
//		err := createRightToUse(rightToUseEntity, contractMapper.Id, executor)
//		if err != nil {
//			return 0, err
//		}
//	}
//
//	return contractMapper.Id, nil
//}
//
//func createRightToUse(rightToUseEntity *contract.RightToUseEntity, contractId int, executor gorp.SqlExecutor) error {
//	rightToUseMapper := RightToUseMapper{}
//	rightToUseMapper.ContractId = contractId
//	rightToUseMapper.ValidFrom = rightToUseEntity.ValidFrom()
//	rightToUseMapper.ValidTo = rightToUseEntity.ValidTo()
//	err := executor.Insert(&rightToUseMapper)
//	if err != nil {
//		return errors.Wrapf(err, "right_to_useテーブルへの保存に失敗しました。rightToUseMapper: %+v", rightToUseMapper)
//	}
//	////// rightToUseActiveの保存
//	activeMapper := RightToUseActiveMapper{}
//	activeMapper.RightToUseId = rightToUseMapper.Id
//	err = executor.Insert(&activeMapper)
//	if err != nil {
//		return errors.Wrapf(err, "right_to_use_rightテーブルへの保存に失敗しました。activeMapper: %+v", activeMapper)
//	}
//	return nil
//}
//
//func (r *ContractRepository) GetByIds(ids []int, executor gorp.SqlExecutor) (contracts []*contract.ContractEntity, err error) {
//	if len(ids) == 0 {
//		return nil, errors.Errorf("idsが空スライスです。ids: %+v", ids)
//	}
//	// データ取得
//	// contracts取得
//	// データマッパー用意
//	var contractMappers []*ContractMapper
//	// idsをインターフェース型に変更
//	idsInterfaceType := make([]interface{}, 0, len(ids))
//	preparedStatement := make([]string, 0, len(ids))
//	for i, id := range ids {
//		idsInterfaceType = append(idsInterfaceType, id)
//		preparedStatement = append(preparedStatement, "$"+strconv.Itoa(int(i)+1))
//	}
//	// sql作成
//	contractQuery := `
//select
//       c.id as id,
//       c.user_id as user_id,
//       c.product_id as product_id,
//       c.contract_date as contract_date,
//       c.billing_start_date as billing_start_date,
//       c.created_at as created_at,
//       c.updated_at as updated_at
//from contracts c
//where c.id IN (%v)
//order by c.id
//`
//	contractQuery = fmt.Sprintf(contractQuery, strings.Join(preparedStatement, ", "))
//	// sqlとデータマッパーでクエリ実行
//	_, err = executor.Select(&contractMappers, contractQuery, idsInterfaceType...)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, nil
//		} else {
//			return nil, errors.Wrapf(err, "契約情報取得失敗。contractQuery: %v, ids: %v", contractQuery, ids)
//		}
//	}
//
//	if len(contractMappers) == 0 {
//		// データが無い時
//		return nil, nil
//	}
//
//	retContracts := make([]*contract.ContractEntity, 0, len(contractMappers))
//	for _, mapper := range contractMappers {
//		// 使用権データを取得
//		rightToUseEntities, err := getRightToUseEntitiesByContractId(mapper.Id, executor)
//		if err != nil {
//			return nil, err
//		}
//		// contractEntityを作成
//		contractEntity, err := contract.NewContractEntityWithData(
//			mapper.Id,
//			mapper.UserId,
//			mapper.ProductId,
//			mapper.ContractDate,
//			mapper.BillingStartDate,
//			mapper.CreatedAt,
//			mapper.UpdatedAt,
//			rightToUseEntities,
//		)
//		if err != nil {
//			return nil, err
//		}
//		retContracts = append(retContracts, contractEntity)
//	}
//	return retContracts, nil
//}
//
//func getRightToUseEntitiesByContractId(contractId int, executor gorp.SqlExecutor) ([]*contract.RightToUseEntity, error) {
//	query := `
//SELECT
//	right_to_use.*,
//    COALESCE(bd.id, 0) as bill_detail_id
//FROM right_to_use
//INNER JOIN right_to_use_active rtua on right_to_use.id = rtua.right_to_use_id
//LEFT OUTER JOIN bill_details bd ON bd.right_to_use_id = right_to_use.id
//
//WHERE contract_id = $1
//ORDER BY id
//`
//	var mappers []*rightToUseMapper
//	_, err := executor.Select(&mappers, query, contractId)
//	if err != nil {
//		return nil, errors.Wrapf(err, "使用権データの取得に失敗しました。query: %v, contractId: %v", query, contractId)
//	}
//
//	retEntities := make([]*contract.RightToUseEntity, 0, len(mappers))
//	for _, mapper := range mappers {
//		entity := contract.NewRightToUseEntityWithData(
//			mapper.Id,
//			mapper.ValidFrom,
//			mapper.ValidTo,
//			mapper.BillDetailId,
//			mapper.CreatedAt,
//			mapper.UpdatedAt)
//		retEntities = append(retEntities, entity)
//	}
//
//	return retEntities, nil
//}
//
//type rightToUseMapper struct {
//	RightToUseMapper
//	BillDetailId int `db:"bill_detail_id"`
//}
//
//func (r *ContractRepository) GetById(id int, executor gorp.SqlExecutor) (contract *contract.ContractEntity, err error) {
//	contracts, err := r.GetByIds([]int{id}, executor)
//	if err != nil {
//		return nil, err
//	}
//	if contracts == nil {
//		// データが無い時
//		return nil, nil
//	} else {
//		return contracts[0], nil
//	}
//}
//
//func (r *ContractRepository) GetBillingTargetByBillingDate(billingDate time.Time, executor gorp.SqlExecutor) ([]*contract.ContractEntity, error) {
//	// 対象contractId取得クエリを用意する
//	query := `
//SELECT
//       c.id
//FROM contracts c
//    INNER JOIN right_to_use rtu on c.id = rtu.contract_id
//    INNER JOIN right_to_use_active rtua on rtu.id = rtua.right_to_use_id
//    LEFT OUTER JOIN bill_details bd on rtu.id = bd.right_to_use_id
//WHERE bd.id IS NULL
//  AND rtu.valid_from <= $1
//  AND c.billing_start_date <= $1
//`
//	// データ取得実行する
//	var targetIds []int
//	_, err := executor.Select(&targetIds, query, billingDate)
//	if err != nil {
//		return nil, errors.Wrapf(err, "請求対象契約のidの取得に失敗しました。query: %v, billingDate: %+v", query, billingDate)
//	}
//
//	if len(targetIds) == 0 {
//		// データがなかった時
//		return []*contract.ContractEntity{}, nil
//	} else {
//		contracts, err := r.GetByIds(targetIds, executor)
//		if err != nil {
//			return nil, err
//		}
//		return contracts, nil
//	}
//}
//
///*
//渡した日（実行日）から5日以内に終了し、かつ、まだ次の期間の使用権データが存在しない使用権をもつ契約集約を全て返す
//
//例）実行日が6/1の場合
//使用権の終了日が6/1の使用権=> 返る
//使用権の終了日が6/6の使用権=> 返る
//使用権の終了日が6/7の使用権=> 返らない
//使用権の終了日が6/1だが、次（6/2 ~ 7/1の期間）の使用権が存在する=> 返らない
//*/
//func (r *ContractRepository) GetRecurTargets(executeDate time.Time, executor gorp.SqlExecutor) ([]*contract.ContractEntity, error) {
//	from := executeDate
//	to := executeDate.AddDate(0, 0, 5)
//
//	query := `
//	WITH tmp_t AS (
//	   SELECT *, row_number() over (partition by contract_id order by valid_to DESC) AS num FROM right_to_use
//	)
//	SELECT
//		contract_id
//	FROM tmp_t
//	WHERE num = 1
//	AND $1 <= tmp_t.valid_to
//	AND tmp_t.valid_to < $2
//	GROUP BY contract_id
//	ORDER BY contract_id;
//	;`
//
//	var contractIds []int
//	var _, err = executor.Select(&contractIds, query, from, to)
//	if err != nil {
//		return nil, errors.Wrapf(err, "継続処理対象使用権をもつ契約IDの取得に失敗しました。query: %v, from: %v, to: %v", query, from, to)
//	}
//	if len(contractIds) == 0 {
//		// 対象がない
//		return []*contract.ContractEntity{}, nil
//	}
//
//	// 契約集約を取得
//	contracts, err := r.GetByIds(contractIds, executor)
//	if err != nil {
//		return nil, errors.Wrapf(err, "継続処理対象使用権をもつ契約集約の取得に失敗しました。contractIds: %v", contractIds)
//	}
//
//	return contracts, nil
//}
//
//// contractEntityを更新する（まだ使用権の追加しか対応してない）
//func (r *ContractRepository) Update(contractEntity *contract.ContractEntity, executor gorp.SqlExecutor) error {
//	// idがzeroの使用権があれば新規登録する
//	for _, rightToUse := range contractEntity.RightToUses() {
//		if rightToUse.Id() == 0 {
//			err := createRightToUse(rightToUse, contractEntity.Id(), executor)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	// アーカイブ指定の使用権があればhistoryテーブルに移す
//	for _, id := range contractEntity.GetToArchiveRightToUseIds() {
//		err := executeArchiveRightToUse(id, executor)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//// 指定idの使用権をright_to_use_activeからhistoryテーブルへ移す
//func executeArchiveRightToUse(id int, executor gorp.SqlExecutor) error {
//	// right_to_use_historyへのinsert
//	history := RightToUseHistoryMapper{}
//	history.RightToUseId = id
//	err := executor.Insert(&history)
//	if err != nil {
//		return errors.Wrapf(err, "right_to_use_historyテーブルへのinsert失敗。history: %+v", history)
//	}
//
//	// right_to_use_activeからの削除
//	active := RightToUseActiveMapper{}
//	active.RightToUseId = id
//	_, err = executor.Delete(&active)
//	if err != nil {
//		return errors.Wrapf(err, "right_to_use_activeテーブルからのdelete失敗。active: %+v", active)
//	}
//	return nil
//}
//
///*
//渡した基準日時点で期限が切れているactiveな使用権を持っている契約エンティティを返す
//*/
//func (r *ContractRepository) GetHavingExpiredRightToUseContractIds(baseDate time.Time, executor gorp.SqlExecutor) ([]int, error) {
//	query := `
//SELECT DISTINCT c.id
//FROM contracts c
//    INNER JOIN right_to_use rtu ON rtu.contract_id = c.id
//    INNER JOIN right_to_use_active rtua ON rtua.right_to_use_id = rtu.id
//WHERE rtu.valid_to <= $1
//ORDER BY c.id
//`
//	contractIds := []int{}
//	_, err := executor.Select(&contractIds, query, baseDate)
//	if err != nil {
//		return nil, errors.Wrapf(err, "期限切れ使用権保持契約idの取得失敗。query: %v", query)
//	}
//
//	return contractIds, nil
//}
//
//type ContractMapper struct {
//	Id               int       `db:"id"`
//	UserId           int       `db:"user_id"`
//	ProductId        int       `db:"product_id"`
//	ContractDate     time.Time `db:"contract_date"`
//	BillingStartDate time.Time `db:"billing_start_date"`
//	CreatedAtUpdatedAtMapper
//}
//
//func (c *ContractMapper) SetDataToEntity(entity interface{}) error {
//	value, ok := entity.(*contract.ContractEntity)
//	if !ok {
//		return errors.Errorf("*entities.ContractEntityではないものが渡された。entity: %T", entity)
//	}
//	err := value.LoadData(
//		c.Id,
//		c.UserId,
//		c.ProductId,
//		c.ContractDate,
//		c.BillingStartDate,
//		c.CreatedAt,
//		c.UpdatedAt,
//		[]*contract.RightToUseEntity{},
//	)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//type RightToUseMapper struct {
//	Id         int       `db:"id"`
//	ContractId int       `db:"contract_id"`
//	ValidFrom  time.Time `db:"valid_from"`
//	ValidTo    time.Time `db:"valid_to"`
//	CreatedAtUpdatedAtMapper
//}
//
//type RightToUseActiveMapper struct {
//	RightToUseId int `db:"right_to_use_id"`
//	CreatedAtUpdatedAtMapper
//}
//
//type RightToUseHistoryMapper struct {
//	RightToUseId int `db:"right_to_use_id"`
//	CreatedAtUpdatedAtMapper
//}

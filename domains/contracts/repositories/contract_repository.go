package repositories

import (
	"database/sql"
	"fmt"
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/data_mappers"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"strconv"
	"strings"
	"time"
)

type ContractRepository struct {
	*BaseRepository
}

func NewContractRepository() *ContractRepository {
	return &ContractRepository{
		&BaseRepository{},
	}
}

// 契約エンティティを新規保存する
func (r *ContractRepository) Create(contractEntity *entities.ContractEntity, executor gorp.SqlExecutor) (savedId int, err error) {
	////// contractの保存
	// data_mapperオブジェクトに詰め替え
	contractMapper := data_mappers.ContractMapper{
		UserId:                   contractEntity.UserId(),
		ProductId:                contractEntity.ProductId(),
		ContractDate:             contractEntity.ContractDate(),
		BillingStartDate:         contractEntity.BillingStartDate(),
		CreatedAtUpdatedAtMapper: data_mappers.CreatedAtUpdatedAtMapper{},
	}
	// 新規保存実行
	err = executor.Insert(&contractMapper)
	if err != nil {
		return 0, errors.Wrapf(err, "contractsテーブルへの保存に失敗しました。%v", contractEntity)
	}

	// 使用権の保存
	rightToUses := contractEntity.RightToUses()
	for _, rightToUseEntity := range rightToUses {
		////// rightToUseの保存
		err := createRightToUse(rightToUseEntity, contractMapper.Id, executor)
		if err != nil {
			return 0, err
		}
	}

	return contractMapper.Id, nil
}

func createRightToUse(rightToUseEntity *entities.RightToUseEntity, contractId int, executor gorp.SqlExecutor) error {
	rightToUseMapper := data_mappers.RightToUseMapper{}
	rightToUseMapper.ContractId = contractId
	rightToUseMapper.ValidFrom = rightToUseEntity.ValidFrom()
	rightToUseMapper.ValidTo = rightToUseEntity.ValidTo()
	err := executor.Insert(&rightToUseMapper)
	if err != nil {
		return errors.Wrapf(err, "right_to_useテーブルへの保存に失敗しました。rightToUseMapper: %+v", rightToUseMapper)
	}
	////// rightToUseActiveの保存
	activeMapper := data_mappers.RightToUseActiveMapper{}
	activeMapper.RightToUseId = rightToUseMapper.Id
	err = executor.Insert(&activeMapper)
	if err != nil {
		return errors.Wrapf(err, "right_to_use_rightテーブルへの保存に失敗しました。activeMapper: %+v", activeMapper)
	}
	return nil
}

func (r *ContractRepository) GetByIds(ids []int, executor gorp.SqlExecutor) (contracts []*entities.ContractEntity, products []*entities.ProductEntity, users []interface{}, err error) {
	// データ取得
	// データマッパー用意
	var mappers []*data_mappers.ContractView
	// idsをインターフェース型に変更
	idsInterfaceType := make([]interface{}, 0, len(ids))
	preparedStatement := make([]string, 0, len(ids))
	for i, id := range ids {
		idsInterfaceType = append(idsInterfaceType, id)
		preparedStatement = append(preparedStatement, "$"+strconv.Itoa(int(i)+1))
	}
	// sql作成
	query := `
select
       c.id as id,
       c.contract_date as contract_date,
       c.billing_start_date as billing_start_date,
       c.created_at as created_at,
       c.updated_at as updated_at,
       p.id as product_id,
       p.name as product_name,
       ppm.price as product_price,
       p.created_at as product_created_at,
       p.updated_at as product_updated_at,
       u.id as user_id,
       case
           when ui.user_id IS NOT NULL then 'individual'
           when uc.user_id IS NOT NULL then 'corporation'
        end as user_type,
       ui.name as user_individual_name,
       u.created_at as user_individual_created_at,
       u.updated_at as user_individual_updated_at,
       uc.corporation_name as user_corporation_corporation_name,
       uc.contact_person_name as user_corporation_contact_person_name,
       uc.president_name as user_corporation_president_name,
       u.created_at as user_corporation_created_at,
       u.updated_at as user_corporation_updated_at,
       rtu.id as right_to_use_id,
       rtu.valid_from as right_to_use_valid_from,
       rtu.valid_to as right_to_use_valid_to,
       rtua.created_at as right_to_use_active_created_at,
       rtua.updated_at as right_to_use_active_updated_at,
       COALESCE(bd.id, 0) as bill_detail_id
from contracts c
inner join products p on c.product_id = p.id
inner join product_price_monthlies ppm on ppm.product_id = p.id
inner join users u on c.user_id = u.id
inner join right_to_use rtu on rtu.contract_id = c.id
inner join right_to_use_active rtua on rtua.right_to_use_id = rtu.id
left outer join users_individual ui on u.id = ui.user_id
left outer join users_corporation uc on u.id = uc.user_id
left outer join bill_details bd on bd.right_to_use_id = rtu.id
where c.id IN (%v)
order by c.id, right_to_use_id`
	query = fmt.Sprintf(query, strings.Join(preparedStatement, ", "))
	// sqlとデータマッパーでクエリ実行
	_, err = executor.Select(&mappers, query, idsInterfaceType...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil
		} else {
			return nil, nil, nil, errors.Wrapf(err, "契約情報取得失敗。query: %v, ids: %v", query, ids)
		}
	}

	if len(mappers) == 0 {
		// データが無い時
		return nil, nil, nil, nil
	}

	// contract単位でmapperを分割する
	mappersByContract := separateMapper(mappers)
	// contract毎にentityを組み立てる
	for _, mappers := range mappersByContract {
		contract, product, user, err := createEntitiesFromMapper(mappers)
		if err != nil {
			return nil, nil, nil, err
		}
		contracts = append(contracts, contract)
		products = append(products, product)
		users = append(users, user)
	}
	return contracts, products, users, nil
}

// dbから取得したレコードをcontract単位で分割する
func separateMapper(mappers []*data_mappers.ContractView) [][]*data_mappers.ContractView {
	retMapper := [][]*data_mappers.ContractView{}
	index := -1
	prevContractId := 0
	for _, mapper := range mappers {
		if prevContractId != mapper.Id {
			prevContractId = mapper.Id
			index += 1
			retMapper = append(retMapper, []*data_mappers.ContractView{})
		}
		retMapper[index] = append(retMapper[index], mapper)
	}
	return retMapper
}

func (r *ContractRepository) GetById(id int, executor gorp.SqlExecutor) (contract *entities.ContractEntity, product *entities.ProductEntity, user interface{}, err error) {
	contracts, products, users, err := r.GetByIds([]int{id}, executor)
	if err != nil {
		return nil, nil, nil, err
	}
	if contracts == nil {
		// データが無い時
		return nil, nil, nil, nil
	} else {
		return contracts[0], products[0], users[0], nil
	}
}

// dbから取得したレコード情報からエンティティを組み上げる
func createEntitiesFromMapper(mappers []*data_mappers.ContractView) (
	contract *entities.ContractEntity,
	product *entities.ProductEntity,
	user interface{},
	err error,
) {
	// 使用権データ作成
	rightToUseEntities := make([]*entities.RightToUseEntity, 0, len(mappers))
	for _, mapper := range mappers {
		// 使用権エンティティにデータを詰める
		rightToUseEntities = append(rightToUseEntities, createRightToUseFromMapper(mapper))
	}

	// productエンティティにデータを詰める
	product, err = createProductEntityFromMapper(mappers[0])
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "productEntity作成失敗。mappers[0]: %v", mappers[0])
	}
	// contractエンティティにデータを詰める
	contract, err = createContractEntityFromMapper(mappers[0], rightToUseEntities)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "contractEntity作成失敗。mappers[0]: %v", mappers[0])
	}
	// userエンティティにデータを詰める
	user, err = createUserEntityFromMapper(mappers[0])
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "userEntity作成失敗。mappers[0]: %v", mappers[0])
	}

	return contract, product, user, nil
}

func createRightToUseFromMapper(mapper *data_mappers.ContractView) *entities.RightToUseEntity {
	return entities.NewRightToUseEntityWithData(
		mapper.RightToUseId,
		mapper.RightToUseValidFrom,
		mapper.RightToUseValidTo,
		mapper.BillDetailId,
		mapper.RightToUseCreatedAt,
		mapper.RightToUseUpdatedAt,
	)
}

func createProductEntityFromMapper(mapper *data_mappers.ContractView) (*entities.ProductEntity, error) {
	return entities.NewProductEntityWithData(
		mapper.ProductId,
		mapper.ProductName,
		mapper.ProductPrice.String(),
		mapper.ProductCreatedAt,
		mapper.ProductUpdatedAt,
	)
}

func createContractEntityFromMapper(mapper *data_mappers.ContractView, rightToUseEntities []*entities.RightToUseEntity) (*entities.ContractEntity, error) {
	return entities.NewContractEntityWithData(
		mapper.Id,
		mapper.UserId,
		mapper.ProductId,
		mapper.ContractDate,
		mapper.BillingStartDate,
		mapper.CreatedAt,
		mapper.UpdatedAt,
		rightToUseEntities,
	)
}

func createUserEntityFromMapper(mapper *data_mappers.ContractView) (user interface{}, err error) {
	switch mapper.UserType {
	case "individual":
		user, err = entities.NewUserIndividualEntityWithData(mapper.UserId, mapper.UserIndividualName.String, mapper.UserIndividualCreatedAt.Time, mapper.UserIndividualUpdatedAt.Time)
		if err != nil {
			return nil, errors.Wrapf(err, "userIndividualEntity作成失敗。mapper: %v", mapper)
		}
	case "corporation":
		user, err = entities.NewUserCorporationEntityWithData(mapper.UserId, mapper.UserCorporationCorporationName.String, mapper.UserCorporationContractPersonName.String, mapper.UserCorporationPresidentName.String, mapper.UserCorporationCreatedAt.Time, mapper.UserCorporationUpdatedAt.Time)
		if err != nil {
			return nil, errors.Wrapf(err, "userCorporationEntity作成失敗。mapper: %v", mapper)
		}
	default:
		return nil, errors.Errorf("考慮外のUserTypeが来た。mappers[0].UserType: %v, mapper: %v", mapper.UserType, mapper)
	}
	return user, nil
}

func (r *ContractRepository) GetBillingTargetByBillingDate(billingDate time.Time, executor gorp.SqlExecutor) ([]*entities.ContractEntity, error) {
	// 対象contractId取得クエリを用意する
	query := `
SELECT
       c.id
FROM contracts c
    INNER JOIN right_to_use rtu on c.id = rtu.contract_id
    INNER JOIN right_to_use_active rtua on rtu.id = rtua.right_to_use_id
    LEFT OUTER JOIN bill_details bd on rtu.id = bd.right_to_use_id
WHERE bd.id IS NULL
  AND rtu.valid_from <= $1
  AND c.billing_start_date <= $1
`
	// データ取得実行する
	var targetIds []int
	_, err := executor.Select(&targetIds, query, billingDate)
	if err != nil {
		return nil, errors.Wrapf(err, "請求対象契約のidの取得に失敗しました。query: %v, billingDate: %+v", query, billingDate)
	}

	if len(targetIds) == 0 {
		// データがなかった時
		return []*entities.ContractEntity{}, nil
	} else {
		contracts, _, _, err := r.GetByIds(targetIds, executor)
		if err != nil {
			return nil, err
		}
		return contracts, nil
	}
}

/*
渡した日（実行日）から5日以内に終了し、かつ、まだ次の期間の使用権データが存在しない使用権をもつ契約集約を全て返す

例）実行日が6/1の場合
使用権の終了日が6/1の使用権=> 返る
使用権の終了日が6/6の使用権=> 返る
使用権の終了日が6/7の使用権=> 返らない
使用権の終了日が6/1だが、次（6/2 ~ 7/1の期間）の使用権が存在する=> 返らない
*/
func (r *ContractRepository) GetRecurTargets(executeDate time.Time, executor gorp.SqlExecutor) ([]*entities.ContractEntity, error) {
	from := executeDate
	to := executeDate.AddDate(0, 0, 5)

	query := `
	WITH tmp_t AS (
	   SELECT *, row_number() over (partition by contract_id order by valid_to DESC) AS num FROM right_to_use
	)
	SELECT
		contract_id
	FROM tmp_t
	WHERE num = 1
	AND $1 <= tmp_t.valid_to
	AND tmp_t.valid_to < $2
	GROUP BY contract_id
	ORDER BY contract_id;
	;`

	var contractIds []int
	var _, err = executor.Select(&contractIds, query, from, to)
	if err != nil {
		return nil, errors.Wrapf(err, "継続処理対象使用権をもつ契約IDの取得に失敗しました。query: %v, from: %v, to: %v", query, from, to)
	}

	// 契約集約を取得
	contracts, _, _, err := r.GetByIds(contractIds, executor)
	if err != nil {
		return nil, errors.Wrapf(err, "継続処理対象使用権をもつ契約集約の取得に失敗しました。contractIds: %v", contractIds)
	}

	return contracts, nil
}

// contractEntityを更新する（まだ使用権の追加しか対応してない）
func (r *ContractRepository) Update(contractEntity *entities.ContractEntity, executor gorp.SqlExecutor) error {
	// idがzeroの使用権があれば新規登録する
	for _, rightToUse := range contractEntity.RightToUses() {
		if rightToUse.Id() == 0 {
			err := createRightToUse(rightToUse, contractEntity.Id(), executor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

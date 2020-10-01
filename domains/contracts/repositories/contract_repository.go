package repositories

import (
	"database/sql"
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/data_mappers"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
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
		rightToUseMapper := data_mappers.RightToUseMapper{}
		rightToUseMapper.ContractId = contractMapper.Id
		rightToUseMapper.ValidFrom = rightToUseEntity.ValidFrom()
		rightToUseMapper.ValidTo = rightToUseEntity.ValidTo()
		err := executor.Insert(&rightToUseMapper)
		if err != nil {
			return 0, errors.Wrapf(err, "right_to_useテーブルへの保存に失敗しました。rightToUseMapper: %+v", rightToUseMapper)
		}
		////// rightToUseActiveの保存
		activeMapper := data_mappers.RightToUseActiveMapper{}
		activeMapper.RightToUseId = rightToUseMapper.Id
		err = executor.Insert(&activeMapper)
		if err != nil {
			return 0, errors.Wrapf(err, "right_to_use_rightテーブルへの保存に失敗しました。activeMapper: %+v", activeMapper)
		}
	}

	return contractMapper.Id, nil
}

func (r *ContractRepository) GetById(id int, executor gorp.SqlExecutor) (contract *entities.ContractEntity, product *entities.ProductEntity, user interface{}, err error) {
	// データ取得
	// データマッパー用意
	var mappers []*data_mappers.ContractView
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
where c.id = $1
order by right_to_use_id`
	// sqlとデータマッパーでクエリ実行
	_, err = executor.Select(&mappers, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil
		} else {
			return nil, nil, nil, errors.Wrapf(err, "契約情報取得失敗。id: %v", id)
		}
	}

	if len(mappers) == 0 {
		// データが無い時
		return nil, nil, nil, nil
	}

	// 使用権データ作成
	rightToUseEntities := make([]*entities.RightToUseEntity, 0, len(mappers))
	for _, mapper := range mappers {
		// 使用権エンティティにデータを詰める
		entity := entities.NewRightToUseEntityWithData(
			mapper.RightToUseId,
			mapper.RightToUseValidFrom,
			mapper.RightToUseValidTo,
			mapper.BillDetailId,
			mapper.RightToUseCreatedAt,
			mapper.RightToUseUpdatedAt,
		)
		rightToUseEntities = append(rightToUseEntities, entity)
	}

	// productエンティティにデータを詰める
	product, err = entities.NewProductEntityWithData(mappers[0].ProductId, mappers[0].ProductName, mappers[0].ProductPrice.String(), mappers[0].ProductCreatedAt, mappers[0].ProductUpdatedAt)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "productEntity作成失敗。mappers[0]: %v", mappers[0])
	}
	// contractエンティティにデータを詰める
	contract, err = entities.NewContractEntityWithData(mappers[0].Id, mappers[0].UserId, mappers[0].ProductId, mappers[0].ContractDate, mappers[0].BillingStartDate, mappers[0].CreatedAt, mappers[0].UpdatedAt, rightToUseEntities)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "contractEntity作成失敗。mappers[0]: %v", mappers[0])
	}
	// userエンティティにデータを詰める
	switch mappers[0].UserType {
	case "individual":
		user, err = entities.NewUserIndividualEntityWithData(mappers[0].UserId, mappers[0].UserIndividualName.String, mappers[0].UserIndividualCreatedAt.Time, mappers[0].UserIndividualUpdatedAt.Time)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "userIndividualEntity作成失敗。mappers[0]: %v", mappers[0])
		}
	case "corporation":
		user, err = entities.NewUserCorporationEntityWithData(mappers[0].UserId, mappers[0].UserCorporationCorporationName.String, mappers[0].UserCorporationContractPersonName.String, mappers[0].UserCorporationPresidentName.String, mappers[0].UserCorporationCreatedAt.Time, mappers[0].UserCorporationUpdatedAt.Time)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "userCorporationEntity作成失敗。mappers[0]: %v", mappers[0])
		}
	default:
		return nil, nil, nil, errors.Errorf("考慮外のUserTypeが来た。mappers[0].UserType: %v, mappet: %v", mappers[0].UserType, mappers[0])
	}

	return contract, product, user, nil
}

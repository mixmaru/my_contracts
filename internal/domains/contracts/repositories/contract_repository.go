package repositories

import (
	"database/sql"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
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
	// data_mapperオブジェクトに詰め替え
	contractMapper := data_mappers.ContractMapper{
		UserId:                   contractEntity.UserId(),
		ProductId:                contractEntity.ProductId(),
		CreatedAtUpdatedAtMapper: data_mappers.CreatedAtUpdatedAtMapper{},
	}

	// 新規保存実行
	err = executor.Insert(&contractMapper)
	if err != nil {
		return 0, errors.Wrapf(err, "contractsテーブルへの保存に失敗しました。%v", contractEntity)
	}

	return contractMapper.Id, nil
}

func (r *ContractRepository) GetById(id int, executor gorp.SqlExecutor) (contract *entities.ContractEntity, product *entities.ProductEntity, user interface{}, err error) {
	// データ取得
	// データマッパー用意
	var mapper data_mappers.ContractView
	// sql作成
	query :=
		`select
       c.id as id,
       c.created_at as created_at,
       c.updated_at as updated_at,
       p.id as product_id,
       p.name as product_name,
       p.price as product_price,
       p.created_at as product_created_at,
       p.updated_at as product_updated_at,
       u.id as user_id,
       case
           when ui.user_id IS NOT NULL then 'individual'
           when uc.user_id IS NOT NULL then 'corporation'
        end as user_type,
       ui.name as user_individual_name,
       ui.created_at as user_individual_created_at,
       ui.updated_at as user_individual_updated_at,
       uc.contact_person_name as user_corporation_contact_person_name,
       uc.president_name as user_corporation_president_name,
       uc.created_at as user_corporation_created_at,
       uc.updated_at as user_corporation_updated_at
from contracts c
inner join products p on c.product_id = p.id
inner join users u on c.user_id = u.id
left outer join users_individual ui on u.id = ui.user_id
left outer join users_corporation uc on u.id = uc.user_id
where c.id = $1`
	// sqlとデータマッパーでクエリ実行
	err = executor.SelectOne(&mapper, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil
		} else {
			return nil, nil, nil, errors.Wrapf(err, "契約情報取得失敗。id: %v", id)
		}
	}
	// productエンティティにデータを詰める
	product, err = entities.NewProductEntityWithData(mapper.ProductId, mapper.ProductName, mapper.ProductPrice.String(), mapper.ProductCreatedAt, mapper.ProductUpdatedAt)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "productEntity作成失敗。mapper: %v", mapper)
	}
	// contractエンティティにデータを詰める
	contract, err = entities.NewContractEntityWithData(mapper.Id, mapper.UserId, mapper.ProductId, mapper.CreatedAt, mapper.UpdatedAt)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "contractEntity作成失敗。mapper: %v", mapper)
	}
	// userエンティティにデータを詰める
	switch mapper.UserType {
	case "individual":
		user, err = entities.NewUserIndividualEntityWithData(mapper.UserId, mapper.UserIndividualName.String, mapper.UserIndividualCreatedAt.Time, mapper.UserIndividualUpdatedAt.Time)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "userIndividualEntity作成失敗。mapper: %v", mapper)
		}
	case "corporation":
		user, err = entities.NewUserCorporationEntityWithData(mapper.UserId, mapper.UserCorporationContractPersonName.String, mapper.UserCorporationPresidentName.String, mapper.UserCorporationCreatedAt.Time, mapper.UserCorporationUpdatedAt.Time)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "userCorporationEntity作成失敗。mapper: %v", mapper)
		}
	default:
		return nil, nil, nil, errors.Errorf("考慮外のUserTypeが来た。mapper.UserType: %v, mappet: %v", mapper.UserType, mapper)
	}

	return contract, product, user, nil
}

//func (r *ProductRepository) GetByName(name string, executor gorp.SqlExecutor) (*entities.ProductEntity, error) {
//	// データ取得
//	var productRecord data_mappers.ProductMapper
//	var productEntity entities.ProductEntity
//	noRow, err := r.selectOne(executor, &productRecord, &productEntity, "select * from products where name = $1", name)
//	if err != nil {
//		return nil, err
//	}
//	if noRow {
//		return nil, nil
//	}
//	return &productEntity, nil
//}

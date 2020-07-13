package repositories

import (
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

func (r *ContractRepository) GetById(id int, executor gorp.SqlExecutor) (*entities.ContractEntity, error) {
	// データ取得
	var mapper data_mappers.ContractMapper
	var entity entities.ContractEntity
	noRow, err := r.selectOne(executor, &mapper, &entity, "select * from contracts where id = $1", id)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return &entity, nil
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

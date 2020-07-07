package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type ProductRepository struct {
}

// 商品エンティティを新規保存する
func (r *ProductRepository) Save(productEntity *entities.ProductEntity, executor gorp.SqlExecutor) (*entities.ProductEntity, error) {
	// recordオブジェクトに詰め替え
	productRecord := data_mappers.ProductMapper{
		Name:  productEntity.Name(),
		Price: productEntity.Price(),
	}

	// 新規保存実行
	err := executor.Insert(&productRecord)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 再取得
	err = executor.SelectOne(&productRecord, "select * from products where id = $1", productRecord.Id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// エンティティに詰め直し
	err = productEntity.LoadData(
		productRecord.Id,
		productRecord.Name,
		productRecord.Price.String(),
		productRecord.CreatedAt,
		productRecord.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return productEntity, nil
}

func (r *ProductRepository) GetById(id int, executor gorp.SqlExecutor) (*entities.ProductEntity, error) {
	// データ取得
	var productRecord data_mappers.ProductMapper
	var productEntity entities.ProductEntity
	noRow, err := selectOne(executor, &productRecord, &productEntity, "select * from products where id = $1", id)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return &productEntity, nil
}

func (r *ProductRepository) GetByName(name string, executor gorp.SqlExecutor) (*entities.ProductEntity, error) {
	// データ取得
	var productRecord data_mappers.ProductMapper
	var productEntity entities.ProductEntity
	noRow, err := selectOne(executor, &productRecord, &productEntity, "select * from products where name = $1", name)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return &productEntity, nil
}

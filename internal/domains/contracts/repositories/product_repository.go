package repositories

import (
	"database/sql"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	tables2 "github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/tables"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type ProductRepository struct {
}

// 商品エンティティを新規保存する
func (r *ProductRepository) Save(productEntity *entities.ProductEntity, transaction *gorp.Transaction) (*entities.ProductEntity, error) {
	// db接続
	conn, err := db_connection.GetConnectionIfNotTransaction(transaction)
	if err != nil {
		return nil, err
	}
	defer db_connection.CloseConnectionIfNotTransaction(conn)

	// recordオブジェクトに詰め替え
	productRecord := tables2.ProductRecord{
		Name:  productEntity.Name(),
		Price: productEntity.Price(),
	}

	// 新規保存実行
	err = conn.Insert(&productRecord)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 再取得
	err = conn.SelectOne(&productRecord, "select * from products where id = $1", productRecord.Id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// エンティティに詰め直し
	productEntity.LoadData(
		productRecord.Id,
		productRecord.Name,
		productRecord.Price,
		productRecord.CreatedAt,
		productRecord.UpdatedAt,
	)
	return productEntity, nil
}

func (r *ProductRepository) GetById(id int, transaction *gorp.Transaction) (*entities.ProductEntity, error) {
	// db接続
	conn, err := db_connection.GetConnectionIfNotTransaction(transaction)
	if err != nil {
		return nil, err
	}
	defer db_connection.CloseConnectionIfNotTransaction(conn)

	// データ取得
	var productRecord tables2.ProductRecord
	var productEntity entities.ProductEntity
	noRow, err := r.selectOne(conn, &productRecord, &productEntity, "select * from products where id = $1", id)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return &productEntity, nil
}

func (r *ProductRepository) selectOne(executor gorp.SqlExecutor, record tables2.IRecord, entity entities.IBaseEntity, query string, args ...interface{}) (noRow bool, err error) {
	err = executor.SelectOne(record, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			// データがない
			return true, nil
		} else {
			return true, errors.WithStack(err)
		}
	}

	// エンティティに詰める
	err = record.SetDataToEntity(entity)
	if err != nil {
		return true, err
	}

	return false, nil
}

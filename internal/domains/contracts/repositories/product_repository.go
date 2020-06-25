package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/tables"
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
	productRecord := tables.ProductRecord{
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
	err = productEntity.LoadData(
		productRecord.Id,
		productRecord.Name,
		productRecord.Price,
		productRecord.CreatedAt,
		productRecord.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
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
	var productRecord tables.ProductRecord
	var productEntity entities.ProductEntity
	noRow, err := selectOne(conn, &productRecord, &productEntity, "select * from products where id = $1", id)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return &productEntity, nil
}

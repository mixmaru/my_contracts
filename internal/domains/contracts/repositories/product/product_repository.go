package product

import (
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

	// エンティティに詰め直し
	return nil, nil
}

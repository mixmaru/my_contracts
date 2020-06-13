package product

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	tables2 "github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/tables"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type ProductRepository struct {
}

// 商品エンティティを保存する
func (r *ProductRepository) Save(productEntity interface{}, transaction *gorp.Transaction) (interface{}, error) {
	conn, err := db_connection.GetConnectionIfNotTransaction(transaction)
	if err != nil {
		return nil, err
	}
	defer db_connection.CloseConnectionIfNotTransaction(conn)

	productRecord := tables2.ProductRecord{Name: "商品名4", Price: decimal.NewFromFloat(1000)}
	err = conn.Insert(&productRecord)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return nil, nil
}

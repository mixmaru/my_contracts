package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type ProductRepository struct {
	*BaseRepository
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		&BaseRepository{},
	}
}

// 商品エンティティを新規保存する
func (r *ProductRepository) Save(productEntity *entities.ProductEntity, executor gorp.SqlExecutor) (savedId int, err error) {
	// recordオブジェクトに詰め替え
	productRecord := data_mappers.ProductMapper{
		Name: productEntity.Name(),
	}

	// 新規保存実行
	err = executor.Insert(&productRecord)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	price, err := productEntity.MonthlyPrice()
	if err != nil {
		return 0, errors.WithMessagef(err, "月額金額取得失敗。productEntity: %v", productEntity)
	}

	productPriceMonthlyRecord := data_mappers.ProductPriceMonthlyMapper{
		ProductId: productRecord.Id,
		Price:     price,
	}
	err = executor.Insert(&productPriceMonthlyRecord)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return productRecord.Id, nil
}

func (r *ProductRepository) GetById(id int, executor gorp.SqlExecutor) (*entities.ProductEntity, error) {
	// データ取得
	var productRecord data_mappers.ProductMapper
	var productEntity entities.ProductEntity
	noRow, err := r.selectOne(executor, &productRecord, &productEntity, "select * from products where id = $1", id)
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
	noRow, err := r.selectOne(executor, &productRecord, &productEntity, "select * from products where name = $1", name)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return &productEntity, nil
}

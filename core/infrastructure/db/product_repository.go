package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/mixmaru/my_contracts/lib/decimal"
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
func (r *ProductRepository) Save(productEntity *product.ProductEntity, executor gorp.SqlExecutor) (savedId int, err error) {
	// recordオブジェクトに詰め替え
	productRecord := ProductMapper{
		Name: productEntity.Name(),
	}

	// 新規保存実行
	err = executor.Insert(&productRecord)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	price, exist := productEntity.MonthlyPrice()
	if !exist {
		return 0, errors.WithMessagef(err, "月額金額取得失敗。productEntity: %v", productEntity)
	}

	productPriceMonthlyRecord := ProductPriceMonthlyMapper{
		ProductId: productRecord.Id,
		Price:     price,
	}
	err = executor.Insert(&productPriceMonthlyRecord)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return productRecord.Id, nil
}

type ProductMapper struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	CreatedAtUpdatedAtMapper
}

type ProductPriceMonthlyMapper struct {
	ProductId int             `db:"product_id"`
	Price     decimal.Decimal `db:"price"`
	CreatedAtUpdatedAtMapper
}

func (r *ProductRepository) GetById(id int, executor gorp.SqlExecutor) (*product.ProductEntity, error) {
	// データ取得
	var productRecord productGetMapper
	var productEntity product.ProductEntity
	query := `
SELECT
       id,
       name,
       p.created_at,
       p.updated_at,
       CASE
           WHEN ppm.product_id IS NULL THEN false
           ELSE true
       END AS exist_price_monthly,
       ppm.price AS price_monthly
FROM products p
LEFT OUTER JOIN product_price_monthlies ppm on p.id = ppm.product_id
WHERE p.id = $1
`
	noRow, err := r.selectOne(executor, &productRecord, &productEntity, query, id)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return &productEntity, nil
}

func (r *ProductRepository) GetByName(name string, executor gorp.SqlExecutor) ([]*product.ProductEntity, error) {
	// データ取得
	var productRecords []*productGetMapper
	query := `
SELECT
       id,
       name,
       p.created_at,
       p.updated_at,
       CASE
           WHEN ppm.product_id IS NULL THEN false
           ELSE true
       END AS exist_price_monthly,
       ppm.price AS price_monthly
FROM products p
LEFT OUTER JOIN product_price_monthlies ppm on p.id = ppm.product_id
WHERE p.name = $1
ORDER BY id
`
	_, err := executor.Select(&productRecords, query, name)
	if err != nil {
		return nil, errors.Wrapf(err, "商品データの取得失敗。query: %v", query)
	}
	retEntities := make([]*product.ProductEntity, 0, len(productRecords))
	for _, record := range productRecords {
		entity, err := product.NewProductEntityWithData(record.Id, record.Name, record.PriceMonthly.Decimal.String(), record.CreatedAt, record.UpdatedAt)
		if err != nil {
			return nil, errors.Wrapf(err, "商品データEntityの組み立て失敗。record: %+v", record)
		}
		retEntities = append(retEntities, entity)
	}
	return retEntities, nil
}

func (r *ProductRepository) GetByRightToUseId(rightToUseId int, executor gorp.SqlExecutor) (*product.ProductEntity, error) {
	// データ取得
	var productRecord productGetMapper
	var productEntity product.ProductEntity
	query := `
SELECT
       p.id AS id,
       p.name AS name,
       p.created_at,
       p.updated_at,
       CASE
           WHEN ppm.product_id IS NULL THEN false
           ELSE true
       END AS exist_price_monthly,
       ppm.price AS price_monthly
FROM products p
INNER JOIN contracts c ON c.product_id = p.id
INNER JOIN right_to_use rtu ON rtu.contract_id = c.id
LEFT OUTER JOIN product_price_monthlies ppm on p.id = ppm.product_id
WHERE rtu.id = $1
`
	noRow, err := r.selectOne(executor, &productRecord, &productEntity, query, rightToUseId)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return &productEntity, nil
}

type productGetMapper struct {
	ProductMapper

	ExistPriceMonthly bool                `db:"exist_price_monthly"`
	PriceMonthly      decimal.NullDecimal `db:"price_monthly"`
}

func (p *productGetMapper) SetDataToEntity(productEntity interface{}) error {
	entity, ok := productEntity.(*product.ProductEntity)
	if !ok {
		return errors.Errorf("想定外の型が来た。型: %T, productEntity: %+v", productEntity, productEntity)
	}
	err := entity.LoadData(p.Id, p.Name, p.PriceMonthly.Decimal.String(), p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

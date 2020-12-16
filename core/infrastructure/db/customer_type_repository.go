package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type CustomerTypeRepository struct{}

func NewCustomerTypeRepository() *CustomerTypeRepository {
	return &CustomerTypeRepository{}
}

// カスタマータイプを新規保存する
func (r *CustomerTypeRepository) Create(customerTypeEntity *customer.CustomerTypeEntity, executor gorp.SqlExecutor) (savedId int, err error) {
	////// customer_typeの保存
	// mappperに詰める
	customerTypeMapper := CustomerTypeMapper{
		Name: customerTypeEntity.Name(),
	}
	// 保存実行
	if err := executor.Insert(&customerTypeMapper); err != nil {
		return 0, errors.Wrapf(err, "customer_typeテーブルへの保存に失敗しました。%v", customerTypeEntity)
	}

	////// customer_types_customer_propertiesの保存
	// mappperに詰める
	relations := make([]interface{}, 0, len(customerTypeEntity.CustomerPropertyTypeIds()))
	for index, propertyId := range customerTypeEntity.CustomerPropertyTypeIds() {
		ralationMapper := CustomerTypeCustomerPropertyMapper{
			CustomerTypeId:     customerTypeMapper.Id,
			CustomerPropertyId: propertyId,
			Order:              index + 1,
		}
		relations = append(relations, &ralationMapper)
	}
	// 保存実行
	if err := executor.Insert(relations...); err != nil {
		return 0, errors.Wrapf(err, "customer_types_customer_propertiesテーブルへの保存に失敗しました。%v", customerTypeEntity)
	}

	return customerTypeMapper.Id, nil
}

func (r *CustomerTypeRepository) GetById(id int, executor gorp.SqlExecutor) (entity *customer.CustomerTypeEntity, err error) {
	query := `
SELECT
	id,
	name,
    ctcp.customer_property_id
FROM customer_types ct
INNER JOIN customer_types_customer_properties ctcp on ct.id = ctcp.customer_type_id
WHERE ct.id = $1
ORDER BY ctcp.order;
`
	mappers := []*CustomerTypeForLoadMapper{}
	if _, err := executor.Select(&mappers, query, id); err != nil {
		return nil, errors.Wrapf(err, "dbからデータの取得に失敗しました。query: %v, id: %v, mappers: %+v", query, id, mappers)
	}
	if len(mappers) == 0 {
		// データが存在しない
		return nil, nil
	}
	// entityに詰める
	customerPropertyIds := make([]int, 0, len(mappers))
	for _, mapper := range mappers {
		customerPropertyIds = append(customerPropertyIds, mapper.CustomerPropertyId)
	}
	retEntity := customer.NewCustomerTypeEntityWithData(mappers[0].Id, mappers[0].Name, customerPropertyIds)
	return retEntity, nil
}

type CustomerTypeMapper struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	CreatedAtUpdatedAtMapper
}

type CustomerTypeCustomerPropertyMapper struct {
	CustomerTypeId     int `db:"customer_type_id"`
	CustomerPropertyId int `db:"customer_property_id"`
	Order              int `db:"order"`
	CreatedAtUpdatedAtMapper
}

type CustomerTypeForLoadMapper struct {
	Id                 int    `db:"id"`
	Name               string `db:"name"`
	CustomerPropertyId int    `db:"customer_property_id"`
}

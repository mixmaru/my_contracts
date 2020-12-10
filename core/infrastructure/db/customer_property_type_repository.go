package db

import (
	"fmt"
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"strconv"
	"strings"
)

type CustomerPropertyTypeRepository struct{}

func NewCustomerPropertyTypeRepository() *CustomerPropertyTypeRepository {
	return &CustomerPropertyTypeRepository{}
}

// カスタマープロパティタイプを新規保存する
func (r *CustomerPropertyTypeRepository) Create(entities []*customer.CustomerPropertyTypeEntity, executor gorp.SqlExecutor) (savedIds []int, err error) {
	////// costomer_propertiesの保存
	// mappperに詰める
	customerPropertyMappers := make([]interface{}, 0, len(entities))
	for _, entity := range entities {
		mapper := CustomerPropertyMapper{
			Name: entity.Name(),
			Type: int(entity.PropertyType()),
		}
		customerPropertyMappers = append(customerPropertyMappers, &mapper)
	}
	// 保存実行
	if err := executor.Insert(customerPropertyMappers...); err != nil {
		return nil, errors.Wrapf(err, "customer_propertyテーブルへの保存に失敗しました。%v", entities)
	}
	for _, mapperInterface := range customerPropertyMappers {
		mapper, ok := mapperInterface.(*CustomerPropertyMapper)
		if !ok {
			return nil, errors.Errorf("mapperInterfaceの型アサーションに失敗しました。%v", mapperInterface)
		}
		savedIds = append(savedIds, mapper.Id)
	}

	return savedIds, nil
}

// カスタマープロパティタイプIdを取得する
func (r *CustomerPropertyTypeRepository) GetByIds(ids []int, executor gorp.SqlExecutor) (propertyTypes []*customer.CustomerPropertyTypeEntity, err error) {
	idsInterfaceType := make([]interface{}, 0, len(ids))
	preparedStatement := make([]string, 0, len(ids))
	for i, id := range ids {
		idsInterfaceType = append(idsInterfaceType, id)
		preparedStatement = append(preparedStatement, "$"+strconv.Itoa(int(i)+1))
	}

	query := "SELECT id, name, type FROM customer_properties WHERE id IN (%v) ORDER BY id;"
	query = fmt.Sprintf(query, strings.Join(preparedStatement, ", "))

	// データ取得実行
	propertyTypeMappers := []*CustomerPropertyMapper{}
	_, err = executor.Select(&propertyTypeMappers, query, idsInterfaceType...)
	if err != nil {
		return nil, errors.Wrapf(err, "DBからデータ取得に失敗しました。query: %v, idsInterfaceType: %v", query, idsInterfaceType)
	}

	// エンティティに詰める
	for _, mapper := range propertyTypeMappers {
		entity := customer.NewCustomerParamTypeEntityWithData(mapper.Id, mapper.Name, customer.PropertyType(mapper.Type))
		propertyTypes = append(propertyTypes, entity)
	}

	return propertyTypes, nil
}

type CustomerPropertyMapper struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	Type int    `db:"type"`
	CreatedAtUpdatedAtMapper
}

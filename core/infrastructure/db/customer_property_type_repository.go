package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
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
		tmpType, err := propertyTypeStringToInt(entity.ParamType())
		if err != nil {
			return nil, err
		}
		mapper := CustomerPropertyMapper{
			Name: entity.Name(),
			Type: tmpType,
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

const (
	PROPERTY_TYPE_STRING  = 0
	PROPERTY_TYPE_NUMERIC = 1
)

func propertyTypeStringToInt(strType string) (int, error) {
	switch strType {
	case customer.PROPERTY_TYPE_STRING:
		return PROPERTY_TYPE_STRING, nil
	case customer.PROPERTY_TYPE_NUMERIC:
		return PROPERTY_TYPE_NUMERIC, nil
	default:
		return -1, errors.Errorf("想定外の値が渡されました。strType: %v", strType)
	}
}

type CustomerPropertyMapper struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	Type int    `db:"type""`
	CreatedAtUpdatedAtMapper
}

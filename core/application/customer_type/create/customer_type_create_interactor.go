package create

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/application/customer_type"
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type CustomerTypeCreateInteractor struct {
	customerTypeRepository         customer_type.ICustomerTypeRepository
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository
}

func NewCustomerTypeCreateInteractor(
	customerTypeRepository customer_type.ICustomerTypeRepository,
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository,
) *CustomerTypeCreateInteractor {
	return &CustomerTypeCreateInteractor{customerTypeRepository: customerTypeRepository, customerPropertyTypeRepository: customerPropertyTypeRepository}
}

func (c CustomerTypeCreateInteractor) Handle(request *CustomerTypeCreateUseCaseRequest) (*CustomerTypeCreateUseCaseResponse, error) {
	response := &CustomerTypeCreateUseCaseResponse{}

	// トランザクション開始
	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()
	tran, err := conn.Begin()
	if err != nil {
		return nil, err
	}

	// バリデーション
	response.ValidationError = c.validation(request, tran)
	if len(response.ValidationError) > 0 {
		return response, nil
	}

	// Entity用意
	entity := customer.NewCustomerTypeEntity(request.Name, request.CustomerParamTypeIds)

	// カスタマータイプ保存
	savedId, err := c.customerTypeRepository.Create(entity, tran)
	if err != nil {
		tran.Rollback()
		return nil, err
	}

	// リロード
	reloadedEntity, err := c.customerTypeRepository.GetById(savedId, tran)
	if err != nil {
		tran.Rollback()
		return nil, err
	}
	reloadedPropertyEntities, err := c.customerPropertyTypeRepository.GetByIds(reloadedEntity.CustomerPropertyTypeIds(), tran)
	if err != nil {
		tran.Rollback()
		return nil, err
	}

	// コミット
	if err := tran.Commit(); err != nil {
		tran.Rollback()
		return nil, errors.Wrapf(err, "トランザクションコミットに失敗しました。request: %+v", request)
	}

	// 返却dtoに詰める
	response.CustomerTypeDto, err = customer_type.NewCustomerTypeDtoFromEntity(reloadedEntity, reloadedPropertyEntities)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *CustomerTypeCreateInteractor) validation(request *CustomerTypeCreateUseCaseRequest, executor gorp.SqlExecutor) map[string][]string {
	var retMap map[string][]string

	// name重複チェック

	return retMap
}

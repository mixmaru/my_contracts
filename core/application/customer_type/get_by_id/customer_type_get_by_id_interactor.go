package get_by_id

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/application/customer_type"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
)

type CustomerTypeGetByIdInteractor struct {
	customerTypeRepository         customer_type.ICustomerTypeRepository
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository
}

func NewCustomerTypeGetByIdInteractor(customerTypeRepository customer_type.ICustomerTypeRepository, customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository) *CustomerTypeGetByIdInteractor {
	return &CustomerTypeGetByIdInteractor{customerTypeRepository: customerTypeRepository, customerPropertyTypeRepository: customerPropertyTypeRepository}
}

func (i *CustomerTypeGetByIdInteractor) Handle(request *CustomerTypeGetByIdUseCaseRequest) (*CustomerTypeGetByIdUseCaseResponse, error) {
	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	// カスタマータイプデータ取得
	entity, err := i.customerTypeRepository.GetById(request.CustomerTypeId, conn)
	if err != nil {
		return nil, err
	}

	// カスタマープロパティタイプデータ取得
	propertyTypeEntities, err := i.customerPropertyTypeRepository.GetByIds(entity.CustomerPropertyTypeIds(), conn)
	if err != nil {
		return nil, err
	}

	// 返却用dtoに詰める
	response := CustomerTypeGetByIdUseCaseResponse{}
	response.CustomerTypeDto, err = customer_type.NewCustomerTypeDtoFromEntity(entity, propertyTypeEntities)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

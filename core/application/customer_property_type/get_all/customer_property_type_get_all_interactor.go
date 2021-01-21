package get_all

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
)

type CustomerPropertyTypeGetAllInteractor struct {
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository
}

func NewCustomerPropertyTypeGetAllInteractor(customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository) *CustomerPropertyTypeGetAllInteractor {
	return &CustomerPropertyTypeGetAllInteractor{customerPropertyTypeRepository: customerPropertyTypeRepository}
}

func (cont *CustomerPropertyTypeGetAllInteractor) Handle() (*CustomerPropertyTypeGetAllUseCaseResponse, error) {
	response := CustomerPropertyTypeGetAllUseCaseResponse{}

	// リポジトリからデータ取得
	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	entities, err := cont.customerPropertyTypeRepository.GetAll(conn)
	for _, entity := range entities {
		dto, err := customer_property_type.NewCustomerPropertyTypeDtoFromEntity(entity)
		if err != nil {
			return nil, err
		}
		response.CustomerPropertyTypeDtos = append(response.CustomerPropertyTypeDtos, dto)
	}

	return &response, nil
}

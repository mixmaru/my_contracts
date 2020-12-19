package get_by_ids

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
)

type CustomerPropertyTypeGetByIdsInteractor struct {
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository
}

func NewCustomerPropertyTypeGetByIdsInteractor(
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository,
) *CustomerPropertyTypeGetByIdsInteractor {
	return &CustomerPropertyTypeGetByIdsInteractor{
		customerPropertyTypeRepository: customerPropertyTypeRepository,
	}
}

func (i *CustomerPropertyTypeGetByIdsInteractor) Handle(
	request *CustomerPropertyTypeGetByIdsUseCaseRequest,
) (*CustomerPropertyTypeGetByIdsUseCaseResponse, error) {
	response := CustomerPropertyTypeGetByIdsUseCaseResponse{}

	// リポジトリ使ってデータ取得
	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	entities, err := i.customerPropertyTypeRepository.GetByIds(request.Ids, conn)
	if err != nil {
		return nil, err
	}

	// 返却用dtoに詰める
	for _, entity := range entities {
		dto, err := customer_property_type.NewCustomerPropertyTypeDtoFromEntity(entity)
		if err != nil {
			return nil, err
		}
		response.CustomerPropertyTypeDtos = append(response.CustomerPropertyTypeDtos, dto)
	}

	// 返却
	return &response, nil
}

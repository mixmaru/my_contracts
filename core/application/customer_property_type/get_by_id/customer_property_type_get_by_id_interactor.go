package get_by_id

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
)

type CustomerPropertyTypeGetByIdInteractor struct {
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository
}

func NewCustomerPropertyTypeGetByIdInteractor(
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository,
) *CustomerPropertyTypeGetByIdInteractor {
	return &CustomerPropertyTypeGetByIdInteractor{
		customerPropertyTypeRepository: customerPropertyTypeRepository,
	}
}

func (i *CustomerPropertyTypeGetByIdInteractor) Handle(
	request *CustomerPropertyTypeGetByIdUseCaseRequest,
) (*CustomerPropertyTypeGetByIdUseCaseResponse, error) {
	response := CustomerPropertyTypeGetByIdUseCaseResponse{}

	// リポジトリ使ってデータ取得
	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	entities, err := i.customerPropertyTypeRepository.GetByIds([]int{request.Id}, conn)
	if err != nil {
		return nil, err
	}
	// データがなければゼロ値で返す
	if len(entities) == 0 {
		return &response, nil
	} else {
		// 返却用dtoに詰める
		response.CustomerPropertyTypeDto, err = customer_property_type.NewCustomerPropertyTypeDtoFromEntity(entities[0])
		if err != nil {
			return nil, err
		}

		// 返却
		return &response, nil
	}
}

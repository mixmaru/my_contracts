package get_by_id

import (
	"github.com/mixmaru/my_contracts/core/application/customer"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
)

type CustomerGetByIdInteractor struct {
	customerRepository customer.ICustomerRepository
}

func NewCustomerGetByIdInteractor(customerRepository customer.ICustomerRepository) *CustomerGetByIdInteractor {
	return &CustomerGetByIdInteractor{customerRepository: customerRepository}
}

func (c CustomerGetByIdInteractor) Handle(request CustomerGetByIdUseCaseRequest) (CustomerGetByIdResponse, error) {
	response := CustomerGetByIdResponse{}

	////// idを使ってリポジトリからデータを取得
	conn, err := db.GetConnection()
	if err != nil {
		return CustomerGetByIdResponse{}, err
	}
	defer conn.Db.Close()
	entity, err := c.customerRepository.GetById(request.CustomerId, conn)
	if err != nil {
		return CustomerGetByIdResponse{}, err
	}

	///// dtoにつめる
	response.CustomerDto = customer.NewCustomerDtoFromEntity(entity)
	return response, nil
}

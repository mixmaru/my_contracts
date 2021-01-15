package create

import (
	customer_app "github.com/mixmaru/my_contracts/core/application/customer"
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
)

type CustomerCreateInteractor struct {
	customerRepository customer_app.ICustomerRepository
}

func NewCustomerCreateInteractor(customerRepository customer_app.ICustomerRepository) *CustomerCreateInteractor {
	return &CustomerCreateInteractor{customerRepository: customerRepository}
}

func (c CustomerCreateInteractor) Handle(request *CustomerCreateUseCaseRequest) (*CustomerCreateUseCaseResponse, error) {
	response := CustomerCreateUseCaseResponse{}

	// バリデーションする

	// entityをつくる
	newEntity := customer.NewCustomerEntity(request.Name, request.CustomerTypeId, request.Properties)

	// repositoryで保存する
	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()
	tran, err := conn.Begin()
	if err != nil {
		return nil, err
	}
	savedId, err := c.customerRepository.Create(newEntity, tran)
	if err != nil {
		tran.Rollback()
		return nil, err
	}

	// 再読込する
	savedEntity, err := c.customerRepository.GetById(savedId, tran)
	if err != nil {
		tran.Rollback()
		return nil, err
	}
	err = tran.Commit()
	if err != nil {
		tran.Rollback()
		return nil, err
	}

	// dtoに詰める
	response.CustomerDto = customer_app.NewCustomerDtoFromEntity(savedEntity)
	return &response, nil
}

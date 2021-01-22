package create

import (
	customer_app "github.com/mixmaru/my_contracts/core/application/customer"
	"github.com/mixmaru/my_contracts/core/application/customer_type"
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"gopkg.in/gorp.v2"
)

type CustomerCreateInteractor struct {
	customerRepository     customer_app.ICustomerRepository
	customerTypeRepository customer_type.ICustomerTypeRepository
}

func NewCustomerCreateInteractor(customerRepository customer_app.ICustomerRepository, customerTypeRepository customer_type.ICustomerTypeRepository) *CustomerCreateInteractor {
	return &CustomerCreateInteractor{customerRepository: customerRepository, customerTypeRepository: customerTypeRepository}
}

func (c *CustomerCreateInteractor) Handle(request *CustomerCreateUseCaseRequest) (*CustomerCreateUseCaseResponse, error) {
	response := CustomerCreateUseCaseResponse{}

	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()
	tran, err := conn.Begin()
	if err != nil {
		return nil, err
	}

	// バリデーションする
	customerTypeIdValidErrors, err := c.validateCustomerTypeId(request.CustomerTypeId, tran)
	if err != nil {
		return nil, err
	}
	if len(customerTypeIdValidErrors) > 0 {
		response.ValidationErrors = map[string][]string{
			"customer_type_id": customerTypeIdValidErrors,
		}
	}
	if len(response.ValidationErrors) > 0 {
		return &response, nil
	}

	// entityをつくる
	newEntity := customer.NewCustomerEntity(request.Name, request.CustomerTypeId, request.Properties)

	// repositoryで保存する
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

func (c *CustomerCreateInteractor) validateCustomerTypeId(customerTypeId int, executor gorp.SqlExecutor) ([]string, error) {
	////// 存在チェック
	// 取得してみる
	customerTypeEntity, err := c.customerTypeRepository.GetByIdForUpdate(customerTypeId, executor)
	if err != nil {
		return nil, err
	}
	if customerTypeEntity == nil {
		return []string{"存在しないIDです"}, nil
	}
	return nil, nil
}

package create

import "github.com/mixmaru/my_contracts/core/domain/models/customer"

type CustomerCreateInteractor struct {
}

func NewCustomerCreateInteractor() *CustomerCreateInteractor {
	return &CustomerCreateInteractor{}
}

func (c CustomerCreateInteractor) Handle(request *CustomerCreateUseCaseRequest) (*CustomerCreateUseCaseResponse, error) {
	response := CustomerCreateUseCaseResponse{}

	// バリデーションする

	// entityをつくる
	newEntity := customer.NewCustomerEntity(request.Name, request.CustomerTypeId, request.Properties)

	// repositoryで保存する

	// 再読込する

	// dtoに詰める

	return &response, nil
}

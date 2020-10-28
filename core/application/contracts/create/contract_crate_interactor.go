package create

import (
	"github.com/mixmaru/my_contracts/core/application/contracts"
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/application/users"
	"github.com/mixmaru/my_contracts/core/domain/services"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/pkg/errors"
)

type ContractCreateInteractor struct {
	userRepository     users.IUserRepository
	productRepository  products.IProductRepository
	contractRepository contracts.IContractRepository
}

func NewContractCreateInteractor(userRepository users.IUserRepository, productRepository products.IProductRepository, contractRepository contracts.IContractRepository) *ContractCreateInteractor {
	return &ContractCreateInteractor{userRepository: userRepository, productRepository: productRepository, contractRepository: contractRepository}
}

func (c *ContractCreateInteractor) Handle(request *ContractCreateUseCaseRequest) (*ContractCreateUseCaseResponse, error) {
	response := &ContractCreateUseCaseResponse{}

	// トランザクション開始
	conn, err := db.GetConnection()
	if err != nil {
		return response, err
	}
	defer conn.Db.Close()
	tran, err := conn.Begin()
	if err != nil {
		return response, errors.WithStack(err)
	}

	// ドメインサービス作成
	contractDomainService := services.NewContractDomainService(c.userRepository, c.productRepository)
	contractEntity, validationErrors, err := contractDomainService.CreateContract(request.UserId, request.ProductId, request.ContractDateTime, tran)
	if err != nil {
		tran.Rollback()
		return response, err
	}
	if len(validationErrors) > 0 {
		tran.Rollback()
		response.ValidationErrors = validationErrors
		return response, nil
	}
	// 契約保存
	savedContractId, err := c.contractRepository.Create(contractEntity, tran)
	if err != nil {
		tran.Rollback()
		return response, err
	}
	// 再読込
	savedContractEntity, err := c.contractRepository.GetById(savedContractId, tran)
	if err != nil {
		tran.Rollback()
		return response, err
	}
	err = tran.Commit()
	if err != nil {
		return response, errors.Wrapf(err, "コミットに失敗した。request: %v", request)
	}
	// dtoに詰める
	response.ContractDto = contracts.NewContractDtoFromEntity(savedContractEntity)

	return response, nil
}

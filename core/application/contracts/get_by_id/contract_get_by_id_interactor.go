package get_by_id

import (
	"github.com/mixmaru/my_contracts/core/application/contracts"
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/application/users"
	"github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/pkg/errors"
)

type ContractGetByIdInteractor struct {
	contractRepository contracts.IContractRepository
	productRepository  products.IProductRepository
	userRepository     users.IUserRepository
}

func NewContractGetByIdInteractor(contractRepository contracts.IContractRepository, productRepository products.IProductRepository, userRepository users.IUserRepository) *ContractGetByIdInteractor {
	return &ContractGetByIdInteractor{contractRepository: contractRepository, productRepository: productRepository, userRepository: userRepository}
}

func (c ContractGetByIdInteractor) Handle(request *ContractGetByIdUseCaseRequest) (*ContractGetByIdUseCaseResponse, error) {
	response := &ContractGetByIdUseCaseResponse{}

	conn, err := db.GetConnection()
	if err != nil {
		return response, err
	}
	defer conn.Db.Close()

	// リポジトリつかってデータ取得
	contractEntity, err := c.contractRepository.GetById(request.ContractId, conn)
	if err != nil {
		return response, err
	}
	if contractEntity == nil {
		// データがない
		return response, nil
	}
	// 商品データ
	productEntity, err := c.productRepository.GetById(contractEntity.ProductId(), conn)
	if err != nil {
		return response, err
	}
	// ユーザーデータ
	userEntity, err := c.userRepository.GetUserById(contractEntity.UserId(), conn)
	if err != nil {
		return response, err
	}

	// dtoにつめる
	response.ContractDto = contracts.NewContractDtoFromEntity(contractEntity)
	response.ProductDto = products.NewProductDtoFromEntity(productEntity)
	switch userEntity.(type) {
	case *user.UserIndividualEntity:
		response.UserDto = users.NewUserIndividualDtoFromEntity(userEntity.(*user.UserIndividualEntity))
	case *user.UserCorporationEntity:
		response.UserDto = users.NewUserCorporationDtoFromEntity(userEntity.(*user.UserCorporationEntity))
	default:
		return response, errors.Errorf("意図しないUser型が来た。userEntity: %T", userEntity)
	}

	// 返却
	return response, nil
}

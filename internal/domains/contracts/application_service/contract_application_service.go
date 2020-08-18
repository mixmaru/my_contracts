package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/domain_service"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/pkg/errors"
	"time"
)

type ContractApplicationService struct {
	ContractRepository   interfaces.IContractRepository
	UserRepository       interfaces.IUserRepository
	ProductRepository    interfaces.IProductRepository
	RightToUseRepository interfaces.IRightToUseRepository
}

func (c *ContractApplicationService) Register(userId int, productId int, contractDateTime time.Time) (productDto data_transfer_objects.ContractDto, validationErrors map[string][]string, err error) {
	// トランザクション開始
	conn, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, err
	}
	defer conn.Db.Close()
	tran, err := conn.Begin()
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, errors.WithStack(err)
	}

	// ドメインサービス作成
	contractDomainService := domain_service.NewContractDomainService(c.ContractRepository, c.UserRepository, c.ProductRepository, c.RightToUseRepository)
	contractDto, validationErrors, err := contractDomainService.CreateContract(userId, productId, contractDateTime, tran)
	if err != nil {
		tran.Rollback()
		return data_transfer_objects.ContractDto{}, nil, err
	}
	if len(validationErrors) > 0 {
		tran.Rollback()
		return data_transfer_objects.ContractDto{}, validationErrors, nil
	}
	err = tran.Commit()
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, errors.Wrapf(err, "コミットに失敗した。userId: %v, productId: %v, contractDateTime: %v", userId, productId, contractDateTime)
	}

	return contractDto, nil, nil
}

func (c *ContractApplicationService) GetById(id int) (contractDto data_transfer_objects.ContractDto, productDto data_transfer_objects.ProductDto, userDto interface{}, err error) {
	conn, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.ContractDto{}, data_transfer_objects.ProductDto{}, nil, err
	}
	defer conn.Db.Close()

	// リポジトリつかってデータ取得
	contractEntity, productEntity, userEntity, err := c.ContractRepository.GetById(id, conn)
	if err != nil {
		return data_transfer_objects.ContractDto{}, data_transfer_objects.ProductDto{}, nil, err
	}
	if contractEntity == nil {
		// データがない
		return data_transfer_objects.ContractDto{}, data_transfer_objects.ProductDto{}, nil, nil
	}

	// dtoにつめる
	contractDto = data_transfer_objects.NewContractDtoFromEntity(contractEntity)
	productDto = data_transfer_objects.NewProductDtoFromEntity(productEntity)
	switch userEntity.(type) {
	case *entities.UserIndividualEntity:
		userDto = data_transfer_objects.NewUserIndividualDtoFromEntity(userEntity.(*entities.UserIndividualEntity))
	case *entities.UserCorporationEntity:
		userDto = data_transfer_objects.NewUserCorporationDtoFromEntity(userEntity.(*entities.UserCorporationEntity))
	default:
		return data_transfer_objects.ContractDto{}, data_transfer_objects.ProductDto{}, nil, errors.Errorf("意図しないUser型が来た。userEntity: %t", userEntity)
	}

	// 返却
	return contractDto, productDto, userDto, nil
}

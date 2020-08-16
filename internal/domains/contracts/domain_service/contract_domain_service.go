package domain_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/utils"
	"gopkg.in/gorp.v2"
	"time"
)

type ContractDomainService struct {
	contractRepository interfaces.IContractRepository
	userRepository     interfaces.IUserRepository
	productRepository  interfaces.IProductRepository
}

func NewContractDomainService(contractRepository interfaces.IContractRepository, userRepository interfaces.IUserRepository, productRepository interfaces.IProductRepository) *ContractDomainService {
	return &ContractDomainService{
		contractRepository: contractRepository,
		userRepository:     userRepository,
		productRepository:  productRepository,
	}
}

func (c *ContractDomainService) CreateContract(userId, productId int, executeDate time.Time, executor gorp.SqlExecutor) (contractDto data_transfer_objects.ContractDto, validationErrors map[string][]string, err error) {
	// 入力値バリデーション
	validationErrors, err = c.registerValidation(userId, productId, executor)
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, err
	}
	if len(validationErrors) > 0 {
		return data_transfer_objects.ContractDto{}, validationErrors, nil
	}

	// entityを作成
	billingStartDate := c.calculateBillingStartDate(executeDate, 1, utils.CreateJstLocation())

	entity := entities.NewContractEntity(userId, productId, executeDate, billingStartDate)

	// リポジトリで保存
	savedId, err := c.contractRepository.Create(entity, executor)
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, err
	}

	// 再読込
	savedEntity, _, _, err := c.contractRepository.GetById(savedId, executor)
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, err
	}

	// dtoに詰める
	dto := data_transfer_objects.NewContractDtoFromEntity(savedEntity)

	// 返却
	return dto, nil, nil

	//// トランザクション開始
	//conn, err := db_connection.GetConnection()
	//if err != nil {
	//	return data_transfer_objects.ContractDto{}, nil, err
	//}
	//defer conn.Db.Close()
	//tran, err := conn.Begin()
	//if err != nil {
	//	return data_transfer_objects.ContractDto{}, nil, errors.WithStack(err)
	//}
	//
	//// 入力値バリデーション
	//validationErrors, err = c.registerValidation(userId, productId, tran)
	//if err != nil {
	//	return data_transfer_objects.ContractDto{}, nil, err
	//}
	//if len(validationErrors) > 0 {
	//	return data_transfer_objects.ContractDto{}, validationErrors, nil
	//}
	//
	//// entityを作成
	//billingStartDate := c.calculateBillingStartDate(contractDateTime, 1, utils.CreateJstLocation())
	//
	//entity := entities.NewContractEntity(userId, productId, contractDateTime, billingStartDate)
	//
	//// リポジトリで保存
	//savedId, err := c.ContractRepository.Create(entity, tran)
	//if err != nil {
	//	return data_transfer_objects.ContractDto{}, nil, err
	//}
	//
	//// 再読込
	//savedEntity, _, _, err := c.ContractRepository.GetById(savedId, tran)
	//if err != nil {
	//	return data_transfer_objects.ContractDto{}, nil, err
	//}
	//
	//err = tran.Commit()
	//if err != nil {
	//	return data_transfer_objects.ContractDto{}, nil, errors.WithStack(err)
	//}
	//
	//// dtoに詰める
	//dto := data_transfer_objects.NewContractDtoFromEntity(savedEntity)
	//
	//// 返却
	//return dto, nil, nil

	//return nil, nil
}

func (c *ContractDomainService) registerValidation(userId int, productId int, executor gorp.SqlExecutor) (validationErrors map[string][]string, err error) {
	validationErrors = map[string][]string{}

	// userの存在チェック
	user, err := c.userRepository.GetUserById(userId, executor)
	if err != nil {
		return nil, err
	}
	if user == nil {
		validationErrors["user_id"] = []string{
			"存在しません",
		}
	}

	// productの存在チェック
	product, err := c.productRepository.GetById(productId, executor)
	if err != nil {
		return nil, err
	}
	if product == nil {
		validationErrors["product_id"] = []string{
			"存在しません",
		}
	}

	return validationErrors, nil
}

func (c *ContractDomainService) calculateBillingStartDate(contractDate time.Time, freeDays int, locale *time.Location) time.Time {
	addFreeDays := contractDate.AddDate(0, 0, freeDays).In(locale)
	return time.Date(addFreeDays.Year(), addFreeDays.Month(), addFreeDays.Day(), 0, 0, 0, 0, locale)
}

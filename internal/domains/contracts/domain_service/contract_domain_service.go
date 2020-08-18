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
	contractRepository   interfaces.IContractRepository
	userRepository       interfaces.IUserRepository
	productRepository    interfaces.IProductRepository
	rightToUseRepository interfaces.IRightToUseRepository
}

func NewContractDomainService(
	contractRepository interfaces.IContractRepository,
	userRepository interfaces.IUserRepository,
	productRepository interfaces.IProductRepository,
	rightToUseRepository interfaces.IRightToUseRepository,
) *ContractDomainService {
	return &ContractDomainService{
		contractRepository:   contractRepository,
		userRepository:       userRepository,
		productRepository:    productRepository,
		rightToUseRepository: rightToUseRepository,
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

	// 課金開始日算出
	billingStartDate := c.calculateBillingStartDate(executeDate, 1, utils.CreateJstLocation())

	// 契約の作成
	savedContractId, err := c.createNewContract(userId, productId, executeDate, billingStartDate, executor)
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, err
	}

	// 使用権の作成
	validTo := billingStartDate.AddDate(0, 1, 0)
	_, err = c.createNewRightToUse(savedContractId, executeDate, validTo, executor)

	// 再読込
	savedEntity, _, _, err := c.contractRepository.GetById(savedContractId, executor)
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, err
	}

	// dtoに詰める
	dto := data_transfer_objects.NewContractDtoFromEntity(savedEntity)

	// 返却
	return dto, nil, nil
}

func (c *ContractDomainService) createNewContract(userId, productId int, executeDate, billingStartDate time.Time, executor gorp.SqlExecutor) (savedContractId int, err error) {
	// 契約entityを作成
	entity := entities.NewContractEntity(userId, productId, executeDate, billingStartDate)

	// リポジトリで保存
	savedContractId, err = c.contractRepository.Create(entity, executor)
	if err != nil {
		return 0, err
	}
	return savedContractId, nil
}

func (c *ContractDomainService) createNewRightToUse(contractId int, validFrom, validTo time.Time, executor gorp.SqlExecutor) (savedRightToUseId int, err error) {
	// 使用権entityを作成
	rightToUseEntity := entities.NewRightToUseEntity(contractId, validFrom, validTo)

	// リポジトリで保存
	savedRightToUseId, err = c.rightToUseRepository.Create(rightToUseEntity, executor)
	if err != nil {
		return 0, err
	}
	return savedRightToUseId, nil
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

package domain_service

import (
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
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
	billingStartDate := c.calculateBillingStartDate(executeDate, 0, utils.CreateJstLocation())

	// 契約の作成
	savedContractId, err := c.createNewContract(userId, productId, executeDate, billingStartDate, executor)
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, err
	}

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
	// 使用権entityを作成
	jst := utils.CreateJstLocation()
	executeDateJst := executeDate.In(jst)
	validTo := utils.CreateJstTime(executeDateJst.Year(), executeDateJst.Month()+1, executeDateJst.Day(), 0, 0, 0, 0)
	rightToUseEntity := entities.NewRightToUseEntity(executeDate, validTo)

	// 契約entityを作成
	entity := entities.NewContractEntity(userId, productId, executeDate, billingStartDate, []*entities.RightToUseEntity{rightToUseEntity})

	// リポジトリで保存
	savedContractId, err = c.contractRepository.Create(entity, executor)
	if err != nil {
		return 0, err
	}
	return savedContractId, nil
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

/*
渡された使用権の次の期間の使用権を作成して永続化する。永続化したデータを返却する
*/
func (c *ContractDomainService) CreateNextTermRightToUse(currentRightToUse *entities.RightToUseEntity, executor gorp.SqlExecutor) (*entities.RightToUseEntity, error) {
	jst := utils.CreateJstLocation()

	// 商品集約を取得する
	product, err := c.productRepository.GetByRightToUseId(currentRightToUse.Id(), executor)
	if err != nil {
		return nil, err
	}

	// 商品の期間（年、月、カスタム期間）を取得する
	termType, err := product.GetTermType()
	if err != nil {
		return nil, err
	}

	// 期間から次の期間を算出する
	from := currentRightToUse.ValidTo().In(jst)
	var to time.Time
	switch termType {
	case entities.TermMonthly:
		to = from.AddDate(0, 1, 0)
	case entities.TermYearly:
		to = from.AddDate(1, 0, 0)
	case entities.TermCustom:
		term, err := product.GetCustomTerm()
		if err != nil {
			return nil, err
		}
		to = from.AddDate(0, 0, term)
	case entities.TermLump:
		return nil, errors.Errorf("一括購入タイプの商品は継続処理できません。termType: %v, product: %+v", termType, product)
	default:
		return nil, errors.Errorf("考慮外のtermType。termType: %v, product: %+v", termType, product)
	}

	// entityを作る
	nextTermEntity := entities.NewRightToUseEntity(currentRightToUse.ContractId(), from, to)

	// 保存する。
	nextTermEntityId, err := c.rightToUseRepository.Create(nextTermEntity, executor)
	if err != nil {
		return nil, err
	}

	// 再読込
	nextTermEntity, err = c.rightToUseRepository.GetById(nextTermEntityId, executor)
	if err != nil {
		return nil, err
	}

	return nextTermEntity, nil
}

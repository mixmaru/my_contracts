package services

import (
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/application/users"
	"github.com/mixmaru/my_contracts/core/domain/models/contract"
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"time"
)

type ContractDomainService struct {
	userRepository    users.IUserRepository
	productRepository products.IProductRepository
}

func NewContractDomainService(
	userRepository users.IUserRepository,
	productRepository products.IProductRepository,
) *ContractDomainService {
	return &ContractDomainService{
		userRepository:    userRepository,
		productRepository: productRepository,
	}
}

func (c *ContractDomainService) CreateContract(userId, productId int, executeDate time.Time, executor gorp.SqlExecutor) (contract *contract.ContractEntity, validationErrors map[string][]string, err error) {
	// 入力値バリデーション
	validationErrors, err = c.registerValidation(userId, productId, executor)
	if err != nil {
		return nil, nil, err
	}
	if len(validationErrors) > 0 {
		return nil, validationErrors, nil
	}

	// 課金開始日算出
	billingStartDate := c.calculateBillingStartDate(executeDate, 0, utils.CreateJstLocation())

	// 契約の作成
	contract, err = c.createNewContract(userId, productId, executeDate, billingStartDate)
	if err != nil {
		return nil, nil, err
	}
	return contract, nil, nil
}

func (c *ContractDomainService) createNewContract(userId, productId int, executeDate, billingStartDate time.Time) (*contract.ContractEntity, error) {
	// 使用権entityを作成
	jst := utils.CreateJstLocation()
	executeDateJst := executeDate.In(jst)
	validTo := utils.CreateJstTime(executeDateJst.Year(), executeDateJst.Month()+1, executeDateJst.Day(), 0, 0, 0, 0)
	rightToUseEntity := contract.NewRightToUseEntity(executeDate, validTo)

	// 契約entityを作成
	entity := contract.NewContractEntity(userId, productId, executeDate, billingStartDate, []*contract.RightToUseEntity{rightToUseEntity})

	return entity, nil
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
渡された契約集約の使用権を元に、次の期間の使用権を作成して返す（永続化はしない）
*/
func CreateNextTermRightToUse(currentRightToUse *contract.RightToUseEntity, productEntity *product.ProductEntity) (*contract.RightToUseEntity, error) {
	// 商品の期間（年、月、カスタム期間）を取得する
	termType, err := productEntity.GetTermType()
	if err != nil {
		return nil, err
	}

	// 期間から次の期間を算出する
	jst := utils.CreateJstLocation()
	from := currentRightToUse.ValidTo().In(jst)
	var to time.Time
	switch termType {
	case product.TermMonthly:
		to = from.AddDate(0, 1, 0)
	case product.TermYearly:
		to = from.AddDate(1, 0, 0)
	case product.TermCustom:
		term, err := productEntity.GetCustomTerm()
		if err != nil {
			return nil, err
		}
		to = from.AddDate(0, 0, term)
	case product.TermLump:
		return nil, errors.Errorf("一括購入タイプの商品は継続処理できません。termType: %v, productEntity: %+v", termType, productEntity)
	default:
		return nil, errors.Errorf("考慮外のtermType。termType: %v, productEntity: %+v", termType, productEntity)
	}

	// entityを作る
	nextTermEntity := contract.NewRightToUseEntity(from, to)
	return nextTermEntity, err
}

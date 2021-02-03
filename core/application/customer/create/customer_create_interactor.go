package create

import (
	"fmt"
	"strconv"
	"strings"

	customer_app "github.com/mixmaru/my_contracts/core/application/customer"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/application/customer_type"
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"gopkg.in/gorp.v2"
)

type CustomerCreateInteractor struct {
	customerRepository             customer_app.ICustomerRepository
	customerTypeRepository         customer_type.ICustomerTypeRepository
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository
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
	response.ValidationErrors, err = c.validation(request, tran)
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

func (c *CustomerCreateInteractor) validation(request *CustomerCreateUseCaseRequest, executor gorp.SqlExecutor) (map[string][]string, error) {
	validationErrors := map[string][]string{}

	// 取得してみる
	customerTypeEntity, err := c.customerTypeRepository.GetByIdForUpdate(request.CustomerTypeId, executor)
	if err != nil {
		return nil, err
	}
	if customerTypeEntity == nil {
		validationErrors["customer_type_id"] = []string{"存在しないIDです"}
	}

	// customer_type_idのバリデーションエラーがない場合のみチェック
	if len(validationErrors["customer_type_id"]) == 0 {
		// propertiesバリデーション
		// リクエストされたpropertyTypeIdにcustomerTypeに存在しないIDがないかどうかチェック
		badPropertyTypeIds := []string{}
        REQUEST_PROPERTY_ROOP:
		for propertyTypeId := range request.Properties {
			for _, existedId := range customerTypeEntity.CustomerPropertyTypeIds() {
				if propertyTypeId == existedId {
					continue REQUEST_PROPERTY_ROOP
				}
			}
			badPropertyTypeIds = append(badPropertyTypeIds, strconv.Itoa(propertyTypeId))
		}
		// あればバリデーションエラー
		if len(badPropertyTypeIds) > 0 {
			validationErrors["customer_property_type"] = []string{fmt.Sprintf("%vは存在しないIDです", strings.Join(badPropertyTypeIds, ", "))}
		}
	}

	return validationErrors, nil
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

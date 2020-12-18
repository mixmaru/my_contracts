package create

import (
	"fmt"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/application/customer_type"
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"strconv"
	"strings"
)

type CustomerTypeCreateInteractor struct {
	customerTypeRepository         customer_type.ICustomerTypeRepository
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository
}

func NewCustomerTypeCreateInteractor(
	customerTypeRepository customer_type.ICustomerTypeRepository,
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository,
) *CustomerTypeCreateInteractor {
	return &CustomerTypeCreateInteractor{customerTypeRepository: customerTypeRepository, customerPropertyTypeRepository: customerPropertyTypeRepository}
}

const MAX_RETRY_NUM int = 2

func (c CustomerTypeCreateInteractor) Handle(request *CustomerTypeCreateUseCaseRequest) (*CustomerTypeCreateUseCaseResponse, error) {
	response := &CustomerTypeCreateUseCaseResponse{}

	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	// 事前にnameの重複チェックを行っても別トランザクションから先に挿入されると重複する可能性がある
	// nameの重複が起こったら1回まで再実行する
	execCount := 0
	var reloadedEntity *customer.CustomerTypeEntity
	var reloadedPropertyEntities []*customer.CustomerPropertyTypeEntity
	for {
		execCount++
		if execCount > MAX_RETRY_NUM {
			return nil, errors.Errorf("実行回数がMAX_RETRY_NUMを超えました。execCount: %v, MAX_RETRY_NUM: %v, request: %+v", execCount, MAX_RETRY_NUM, request)
		}

		tran, err := conn.Begin()
		if err != nil {
			return nil, err
		}

		// バリデーション
		response.ValidationError, err = c.validation(request, tran)
		if err != nil {
			return nil, err
		}
		if len(response.ValidationError) > 0 {
			return response, nil
		}

		// Entity用意
		entity := customer.NewCustomerTypeEntity(request.Name, request.CustomerPropertyTypeIds)

		// カスタマータイプ保存
		savedId, err := c.customerTypeRepository.Create(entity, tran)
		if err != nil {
			tran.Rollback()
			if execCount < MAX_RETRY_NUM {
				continue
			}
			return nil, err
		}

		// リロード
		reloadedEntity, err = c.customerTypeRepository.GetById(savedId, tran)
		if err != nil {
			tran.Rollback()
			return nil, err
		}
		reloadedPropertyEntities, err = c.customerPropertyTypeRepository.GetByIds(reloadedEntity.CustomerPropertyTypeIds(), tran)
		if err != nil {
			tran.Rollback()
			return nil, err
		}

		// コミット
		if err := tran.Commit(); err != nil {
			tran.Rollback()
			return nil, errors.Wrapf(err, "トランザクションコミットに失敗しました。request: %+v", request)
		}
		break
	}

	// 返却dtoに詰める
	response.CustomerTypeDto, err = customer_type.NewCustomerTypeDtoFromEntity(reloadedEntity, reloadedPropertyEntities)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *CustomerTypeCreateInteractor) validation(request *CustomerTypeCreateUseCaseRequest, executor gorp.SqlExecutor) (map[string][]string, error) {
	retMap := map[string][]string{}

	// name重複チェック
	errorMessage, err := c.duplicatedNameValidation(request.Name, executor)
	if err != nil {
		return retMap, err
	}
	if errorMessage != "" {
		retMap["name"] = []string{errorMessage}
	}

	// CustomerPropertyTypeIdsの存在チェック
	errorMessage, err = c.existCustomerPropertyTypeIdsValidation(request.CustomerPropertyTypeIds, executor)
	if err != nil {
		return retMap, err
	}
	if errorMessage != "" {
		retMap["customer_property_ids"] = []string{errorMessage}
	}

	return retMap, err
}

func (c *CustomerTypeCreateInteractor) duplicatedNameValidation(name string, executor gorp.SqlExecutor) (string, error) {
	// name重複チェック
	entity, err := c.customerTypeRepository.GetByName(name, executor)
	if err != nil {
		return "", err
	}
	if entity != nil {
		return "既に存在する名前です", nil
	} else {
		return "", nil
	}
}

func (c CustomerTypeCreateInteractor) existCustomerPropertyTypeIdsValidation(ids []int, executor gorp.SqlExecutor) (string, error) {
	entities, err := c.customerPropertyTypeRepository.GetByIds(ids, executor)
	if err != nil {
		return "", err
	}

	notExistIds := []string{}
OUTER:
	for _, requestId := range ids {
		for _, entity := range entities {
			if requestId == entity.Id() {
				continue OUTER
			}
		}
		notExistIds = append(notExistIds, strconv.Itoa(requestId))
	}
	if len(notExistIds) > 0 {
		return fmt.Sprintf("%s は存在しないidです", strings.Join(notExistIds, ", ")), nil
	} else {
		return "", nil
	}
}

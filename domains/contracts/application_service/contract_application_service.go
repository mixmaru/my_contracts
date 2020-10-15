package application_service

import (
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/domains/contracts/domain_service"
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/db_connection"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"time"
)

type ContractApplicationService struct {
	contractRepository interfaces.IContractRepository
	userRepository     interfaces.IUserRepository
	productRepository  interfaces.IProductRepository
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
	contractDomainService := domain_service.NewContractDomainService(c.userRepository, c.productRepository)
	contractEntity, validationErrors, err := contractDomainService.CreateContract(userId, productId, contractDateTime, tran)
	if err != nil {
		tran.Rollback()
		return data_transfer_objects.ContractDto{}, nil, err
	}
	if len(validationErrors) > 0 {
		tran.Rollback()
		return data_transfer_objects.ContractDto{}, validationErrors, nil
	}
	// 契約保存
	savedContractId, err := c.contractRepository.Create(contractEntity, tran)
	if err != nil {
		tran.Rollback()
		return data_transfer_objects.ContractDto{}, nil, err
	}
	// 再読込
	savedContractEntity, _, _, err := c.contractRepository.GetById(savedContractId, tran)
	if err != nil {
		tran.Rollback()
		return data_transfer_objects.ContractDto{}, nil, err
	}
	err = tran.Commit()
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, errors.Wrapf(err, "コミットに失敗した。userId: %v, productId: %v, contractDateTime: %v", userId, productId, contractDateTime)
	}
	// dtoに詰める
	contractDto := data_transfer_objects.NewContractDtoFromEntity(savedContractEntity)

	return contractDto, nil, nil
}

func (c *ContractApplicationService) GetById(id int) (contractDto data_transfer_objects.ContractDto, productDto data_transfer_objects.ProductDto, userDto interface{}, err error) {
	conn, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.ContractDto{}, data_transfer_objects.ProductDto{}, nil, err
	}
	defer conn.Db.Close()

	// リポジトリつかってデータ取得
	contractEntity, productEntity, userEntity, err := c.contractRepository.GetById(id, conn)
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

/*
渡した実行日から5日以内に期間終了である使用権に対して、次の期間の使用権データを作成して永続化して返却する
*/
func (c *ContractApplicationService) CreateNextRightToUse(executeDate time.Time) (nextTermContracts []data_transfer_objects.ContractDto, err error) {
	db, err := db_connection.GetConnection()
	if err != nil {
		return nil, err
	}
	defer db.Db.Close()
	// 使用権更新の対象契約を取得
	contracts, err := c.contractRepository.GetRecurTargets(executeDate, db)
	if err != nil {
		return nil, err
	}

	nextTermContracts = make([]data_transfer_objects.ContractDto, 0, len(contracts))
	// 次の使用権を作成して更新する
	for _, contract := range contracts {
		if len(contract.RightToUses()) >= 2 {
			return nil, errors.Errorf("使用権が2つ以上ある（既に次期使用権がある可能性がある） contract: %+v", contract)
		}

		tran, err := db.Begin()
		if err != nil {
			return nil, errors.Wrap(err, "トランザクション開始失敗")
		}

		product, err := c.productRepository.GetById(contract.ProductId(), tran)
		if err != nil {
			return nil, err
		}
		nextTermRightToUse, err := domain_service.CreateNextTermRightToUse(contract.RightToUses()[0], product)
		if err != nil {
			return nil, err
		}
		contract.AddNextTermRightToUses(nextTermRightToUse)
		// contractの保存実行
		err = c.contractRepository.Update(contract, tran)
		if err != nil {
			return nil, err
		}
		// リロード
		reloadedContract, _, _, err := c.contractRepository.GetById(contract.Id(), tran)
		if err != nil {
			return nil, err
		}
		err = tran.Commit()
		if err != nil {
			return nil, errors.Wrapf(err, "コミットに失敗しました")
		}
		nextTermContracts = append(nextTermContracts, data_transfer_objects.NewContractDtoFromEntity(reloadedContract))
	}
	return nextTermContracts, nil
}

/*
渡した基準日に期限が切れている使用権をアーカイブ処理し、処理した使用権dtoを返す
*/
func (c *ContractApplicationService) ArchiveExpiredRightToUse(baseDate time.Time) (archivedRightToUse []data_transfer_objects.RightToUseDto, err error) {
	db, err := db_connection.GetConnection()
	if err != nil {
		return nil, err
	}
	defer db.Db.Close()

	retDto := []data_transfer_objects.RightToUseDto{}

	// 対象取得
	targetContractIds, err := c.contractRepository.GetHavingExpiredRightToUseContractIds(baseDate, db)
	if err != nil {
		return retDto, err
	}

	// アーカイブ処理（同時実行時の更新失敗（別トランザクションがアーカイブした等）に備えて3回までリトライする）
	for _, contractId := range targetContractIds {
		count := 0
		for {
			count++
			dtos, err := c.execArchive(contractId, baseDate, db)
			retDto = append(retDto, dtos...)
			if err != nil {
				if count < 3 {
					continue
				} else {
					return retDto, err
				}
			} else {
				break
			}
		}
	}

	return retDto, nil
}

func (c *ContractApplicationService) execArchive(contractId int, baseDate time.Time, executor *gorp.DbMap) ([]data_transfer_objects.RightToUseDto, error) {
	retDtos := []data_transfer_objects.RightToUseDto{}

	// トランザクション開始
	tran, err := executor.Begin()
	if err != nil {
		return retDtos, errors.Wrapf(err, "トランザクション開始失敗")
	}
	_, err = tran.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	if err != nil {
		return retDtos, errors.Wrapf(err, "トランザクション分離レベル切り替え失敗。")
	}

	// データ取得
	contractEntity, _, _, err := c.contractRepository.GetById(contractId, tran)
	if err != nil {
		return retDtos, err
	}
	// アーカイブ処理
	contractEntity.ArchiveRightToUseByValidTo(baseDate)
	err = c.contractRepository.Update(contractEntity, executor)
	if err != nil {
		return retDtos, err
	}
	// 返却dtoを用意
	for _, entity := range contractEntity.GetToArchiveRightToUses() {
		retDtos = append(retDtos, data_transfer_objects.NewRightToUseDtoFromEntity(entity))
	}
	// コミット
	err = tran.Commit()
	if err != nil {
		return retDtos, errors.Wrapf(err, "コミット失敗")
	}

	return retDtos, nil
}

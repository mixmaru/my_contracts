package archive_expired_right_to_use

import (
	"github.com/mixmaru/my_contracts/core/application/contracts"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"time"
)

type ContractArchiveExpiredRightToUseInteractor struct {
	contractRepository contracts.IContractRepository
}

func NewContractArchiveExpiredRightToUseInteractor(contractRepository contracts.IContractRepository) *ContractArchiveExpiredRightToUseInteractor {
	return &ContractArchiveExpiredRightToUseInteractor{contractRepository: contractRepository}
}

func (c *ContractArchiveExpiredRightToUseInteractor) Handle(request *ContractArchiveExpiredRightToUseUseCaseRequest) (*ContractArchiveExpiredRightToUseUseCaseResponse, error) {
	response := &ContractArchiveExpiredRightToUseUseCaseResponse{}

	db, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer db.Db.Close()

	response.ArchivedRightToUse = []contracts.RightToUseDto{}

	// 対象取得
	targetContractIds, err := c.contractRepository.GetHavingExpiredRightToUseContractIds(request.BaseDate, db)
	if err != nil {
		return response, err
	}

	// アーカイブ処理（同時実行時の更新失敗（別トランザクションがアーカイブした等）に備えて3回までリトライする）
	for _, contractId := range targetContractIds {
		count := 0
		for {
			count++

			// トランザクション開始
			tran, err := db.Begin()
			if err != nil {
				return response, errors.Wrapf(err, "トランザクション開始失敗")
			}
			_, err = tran.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
			if err != nil {
				return response, errors.Wrapf(err, "トランザクション分離レベル切り替え失敗。")
			}

			// アーカイブ実行
			dtos, err := c.execArchive(contractId, request.BaseDate, tran)
			response.ArchivedRightToUse = append(response.ArchivedRightToUse, dtos...)
			if err != nil {
				// 失敗してたら3回までリトライ
				tran.Rollback()
				if count < 3 {
					continue
				} else {
					// 3回以上だったらエラーにする
					return response, err
				}
			}

			//コミット
			err = tran.Commit()
			if err != nil {
				return response, errors.Wrapf(err, "コミット失敗")
			}
			break
		}
	}

	return response, nil
}

func (c *ContractArchiveExpiredRightToUseInteractor) execArchive(contractId int, baseDate time.Time, executor gorp.SqlExecutor) ([]contracts.RightToUseDto, error) {
	retDtos := []contracts.RightToUseDto{}

	// データ取得
	contractEntity, err := c.contractRepository.GetById(contractId, executor)
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
		retDtos = append(retDtos, contracts.NewRightToUseDtoFromEntity(entity))
	}

	return retDtos, nil
}

package create_next_right_to_use

import (
	"github.com/mixmaru/my_contracts/core/application/contracts"
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/domain/services"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/pkg/errors"
)

type ContractCreateNextRightToUseInteractor struct {
	contractRepository contracts.IContractRepository
	productRepository  products.IProductRepository
}

func NewContractCreateNextRightToUseInteractor(contractRepository contracts.IContractRepository, productRepository products.IProductRepository) *ContractCreateNextRightToUseInteractor {
	return &ContractCreateNextRightToUseInteractor{contractRepository: contractRepository, productRepository: productRepository}
}

/*
渡した実行日から5日以内に期間終了である使用権に対して、次の期間の使用権データを作成して永続化して返却する
*/
func (c ContractCreateNextRightToUseInteractor) Handle(request *ContractCreateNextRightToUseUseCaseRequest) (*ContractCreateNextRightToUseUseCaseResponse, error) {
	response := &ContractCreateNextRightToUseUseCaseResponse{}

	db, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer db.Db.Close()
	// 使用権更新の対象契約を取得
	contractEntities, err := c.contractRepository.GetRecurTargets(request.ExecuteDate, db)
	if err != nil {
		return nil, err
	}

	response.NextTermContracts = make([]contracts.ContractDto, 0, len(contractEntities))
	// 次の使用権を作成して更新する
	for _, contractEntity := range contractEntities {
		if len(contractEntity.RightToUses()) >= 2 {
			return nil, errors.Errorf("使用権が2つ以上ある（既に次期使用権がある可能性がある） contractEntity: %+v", contractEntity)
		}

		tran, err := db.Begin()
		if err != nil {
			return nil, errors.Wrap(err, "トランザクション開始失敗")
		}

		product, err := c.productRepository.GetById(contractEntity.ProductId(), tran)
		if err != nil {
			return nil, err
		}
		nextTermRightToUse, err := services.CreateNextTermRightToUse(contractEntity.RightToUses()[0], product)
		if err != nil {
			return nil, err
		}
		contractEntity.AddNextTermRightToUses(nextTermRightToUse)
		// contractEntityの保存実行
		err = c.contractRepository.Update(contractEntity, tran)
		if err != nil {
			return nil, err
		}
		// リロード
		reloadedContract, err := c.contractRepository.GetById(contractEntity.Id(), tran)
		if err != nil {
			return nil, err
		}
		err = tran.Commit()
		if err != nil {
			return nil, errors.Wrapf(err, "コミットに失敗しました")
		}
		//nextTermContracts = append(nextTermContracts, data_transfer_objects.NewContractDtoFromEntity(reloadedContract))
		response.NextTermContracts = append(response.NextTermContracts, contracts.NewContractDtoFromEntity(reloadedContract))
	}
	return response, nil
}

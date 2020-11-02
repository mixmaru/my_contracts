package archive_expired_right_to_use

import (
	"github.com/golang/mock/gomock"
	contract_create "github.com/mixmaru/my_contracts/core/application/contracts/create"
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/application/products/create"
	"github.com/mixmaru/my_contracts/core/application/users"
	user_create "github.com/mixmaru/my_contracts/core/application/users/create"
	"github.com/mixmaru/my_contracts/core/domain/models/contract"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/core/infrastructure/mock/mock_contract"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
	"time"
)

func TestContractArchiveExpiredRightToUseInteractor_Handle(t *testing.T) {
	userRep := db.NewUserRepository()
	productRep := db.NewProductRepository()
	contractRep := db.NewContractRepository()

	t.Run("渡した基準日に期限が切れている使用権をアーカイブ処理し、処理した使用権dtoを返す", func(t *testing.T) {
		////// 準備
		// 事前に存在するデータを削除しておく
		conn, err := db.GetConnection()
		assert.NoError(t, err)
		deleteSql := `
DELETE FROM discount_apply_contract_updates;
DELETE FROM bill_details;
DELETE FROM right_to_use_active;
DELETE FROM right_to_use_history;
DELETE FROM right_to_use;
DELETE FROM contracts;
`
		_, err = conn.Exec(deleteSql)
		assert.NoError(t, err)

		user := createUser()
		product := createProduct()
		contractCreateIntractor := contract_create.NewContractCreateInteractor(userRep, productRep, contractRep)
		contractCreateResponse1, err := contractCreateIntractor.Handle(contract_create.NewContractCreateUseCaseRequest(
			user.Id,
			product.Id,
			utils.CreateJstTime(2020, 5, 1, 3, 0, 0, 0)),
		)
		if err != nil || len(contractCreateResponse1.ValidationErrors) > 0 {
			panic("データ作成失敗")
		}

		contractCreateResponse2, err := contractCreateIntractor.Handle(contract_create.NewContractCreateUseCaseRequest(
			user.Id,
			product.Id,
			utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0)),
		)
		if err != nil || len(contractCreateResponse2.ValidationErrors) > 0 {
			panic("データ作成失敗")
		}

		contractCreateResponse3, err := contractCreateIntractor.Handle(contract_create.NewContractCreateUseCaseRequest(
			user.Id,
			product.Id,
			utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)),
		)
		if err != nil || len(contractCreateResponse3.ValidationErrors) > 0 {
			panic("データ作成失敗")
		}

		////// 実行
		interactor := NewContractArchiveExpiredRightToUseInteractor(contractRep)
		response, err := interactor.Handle(NewContractArchiveExpiredRightToUseUseCaseRequest(utils.CreateJstTime(2020, 7, 2, 0, 0, 0, 0)))
		assert.NoError(t, err)

		////// 検証
		dtos := response.ArchivedRightToUses
		assert.Len(t, dtos, 2)
		assert.Equal(t, dtos[0], contractCreateResponse1.ContractDto.RightToUseDtos[0])
		assert.Equal(t, dtos[1], contractCreateResponse2.ContractDto.RightToUseDtos[0])
	})

	t.Run("同時実行テスト", func(t *testing.T) {
		t.Run("別トランザクションが先に同じデータを処理して失敗した場合再トライしてスキップされる", func(t *testing.T) {
			////// 準備
			// 事前に存在するデータを削除しておく
			conn, err := db.GetConnection()
			assert.NoError(t, err)
			deleteSql := `
			DELETE FROM discount_apply_contract_updates;
			DELETE FROM bill_details;
			DELETE FROM right_to_use_active;
			DELETE FROM right_to_use_history;
			DELETE FROM right_to_use;
			DELETE FROM contracts;
			`
			_, err = conn.Exec(deleteSql)
			assert.NoError(t, err)

			user := createUser()
			product := createProduct()
			contractCrateInteractor := contract_create.NewContractCreateInteractor(userRep, productRep, contractRep)
			contractCrateResponse1, err := contractCrateInteractor.Handle(contract_create.NewContractCreateUseCaseRequest(
				user.Id,
				product.Id,
				utils.CreateJstTime(2020, 5, 1, 3, 0, 0, 0)),
			)
			if err != nil || len(contractCrateResponse1.ValidationErrors) > 0 {
				panic("データ作成失敗")
			}

			contractCrateResponse2, err := contractCrateInteractor.Handle(contract_create.NewContractCreateUseCaseRequest(
				user.Id,
				product.Id,
				utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0)),
			)
			if err != nil || len(contractCrateResponse2.ValidationErrors) > 0 {
				panic("データ作成失敗")
			}

			contractCrateResponse3, err := contractCrateInteractor.Handle(contract_create.NewContractCreateUseCaseRequest(
				user.Id,
				product.Id,
				utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)),
			)
			if err != nil || len(contractCrateResponse3.ValidationErrors) > 0 {
				panic("データ作成失敗")
			}

			// モックリポジトリ
			contractRep := db.NewContractRepository()
			ctrl := gomock.NewController(t)
			contractRepMock := mock_contract.NewMockIContractRepository(ctrl)
			contractRepMock.EXPECT().GetHavingExpiredRightToUseContractIds(gomock.Any(), gomock.Any()).DoAndReturn(
				func(baseDate time.Time, executor gorp.SqlExecutor) ([]int, error) {
					return contractRep.GetHavingExpiredRightToUseContractIds(baseDate, executor)
				}).AnyTimes()
			contractRepMock.EXPECT().GetById(gomock.Any(), gomock.Any()).DoAndReturn(
				func(id int, executor gorp.SqlExecutor) (contract *contract.ContractEntity, err error) {
					return contractRep.GetById(id, executor)
				}).AnyTimes()
			count := 0
			contractRepMock.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
				func(contractEntity *contract.ContractEntity, executor gorp.SqlExecutor) error {
					count++
					if count == 2 {
						// ２回目はエラーを返す（別トランザクションが既に更新をかけていたという想定）
						conn, err := db.GetConnection()
						if err != nil {
							return errors.New("別トランザクション用db接続に失敗")
						}
						err = contractRep.Update(contractEntity, conn) // 別トランザクションが更新をかけていた事を再現する
						if err != nil {
							return errors.New("更新に失敗した。")
						}
						return errors.New("先にやられた")
					} else {
						return contractRep.Update(contractEntity, executor)
					}
				}).AnyTimes()

			////// 実行
			interactor := NewContractArchiveExpiredRightToUseInteractor(contractRepMock)
			response, err := interactor.Handle(NewContractArchiveExpiredRightToUseUseCaseRequest(utils.CreateJstTime(2020, 7, 2, 0, 0, 0, 0)))
			assert.NoError(t, err)

			////// 検証
			assert.Len(t, response.ArchivedRightToUses, 1)
			assert.Equal(t, response.ArchivedRightToUses[0], contractCrateResponse1.ContractDto.RightToUseDtos[0])
		})
	})
}

func createProduct() products.ProductDto {
	interactor := create.NewProductCreateInteractor(db.NewProductRepository())
	response, err := interactor.Handle(create.NewProductCreateUseCaseRequest("商品", "2000"))
	if err != nil || len(response.ValidationError) > 0 {
		panic("データ作成失敗")
	}

	return response.ProductDto
}

func createUser() users.UserIndividualDto {
	interactor := user_create.NewUserIndividualCreateInteractor(db.NewUserRepository())
	response, err := interactor.Handle(user_create.NewUserIndividualCreateUseCaseRequest("個人たろう"))
	if err != nil || len(response.ValidationErrors) > 0 {
		panic("データ作成失敗")
	}
	return response.UserDto
}

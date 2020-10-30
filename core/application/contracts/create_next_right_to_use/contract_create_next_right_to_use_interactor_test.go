package create_next_right_to_use

import (
	create3 "github.com/mixmaru/my_contracts/core/application/contracts/create"
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/application/products/create"
	"github.com/mixmaru/my_contracts/core/application/users"
	create2 "github.com/mixmaru/my_contracts/core/application/users/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContractCreateNextRightToUseInteractor_Handle(t *testing.T) {
	t.Run("渡した実行日から5日以内に期間終了である使用権に対して、次の期間の使用権データを作成して永続化して返却する", func(t *testing.T) {
		// 事前に影響のあるデータを削除しておく（ちょっと広めに削除）
		conn, err := db.GetConnection()
		assert.NoError(t, err)
		defer conn.Db.Close()
		_, err = conn.Exec("DELETE FROM right_to_use_active WHERE right_to_use_id IN (SELECT id FROM right_to_use WHERE '2020-05-25' <= valid_to AND valid_to <= '2020-06-02')")
		assert.NoError(t, err)
		_, err = conn.Exec("DELETE FROM right_to_use WHERE '2020-05-25' <= valid_to AND valid_to <= '2020-06-02'")
		assert.NoError(t, err)

		////// 準備（2020-05-31が終了日である使用権と2020-05-29が終了日である使用権を作成する）
		user := createUser()
		product := createProduct()
		createInteractor := create3.NewContractCreateInteractor(db.NewUserRepository(), db.NewProductRepository(), db.NewContractRepository())
		response, err := createInteractor.Handle(create3.NewContractCreateUseCaseRequest(user.Id, product.Id, utils.CreateJstTime(2020, 5, 1, 3, 0, 0, 0)))
		if err != nil || len(response.ValidationErrors) > 0 {
			panic("データ作成失敗")
		}
		response, err = createInteractor.Handle(create3.NewContractCreateUseCaseRequest(user.Id, product.Id, utils.CreateJstTime(2020, 4, 30, 0, 0, 0, 0)))
		if err != nil || len(response.ValidationErrors) > 0 {
			panic("データ作成失敗")
		}

		////// 実行
		nextInteractor := NewContractCreateNextRightToUseInteractor(db.NewContractRepository(), db.NewProductRepository())
		actualresponse, err := nextInteractor.Handle(NewContractCreateNextRightToUseUseCaseRequest(utils.CreateJstTime(2020, 5, 28, 0, 10, 0, 0)))
		assert.NoError(t, err)

		////// 検証
		nextTermContracts := actualresponse.NextTermContracts
		assert.Len(t, nextTermContracts, 2)
		// 1つめ
		recurRightToUse1 := nextTermContracts[0].RightToUseDtos[1]
		assert.NotZero(t, recurRightToUse1.Id)
		assert.NotZero(t, recurRightToUse1.CreatedAt)
		assert.NotZero(t, recurRightToUse1.UpdatedAt)
		assert.True(t, recurRightToUse1.ValidFrom.Equal(utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0)))
		assert.True(t, recurRightToUse1.ValidTo.Equal(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
		// 2つめ
		recurRightToUse2 := nextTermContracts[1].RightToUseDtos[1]
		assert.NotZero(t, recurRightToUse2.Id)
		assert.NotZero(t, recurRightToUse2.CreatedAt)
		assert.NotZero(t, recurRightToUse2.UpdatedAt)
		assert.True(t, recurRightToUse2.ValidFrom.Equal(utils.CreateJstTime(2020, 5, 30, 0, 0, 0, 0)))
		assert.True(t, recurRightToUse2.ValidTo.Equal(utils.CreateJstTime(2020, 6, 30, 0, 0, 0, 0)))
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
	interactor := create2.NewUserIndividualCreateInteractor(db.NewUserRepository())
	response, err := interactor.Handle(create2.NewUserIndividualCreateUseCaseRequest("個人たろう"))
	if err != nil || len(response.ValidationErrors) > 0 {
		panic("データ作成失敗")
	}
	return response.UserDto
}

package billing

import (
	"github.com/golang/mock/gomock"
	create2 "github.com/mixmaru/my_contracts/core/application/contracts/create"
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/application/products/create"
	create3 "github.com/mixmaru/my_contracts/core/application/users/create"
	"github.com/mixmaru/my_contracts/core/domain/models/contract"
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/core/infrastructure/mock/mock_product"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
)

func TestBillBillingInteractor_Handle(t *testing.T) {
	productRep := db.NewProductRepository()
	userRep := db.NewUserRepository()
	contractRep := db.NewContractRepository()
	billRep := db.NewBillRepository()

	productCreateInteractor := create.NewProductCreateInteractor(productRep)
	createUserInteractor := create3.NewUserIndividualCreateInteractor(userRep)
	createContractInteractor := create2.NewContractCreateInteractor(userRep, productRep, contractRep)
	billBillingInteractor := NewBillBillingInteractor(productRep, contractRep, billRep)

	t.Run("渡した日付時点で有効な使用権でまだ請求実行データ（billsテーブル）が作成されていないものに請求データを作成する", func(t *testing.T) {
		////// 準備 2020/6/1 ~ 2020/6/30の使用権を作成する
		// 影響するデータを事前削除しておく
		conn, err := db.GetConnection()
		assert.NoError(t, err)
		deleteQuery := `
DELETE FROM bill_details;
DELETE FROM right_to_use_active;
DELETE FROM right_to_use_history;
DELETE FROM right_to_use;
`
		_, err = conn.Exec(deleteQuery)
		assert.NoError(t, err)
		// 商品作成
		productCreateResponse, err := productCreateInteractor.Handle(create.NewProductCreateUseCaseRequest("商品", "1234"))
		assert.NoError(t, err)
		assert.Len(t, productCreateResponse.ValidationError, 0)
		// user作成
		userResponse, err := createUserInteractor.Handle(create3.NewUserIndividualCreateUseCaseRequest("請求実行テスト太郎"))
		assert.NoError(t, err)
		assert.Len(t, userResponse.ValidationErrors, 0)
		user := userResponse.UserDto
		// 契約作成（使用権も自動的に作成される）（課金開始日は6/2からになる。）
		contractResponse, err := createContractInteractor.Handle(create2.NewContractCreateUseCaseRequest(user.Id, productCreateResponse.ProductDto.Id, utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0)))
		assert.NoError(t, err)
		assert.Len(t, contractResponse.ValidationErrors, 0)

		// 実行 2020/6/2で請求実行する（課金開始日が6/2なので、その日を指定）
		billResponse, err := billBillingInteractor.Handle(NewBillBillingUseCaseRequest(utils.CreateJstTime(2020, 6, 2, 0, 0, 0, 0)))
		assert.NoError(t, err)
		billDtos := billResponse.BillDtos

		// 検証 準備で用意した使用権の使用量が請求データに作成されている
		assert.Len(t, billDtos, 1)
		assert.Equal(t, user.Id, billDtos[0].UserId)
		assert.True(t, billDtos[0].BillingDate.Equal(utils.CreateJstTime(2020, 6, 2, 0, 0, 0, 0)))
		assert.False(t, billDtos[0].PaymentConfirmed)
		assert.Zero(t, billDtos[0].PaymentConfirmedAt)
		billDetails := billDtos[0].BillDetails
		assert.Len(t, billDetails, 1)
		assert.Equal(t, "1234", billDetails[0].BillingAmount) // 6/2 ~ 7/2が使用権の期間になってるので日割りにはならない
	})

	t.Run("渡した日時を実行日として_請求を実行する（billsとbill_detailsデータを作成する）", func(t *testing.T) {
		t.Run("2020/7/1を渡すと_7/1時点で使用権が開始されていて克つ_契約の課金開始日以降である使用権の使用量が請求される", func(t *testing.T) {
			////// 準備
			// 事前に同日で実行してすべて請求実行済にしておく。テストのために。
			interactor := NewBillBillingInteractor(productRep, contractRep, billRep)
			_, err := interactor.Handle(NewBillBillingUseCaseRequest(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
			assert.NoError(t, err)
			// テストデータ作成
			user1Id, rightToUse1AId, rightToUse1BId, _ := createTestData(t)

			////// 実行
			response, err := interactor.Handle(NewBillBillingUseCaseRequest(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
			billDtos := response.BillDtos

			////// 検証（billingデータを取得して検証する。1ユーザーの6/1~6/30, 7/1~7/31の請求分がbillsに作成される）
			assert.Len(t, billDtos, 1)
			assert.NotZero(t, billDtos[0].Id)
			assert.NotZero(t, billDtos[0].CreatedAt)
			assert.NotZero(t, billDtos[0].UpdatedAt)
			assert.Equal(t, user1Id, billDtos[0].UserId)
			assert.False(t, billDtos[0].PaymentConfirmed)
			assert.Zero(t, billDtos[0].PaymentConfirmedAt)
			assert.True(t, billDtos[0].BillingDate.Equal(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
			assert.Equal(t, "4000", billDtos[0].TotalAmountExcludingTax)
			// details
			actualDetails := billDtos[0].BillDetails
			assert.Len(t, actualDetails, 2)
			// detail1つめ
			assert.NotZero(t, actualDetails[0].Id)
			assert.NotZero(t, actualDetails[0].CreatedAt)
			assert.NotZero(t, actualDetails[0].UpdatedAt)
			assert.Equal(t, "2000", actualDetails[0].BillingAmount)
			assert.Equal(t, rightToUse1AId, actualDetails[0].RightToUseId)
			// detail2つめ
			assert.NotZero(t, actualDetails[1].Id)
			assert.NotZero(t, actualDetails[1].CreatedAt)
			assert.NotZero(t, actualDetails[1].UpdatedAt)
			assert.Equal(t, "2000", actualDetails[1].BillingAmount)
			assert.Equal(t, rightToUse1BId, actualDetails[1].RightToUseId)
		})

		t.Run("対象がなくて請求データが作成されなかった場合は空スライスが返る", func(t *testing.T) {
			////// 準備
			// 事前に同日で実行してすべて請求実行済にしておく。テストのために。
			interactor := NewBillBillingInteractor(productRep, contractRep, billRep)
			_, err := interactor.Handle(NewBillBillingUseCaseRequest(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
			assert.NoError(t, err)

			////// 実行（既に実行されているので、この実行で新たに作成される請求は無いはず）
			response, err := interactor.Handle(NewBillBillingUseCaseRequest(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
			assert.NoError(t, err)

			////// 検証
			assert.Len(t, response.BillDtos, 0)
		})

		t.Run("処理途中で失敗したときエラーを返すが、処理したものに関しては保存し返却する", func(t *testing.T) {
			////// 準備
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 商品リポジトリのモック作成
			productRepMock := mock_product.NewMockIProductRepository(ctrl)
			callCount := 0
			productRepMock.EXPECT().GetByRightToUseId(gomock.Any(), gomock.Any()).DoAndReturn(func(id int, executor gorp.SqlExecutor) (*product.ProductEntity, error) {
				callCount += 1
				if callCount == 3 {
					// 3回目にエラーを発生させる
					return nil, errors.New("Productデータの取得に失敗しました")
				} else {
					rep := db.NewProductRepository()
					return rep.GetByRightToUseId(id, executor)
				}
			}).AnyTimes()
			// 事前に同日で実行してすべて請求実行済にしておく。テストのために。
			app := NewBillBillingInteractor(productRep, contractRep, billRep)
			_, err := app.Handle(NewBillBillingUseCaseRequest(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
			assert.NoError(t, err)
			// アプリケーションサービス作成
			intractor := NewBillBillingInteractor(productRepMock, contractRep, billRep)
			// テストデータ作成
			_, _, _, _ = createTestData(t)
			_, _, _, _ = createTestData(t)

			////// 実行（途中で失敗するのでエラーがでる）
			response, err := intractor.Handle(NewBillBillingUseCaseRequest(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
			assert.Error(t, err) // エラーが出る。

			////// 検証（処理が完了したものについては返却される）
			assert.Len(t, response.BillDtos, 1)
		})
	})
}

func createTestData(t *testing.T) (userId, rightToUse1Id, rightToUse2Id, rightToUse3Id int) {
	executor, err := db.GetConnection()
	contractRep := db.NewContractRepository()
	userRep := db.NewUserRepository()

	////// 準備（1ユーザーに対して、6/1~6/30, 7/1~7/31, 8/1~8/31の未請求使用権データを作成する）
	// 商品登録
	product := createProduct()
	// user登録
	user1, err := user.NewUserIndividualEntity("ユーザー1")
	assert.NoError(t, err)
	user1Id, err := userRep.SaveUserIndividual(user1, executor)
	assert.NoError(t, err)
	// 契約作成
	contract1 := contract.NewContractEntity(
		user1Id,
		product.Id,
		utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0),
		[]*contract.RightToUseEntity{
			contract.NewRightToUseEntity(
				utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
			),
			contract.NewRightToUseEntity(
				utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
			),
			contract.NewRightToUseEntity(
				utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 9, 1, 0, 0, 0, 0),
			),
		},
	)
	contract1Id, err := contractRep.Create(contract1, executor)
	assert.NoError(t, err)
	// リロード
	savedContract, err := contractRep.GetById(contract1Id, executor)
	rightToUses := savedContract.RightToUses()

	return user1Id, rightToUses[0].Id(), rightToUses[1].Id(), rightToUses[2].Id()
}

func createProduct() products.ProductDto {
	interactor := create.NewProductCreateInteractor(db.NewProductRepository())
	response, err := interactor.Handle(create.NewProductCreateUseCaseRequest("商品", "2000"))
	if err != nil || len(response.ValidationError) > 0 {
		panic("データ作成失敗")
	}

	return response.ProductDto
}

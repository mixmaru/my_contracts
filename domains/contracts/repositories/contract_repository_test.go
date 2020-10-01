package repositories

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContractRepository_Create(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	// userを作成
	userRepository := NewUserRepository()
	userEntity, err := entities.NewUserIndividualEntity("担当太郎")
	assert.NoError(t, err)
	savedUserId, err := userRepository.SaveUserIndividual(userEntity, db)
	assert.NoError(t, err)

	// 商品を登録
	productRepository := NewProductRepository()
	productEntity, err := entities.NewProductEntity("商品", "1000")
	assert.NoError(t, err)
	savedProductId, err := productRepository.Save(productEntity, db)
	assert.NoError(t, err)

	t.Run("UserIdとProductIdと契約日と課金開始日を渡すと契約が新規作成される", func(t *testing.T) {
		////// 準備
		// 使用権データ作成
		rightToUses := make([]*entities.RightToUseEntity, 0, 2)
		rightToUse1 := entities.NewRightToUseEntity(
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
		)
		rightToUse2 := entities.NewRightToUseEntity(
			utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0),
		)
		rightToUses = append(rightToUses, rightToUse1, rightToUse2)
		// 契約データ作成
		contractDate := utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)
		billingStartDate := utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0)
		contractEntity := entities.NewContractEntity(savedUserId, savedProductId, contractDate, billingStartDate, rightToUses)

		////// 実行
		contractRepository := NewContractRepository()
		savedContractId, err := contractRepository.Create(contractEntity, db)

		////// 検証
		assert.NoError(t, err)
		assert.NotZero(t, savedContractId)
		// 使用権が作られているかチェック
		count, err := db.SelectInt(`
SELECT COUNT(1)
FROM right_to_use_active rtua
    INNER JOIN right_to_use rtu ON rtua.right_to_use_id = rtu.id
    INNER JOIN contracts c ON c.id = rtu.contract_id
WHERE c.id = $1
GROUP BY contract_id
;`, savedContractId)
		assert.NoError(t, err)
		assert.Equal(t, 2, int(count))

	})

	t.Run("存在しないuserIdで作成されようとしたとき_エラーが出る", func(t *testing.T) {
		contractRepository := NewContractRepository()
		contractDate := utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)
		billingStartDate := utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0)
		contractEntity := entities.NewContractEntity(-100, savedProductId, contractDate, billingStartDate, []*entities.RightToUseEntity{})

		savedContractId, err := contractRepository.Create(contractEntity, db)

		assert.Error(t, err)
		assert.Zero(t, savedContractId)
	})

	t.Run("存在しないproductIDで作成されようとしたとき_エラーが出る", func(t *testing.T) {
		contractRepository := NewContractRepository()
		contractDate := utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)
		billingStartDate := utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0)
		contractEntity := entities.NewContractEntity(savedUserId, -100, contractDate, billingStartDate, []*entities.RightToUseEntity{})

		savedContractId, err := contractRepository.Create(contractEntity, db)

		assert.Error(t, err)
		assert.Zero(t, savedContractId)
	})
}

func TestContractRepository_GetById(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	// userを作成
	userRepository := NewUserRepository()
	userEntity, err := entities.NewUserIndividualEntity("担当太郎")
	assert.NoError(t, err)
	savedUserId, err := userRepository.SaveUserIndividual(userEntity, db)
	assert.NoError(t, err)

	// 商品を登録
	productRepository := NewProductRepository()
	productEntity, err := entities.NewProductEntity("商品", "1000")
	assert.NoError(t, err)
	savedProductId, err := productRepository.Save(productEntity, db)
	assert.NoError(t, err)

	t.Run("データがある時_Idで契約データと関連する商品データとユーザーデータを返す_一緒に使うことが多い気がするため", func(t *testing.T) {
		r := NewContractRepository()
		// データ登録
		contractEntity := entities.NewContractEntity(
			savedUserId,
			savedProductId,
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0),
			[]*entities.RightToUseEntity{
				entities.NewRightToUseEntity(
					utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
				),
				entities.NewRightToUseEntity(
					utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0),
				),
			},
		)
		savedId, err := r.Create(contractEntity, db)
		assert.NoError(t, err)

		// データ取得
		loadedContract, loadedProduct, loadedUser, err := r.GetById(savedId, db)
		assert.NoError(t, err)

		// contractテスト
		assert.Equal(t, savedId, loadedContract.Id())
		assert.Equal(t, savedUserId, loadedContract.UserId())
		assert.Equal(t, savedProductId, loadedContract.ProductId())
		assert.NotZero(t, loadedContract.CreatedAt())
		assert.NotZero(t, loadedContract.UpdatedAt())
		// rightToUse
		rightToUses := loadedContract.RightToUses()
		assert.Len(t, rightToUses, 2)
		assert.True(t, rightToUses[0].ValidFrom().Equal(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)))
		assert.True(t, rightToUses[0].ValidTo().Equal(utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0)))
		assert.True(t, rightToUses[1].ValidFrom().Equal(utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0)))
		assert.True(t, rightToUses[1].ValidTo().Equal(utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0)))
		// productテスト
		assert.Equal(t, savedProductId, loadedProduct.Id())
		assert.Equal(t, productEntity.Name(), loadedProduct.Name())
		price, exist := loadedProduct.MonthlyPrice()
		assert.True(t, exist)
		assert.Equal(t, "1000", price.String())
		assert.NotZero(t, loadedProduct.CreatedAt())
		assert.NotZero(t, loadedProduct.UpdatedAt())
		// userテスト
		user, ok := loadedUser.(*entities.UserIndividualEntity)
		assert.True(t, ok)
		assert.Equal(t, savedUserId, user.Id())
		assert.Equal(t, "担当太郎", user.Name())
		assert.NotZero(t, user.CreatedAt())
		assert.NotZero(t, user.UpdatedAt())
	})

	t.Run("データがない時はnilが返る", func(t *testing.T) {
		r := NewContractRepository()
		// データ取得
		loadedContract, loadedProduct, loadedUser, err := r.GetById(-100, db)
		assert.NoError(t, err)
		assert.Nil(t, loadedContract)
		assert.Nil(t, loadedProduct)
		assert.Nil(t, loadedUser)
	})
}

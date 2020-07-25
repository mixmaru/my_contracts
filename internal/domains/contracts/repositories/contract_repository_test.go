package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
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

	// 商品を事前に削除
	_, err = db.Exec(
		"delete from contracts " +
			"using products " +
			"where products.id = contracts.product_id " +
			"and products.name = '商品名' ")
	assert.NoError(t, err)
	_, err = db.Exec("delete from products where name = '商品名'")
	assert.NoError(t, err)

	// 商品を登録
	productRepository := NewProductRepository()
	productEntity, err := entities.NewProductEntity("商品名", "1000")
	assert.NoError(t, err)
	savedProductId, err := productRepository.Save(productEntity, db)
	assert.NoError(t, err)

	t.Run("正常系", func(t *testing.T) {
		// 契約作成テスト
		contractRepository := NewContractRepository()
		savedContractId, err := contractRepository.Create(entities.NewContractEntity(savedUserId, savedProductId), db)

		assert.NoError(t, err)
		assert.NotZero(t, savedContractId)
	})

	t.Run("存在しないuserIdで作成されようとしたとき", func(t *testing.T) {
		contractRepository := NewContractRepository()
		savedContractId, err := contractRepository.Create(entities.NewContractEntity(0, savedProductId), db)

		assert.Error(t, err)
		assert.Zero(t, savedContractId)
	})

	t.Run("存在しないproductIDで作成されようとしたとき", func(t *testing.T) {
		contractRepository := NewContractRepository()
		savedContractId, err := contractRepository.Create(entities.NewContractEntity(savedUserId, 0), db)

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

	// 商品を事前に削除
	_, err = db.Exec(
		"delete from contracts " +
			"using products " +
			"where products.id = contracts.product_id " +
			"and products.name = '商品名' ")
	assert.NoError(t, err)
	_, err = db.Exec("delete from products where name = '商品名'")
	assert.NoError(t, err)

	// 商品を登録
	productRepository := NewProductRepository()
	productEntity, err := entities.NewProductEntity("商品名", "1000")
	assert.NoError(t, err)
	savedProductId, err := productRepository.Save(productEntity, db)
	assert.NoError(t, err)

	t.Run("データがある時", func(t *testing.T) {
		r := NewContractRepository()
		// データ登録
		contractEntity := entities.NewContractEntity(savedUserId, savedProductId)
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
		// productテスト
		assert.Equal(t, savedProductId, loadedProduct.Id())
		assert.Equal(t, "商品名", loadedProduct.Name())
		price := loadedProduct.Price()
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

	t.Run("データがない時", func(t *testing.T) {
		r := NewContractRepository()
		// データ取得
		loadedContract, loadedProduct, loadedUser, err := r.GetById(-100, db)
		assert.NoError(t, err)
		assert.Nil(t, loadedContract)
		assert.Nil(t, loadedProduct)
		assert.Nil(t, loadedUser)
	})
}
package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
)

func TestRightToUseRepository_Create(t *testing.T) {
	r := NewRightToUseRepository()
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	////// テスト用契約を作成する
	// 契約の作成
	savedContractId := createPreparedContractData(db)

	t.Run("権利エンティティとdbコネクションを渡すとDBへ新規保存されて_保存Idを返す", func(t *testing.T) {
		// 準備
		entity := entities.NewRightToUseEntity(
			savedContractId,
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
		)

		// 実行
		savedId, err := r.Create(entity, db)
		assert.NoError(t, err)

		//検証
		assert.NotZero(t, savedId)
	})
}

func TestRightToUseRepository_GetById(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	// 事前に使用権を登録する
	r := NewRightToUseRepository()
	savedContractId := createPreparedContractData(db)
	rightToUseEntity := entities.NewRightToUseEntity(
		savedContractId,
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
	)

	// 実行
	savedId, err := r.Create(rightToUseEntity, db)
	assert.NoError(t, err)

	t.Run("データがあればIdを渡すとデータが取得できる", func(t *testing.T) {
		// 実行
		actual, err := r.GetById(savedId, db)
		assert.NoError(t, err)

		// 検証
		assert.Equal(t, savedId, actual.Id())
		assert.Equal(t, savedContractId, actual.ContractId())
		assert.True(t, actual.ValidFrom().Equal(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)))
		assert.True(t, actual.ValidTo().Equal(utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0)))
	})

	t.Run("データがなければidを渡すとnilが返る", func(t *testing.T) {
		// 実行
		actual, err := r.GetById(-100, db)
		assert.NoError(t, err)

		// 検証
		assert.Nil(t, actual)
	})
}

// 使用権データを作成するのに事前に必要なデータを準備する
func createPreparedContractData(executor gorp.SqlExecutor) int {
	// userの作成
	savedUserId := createUser(executor)
	// 商品の作成
	savedProductId := createProduct(executor)
	// 契約の作成
	savedContractId := createContract(savedUserId, savedProductId, executor)
	return savedContractId
}

func createUser(executor gorp.SqlExecutor) int {
	userEntity, err := entities.NewUserIndividualEntity("個人太郎")
	if err != nil {
		panic("userEntity作成失敗")
	}
	userRepository := NewUserRepository()
	savedUserId, err := userRepository.SaveUserIndividual(userEntity, executor)
	if err != nil {
		panic("userEntity保存失敗")
	}
	return savedUserId
}

func createProduct(executor gorp.SqlExecutor) int {
	productEntity, err := entities.NewProductEntity(utils.CreateUniqProductNameForTest(), "1000")
	if err != nil {
		panic("productEntity作成失敗")
	}
	productRepository := NewProductRepository()
	savedProductId, err := productRepository.Save(productEntity, executor)
	if err != nil {
		panic("productEntity保存失敗")
	}
	return savedProductId
}

func createContract(userId, productId int, executor gorp.SqlExecutor) int {
	// 契約の作成
	contractEntity := entities.NewContractEntity(
		userId,
		productId,
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
	)
	contractRepository := NewContractRepository()
	savedContractId, err := contractRepository.Create(contractEntity, executor)
	if err != nil {
		panic("contractEntity保存失敗")
	}
	return savedContractId
}

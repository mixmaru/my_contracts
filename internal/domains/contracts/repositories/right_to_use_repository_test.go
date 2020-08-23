package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRightToUseRepository_Create(t *testing.T) {
	r := NewRightToUseRepository()
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)

	////// テスト用契約を作成する
	// userの作成
	userEntity, err := entities.NewUserIndividualEntity("個人太郎")
	assert.NoError(t, err)
	userRepository := NewUserRepository()
	savedUserId, err := userRepository.SaveUserIndividual(userEntity, db)
	assert.NoError(t, err)

	// 商品の作成
	productEntity, err := entities.NewProductEntity(utils.CreateUniqProductNameForTest(), "1000")
	assert.NoError(t, err)
	productRepository := NewProductRepository()
	savedProductId, err := productRepository.Save(productEntity, db)

	// 契約の作成
	contractEntity := entities.NewContractEntity(
		savedUserId,
		savedProductId,
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
	)
	contractRepository := NewContractRepository()
	savedContractId, err := contractRepository.Create(contractEntity, db)
	assert.NoError(t, err)

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

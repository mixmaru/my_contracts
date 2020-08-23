package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// トランザクションが正しく動作しているかテスト
func TestUserRepository_Transaction(t *testing.T) {
	// db接続
	dbMap, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer dbMap.Db.Close()

	t.Run("コミットするとデータ保存されている", func(t *testing.T) {
		// transaction開始
		tran, err := dbMap.Begin()
		assert.NoError(t, err)

		//データ保存
		user, err := entities.NewUserIndividualEntity("個人太郎")
		assert.NoError(t, err)
		repo := NewUserRepository()
		savedId, err := repo.SaveUserIndividual(user, tran)
		assert.NoError(t, err)

		// コミット
		err = tran.Commit()
		assert.NoError(t, err)

		// データ取得できる
		_, err = repo.GetUserIndividualById(savedId, dbMap)
		assert.NoError(t, err) // sql: no rows in result set エラーが起こらなければ、データが保存されている
	})

	t.Run("ロールバックするとデータ保存されていない", func(t *testing.T) {
		// transaction開始
		tran, err := dbMap.Begin()
		assert.NoError(t, err)

		//データ保存
		user, err := entities.NewUserIndividualEntity("個人太郎")
		assert.NoError(t, err)
		repo := NewUserRepository()
		savedId, err := repo.SaveUserIndividual(user, tran)
		assert.NoError(t, err)

		// ロールバック
		err = tran.Rollback()
		assert.NoError(t, err)

		// データ取得できない
		user, err = repo.GetUserIndividualById(savedId, dbMap)
		assert.Nil(t, user)
	})
}

func TestUserRepository_SaveUserIndividual(t *testing.T) {
	// 登録用データ作成
	user, err := entities.NewUserIndividualEntity("個人太郎")
	assert.NoError(t, err)

	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	// 実行
	repo := NewUserRepository()
	_, err = repo.SaveUserIndividual(user, db)
	assert.NoError(t, err)
}

func TestUserRepository_GetUserIndividualById(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	//　事前にデータ登録する
	user, err := entities.NewUserIndividualEntity("個人太郎")
	assert.NoError(t, err)
	repo := NewUserRepository()
	savedId, err := repo.SaveUserIndividual(user, db)
	assert.NoError(t, err)

	// idで取得して検証
	t.Run("データがある時_idでデータが取得できる", func(t *testing.T) {
		result, err := repo.GetUserIndividualById(savedId, db)
		assert.NoError(t, err)
		assert.Equal(t, user.Name(), result.Name())
		assert.NotEqual(t, time.Time{}, result.CreatedAt())
		assert.NotEqual(t, time.Time{}, result.UpdatedAt())
	})

	t.Run("データが無い時_nilが返る", func(t *testing.T) {
		user, err := repo.GetUserIndividualById(-1, db)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_GetUserCorporationById(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	//　事前にデータ登録する
	savingUser, err := entities.NewUserCorporationEntity("イケてる会社", "担当　太郎", "社長　太郎")
	assert.NoError(t, err)

	repo := NewUserRepository()
	savedId, err := repo.SaveUserCorporation(savingUser, db)
	assert.NoError(t, err)

	// idで取得して検証
	t.Run("データがある時_idでデータが取得できる", func(t *testing.T) {
		result, err := repo.GetUserCorporationById(savedId, db)
		assert.NoError(t, err)
		assert.Equal(t, savedId, result.Id())
		assert.Equal(t, "イケてる会社", result.CorporationName())
		assert.Equal(t, "担当　太郎", result.ContactPersonName())
		assert.Equal(t, "社長　太郎", result.PresidentName())
		assert.NotEqual(t, time.Time{}, result.CreatedAt())
		assert.NotEqual(t, time.Time{}, result.UpdatedAt())
	})

	t.Run("データが無い時_nilが返る", func(t *testing.T) {
		result, err := repo.GetUserCorporationById(-1, db)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestUserRepository_SaveUserCorporation(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	// 保存するデータ作成
	user, err := entities.NewUserCorporationEntity("イケてる会社", "担当太郎", "社長次郎")
	assert.NoError(t, err)

	// 保存実行
	repo := NewUserRepository()
	_, err = repo.SaveUserCorporation(user, db)
	assert.NoError(t, err)
}

func TestUserRepository_getUserCorporationViewById(t *testing.T) {
	// db接続
	dbMap, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer dbMap.Db.Close()

	// 事前にデータ登録
	user, err := entities.NewUserCorporationEntity("イケてる会社", "担当太郎", "社長次郎")
	assert.NoError(t, err)
	repo := NewUserRepository()
	savedId, err := repo.SaveUserCorporation(user, dbMap)
	assert.NoError(t, err)

	// idで取得する
	result, err := repo.getUserCorporationEntityById(savedId, &entities.UserCorporationEntity{}, dbMap)
	assert.NoError(t, err)

	// 検証
	assert.Equal(t, result.Id(), savedId)
	assert.Equal(t, "イケてる会社", result.CorporationName())
	assert.Equal(t, "担当太郎", result.ContactPersonName())
	assert.Equal(t, "社長次郎", result.PresidentName())
	assert.NotEqual(t, time.Time{}, result.CreatedAt())
	assert.NotEqual(t, time.Time{}, result.UpdatedAt())
}

func TestUserRepository_GetUserById(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	//　事前にデータ登録する。個人顧客
	userIndividual, err := entities.NewUserIndividualEntity("個人太郎")
	assert.NoError(t, err)
	repo := NewUserRepository()
	savedIndividualId, err := repo.SaveUserIndividual(userIndividual, db)
	assert.NoError(t, err)

	//　事前にデータ登録する。法人顧客
	userCorporation, err := entities.NewUserCorporationEntity("イケてる会社", "担当太郎", "社長太郎")
	assert.NoError(t, err)
	repo = NewUserRepository()
	savedCorporationId, err := repo.SaveUserCorporation(userCorporation, db)
	assert.NoError(t, err)

	// idで取得して検証
	t.Run("個人顧客データ取得", func(t *testing.T) {
		t.Run("データがある時_interface{}型でUserIndividualEntityが返る", func(t *testing.T) {
			result, err := repo.GetUserById(savedIndividualId, db)
			assert.NoError(t, err)

			loadedIndividual, ok := result.(*entities.UserIndividualEntity)
			assert.True(t, ok)

			assert.Equal(t, savedIndividualId, loadedIndividual.Id())
			assert.Equal(t, "個人太郎", loadedIndividual.Name())
			assert.NotZero(t, loadedIndividual.CreatedAt())
			assert.NotZero(t, loadedIndividual.UpdatedAt())
		})

		t.Run("データが無い時_nilが返る", func(t *testing.T) {
			user, err := repo.GetUserById(-1, db)
			assert.NoError(t, err)
			assert.Nil(t, user)
		})
	})

	t.Run("法人顧客データ取得", func(t *testing.T) {
		t.Run("データがある時_interface{}型でUserCorporationEntityが返る", func(t *testing.T) {
			result, err := repo.GetUserById(savedCorporationId, db)
			assert.NoError(t, err)

			loadedCorporation, ok := result.(*entities.UserCorporationEntity)
			assert.True(t, ok)

			assert.Equal(t, savedCorporationId, loadedCorporation.Id())
			assert.Equal(t, "イケてる会社", loadedCorporation.CorporationName())
			assert.Equal(t, "担当太郎", loadedCorporation.ContactPersonName())
			assert.Equal(t, "社長太郎", loadedCorporation.PresidentName())
			assert.NotZero(t, loadedCorporation.CreatedAt())
			assert.NotZero(t, loadedCorporation.UpdatedAt())
		})

		t.Run("データが無い時_nilが返る", func(t *testing.T) {
			user, err := repo.GetUserById(-1, db)
			assert.NoError(t, err)
			assert.Nil(t, user)
		})
	})
}

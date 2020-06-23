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
		repo := UserRepository{}
		user, err = repo.SaveUserIndividual(user, tran)
		assert.NoError(t, err)

		// コミット
		err = tran.Commit()
		assert.NoError(t, err)

		// データ取得できる
		_, err = repo.GetUserIndividualById(user.Id(), nil)
		assert.NoError(t, err) // sql: no rows in result set エラーが起こらなければ、データが保存されている
	})

	t.Run("ロールバックするとデータ保存されていない", func(t *testing.T) {
		// transaction開始
		tran, err := dbMap.Begin()
		assert.NoError(t, err)

		//データ保存
		user, err := entities.NewUserIndividualEntity("個人太郎")
		assert.NoError(t, err)
		repo := UserRepository{}
		user, err = repo.SaveUserIndividual(user, tran)
		assert.NoError(t, err)

		// ロールバック
		err = tran.Rollback()
		assert.NoError(t, err)

		// データ取得できない
		user, err = repo.GetUserIndividualById(user.Id(), nil)
		assert.Nil(t, user)
	})
}

func TestUserRepository_SaveUserIndividual(t *testing.T) {
	// 登録用データ作成
	user, err := entities.NewUserIndividualEntity("個人太郎")
	assert.NoError(t, err)

	// 実行
	repo := UserRepository{}
	user, err = repo.SaveUserIndividual(user, nil)
	assert.NoError(t, err)
}

func TestUserRepository_GetUserIndividualById(t *testing.T) {
	//　事前にデータ登録する
	user, err := entities.NewUserIndividualEntity("個人太郎")
	assert.NoError(t, err)
	repo := &UserRepository{}
	user, err = repo.SaveUserIndividual(user, nil)
	assert.NoError(t, err)

	// idで取得して検証
	t.Run("データがある時", func(t *testing.T) {
		result, err := repo.GetUserIndividualById(user.Id(), nil)
		assert.NoError(t, err)
		assert.Equal(t, result.Id(), user.Id())
		assert.Equal(t, result.Name(), user.Name())
		assert.NotEqual(t, time.Time{}, user.CreatedAt())
		assert.NotEqual(t, time.Time{}, user.UpdatedAt())
	})

	t.Run("データが無い時", func(t *testing.T) {
		user, err := repo.GetUserIndividualById(-1, nil)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_GetUserCorporationById(t *testing.T) {
	//　事前にデータ登録する
	savingUser := entities.NewUserCorporationEntity()
	err := savingUser.SetContactPersonName("担当　太郎")
	assert.NoError(t, err)
	err = savingUser.SetPresidentName("社長　太郎")
	assert.NoError(t, err)

	repo := &UserRepository{}
	savedUser, err := repo.SaveUserCorporation(savingUser, nil)
	assert.NoError(t, err)

	// idで取得して検証
	t.Run("データがある時", func(t *testing.T) {
		result, err := repo.GetUserCorporationById(savedUser.Id(), nil)
		assert.NoError(t, err)
		assert.Equal(t, savedUser.Id(), result.Id())
		assert.Equal(t, "担当　太郎", result.ContactPersonName())
		assert.Equal(t, "社長　太郎", result.PresidentName())
		assert.NotEqual(t, time.Time{}, result.CreatedAt())
		assert.NotEqual(t, time.Time{}, result.UpdatedAt())
	})

	t.Run("データが無い時", func(t *testing.T) {
		result, err := repo.GetUserCorporationById(-1, nil)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestUserRepository_SaveUserCorporation(t *testing.T) {
	// 保存するデータ作成
	user := entities.NewUserCorporationEntity()
	user.SetContactPersonName("担当太郎")
	user.SetPresidentName("社長次郎")

	// 保存実行
	repo := &UserRepository{}
	user, err := repo.SaveUserCorporation(user, nil)
	assert.NoError(t, err)
}

func TestUserRepository_getUserCorporationViewById(t *testing.T) {
	// db接続
	dbMap, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer dbMap.Db.Close()

	// 事前にデータ登録
	user := entities.NewUserCorporationEntity()
	user.SetContactPersonName("担当太郎")
	user.SetPresidentName("社長次郎")
	repo := &UserRepository{}
	user, err = repo.SaveUserCorporation(user, nil)
	assert.NoError(t, err)

	// idで取得する
	result, err := repo.getUserCorporationEntityById(user.Id(), &entities.UserCorporationEntity{}, dbMap)
	assert.NoError(t, err)

	// 検証
	assert.Equal(t, result.Id(), user.Id())
	assert.Equal(t, "担当太郎", user.ContactPersonName())
	assert.Equal(t, "社長次郎", user.PresidentName())
	assert.NotEqual(t, time.Time{}, user.CreatedAt())
	assert.NotEqual(t, time.Time{}, user.UpdatedAt())
}

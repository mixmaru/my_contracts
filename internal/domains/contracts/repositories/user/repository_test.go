package user

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUser_InitDb(t *testing.T) {
	dbMap, err := InitDb()
	defer dbMap.Db.Close()
	assert.NoError(t, err)
}

// トランザクションが正しく動作しているかテスト
func TestUser_Transaction(t *testing.T) {
	// db接続
	dbMap, err := InitDb()
	assert.NoError(t, err)
	defer dbMap.Db.Close()

	t.Run("コミットするとデータ保存されている", func(t *testing.T) {
		// transaction開始
		tran, err := dbMap.Begin()
		assert.NoError(t, err)

		//データ保存
		user := user.NewUserIndividualEntity()
		user.SetName("個人太郎")
		repo := Repository{}
		err = repo.SaveUserIndividual(user, tran)
		assert.NoError(t, err)

		// コミット
		err = tran.Commit()
		assert.NoError(t, err)

		// データ取得できる
		_, err = repo.GetUserIndividualById(user.Id(), dbMap)
		assert.NoError(t, err) // sql: no rows in result set エラーが起こらなければ、データが保存されている
	})

	t.Run("ロールバックするとデータ保存されていない", func(t *testing.T) {
		// transaction開始
		tran, err := dbMap.Begin()
		assert.NoError(t, err)

		//データ保存
		user := user.NewUserIndividualEntity()
		user.SetName("個人太郎")
		repo := Repository{}
		err = repo.SaveUserIndividual(user, tran)
		assert.NoError(t, err)

		// ロールバック
		err = tran.Rollback()
		assert.NoError(t, err)

		// データ取得できない
		_, err = repo.GetUserIndividualById(user.Id(), dbMap)
		assert.Error(t, err) // sql: no rows in result set エラーが起こる
	})
}

func TestUser_SaveUserIndividual(t *testing.T) {
	// db接続
	dbMap, err := InitDb()
	assert.NoError(t, err)
	defer dbMap.Db.Close()

	// 登録用データ作成
	user := user.NewUserIndividualEntity()
	user.SetName("個人太郎")

	// 実行
	repo := Repository{}
	err = repo.SaveUserIndividual(user, dbMap)
	assert.NoError(t, err)
}

func TestUser_GetUserIndividualById(t *testing.T) {
	// db接続
	dbMap, err := InitDb()
	assert.NoError(t, err)
	defer dbMap.Db.Close()

	//　事前にデータ登録する
	user := user.NewUserIndividualEntity()
	user.SetName("個人太郎")
	repo := &Repository{}
	err = repo.SaveUserIndividual(user, dbMap)
	assert.NoError(t, err)

	// idで取得して検証
	result, err := repo.GetUserIndividualById(user.Id(), dbMap)
	assert.NoError(t, err)
	assert.Equal(t, result.Id(), user.Id())
	assert.Equal(t, result.Name(), user.Name())
}

func TestUser_SaveUserCorporation(t *testing.T) {
	// db接続
	dbMap, err := InitDb()
	assert.NoError(t, err)
	defer dbMap.Db.Close()

	// 保存するデータ作成
	user := user.NewUserCorporationEntity()
	user.SetContactPersonName("担当太郎")
	user.SetPresidentName("社長次郎")

	// 保存実行
	repo := &Repository{}
	err = repo.SaveUserCorporation(user, dbMap)
	assert.NoError(t, err)
}

func TestUser_getUserCorporationViewById(t *testing.T) {
	// db接続
	dbMap, err := InitDb()
	assert.NoError(t, err)
	defer dbMap.Db.Close()

	// 事前にデータ登録
	user := user.NewUserCorporationEntity()
	user.SetContactPersonName("担当太郎")
	user.SetPresidentName("社長次郎")
	repo := &Repository{}
	err = repo.SaveUserCorporation(user, dbMap)
	assert.NoError(t, err)

	// idで取得する
	result, err := repo.getUserCorporationEntityById(user.Id(), dbMap)
	assert.NoError(t, err)

	// 検証
	assert.Equal(t, result.Id(), user.Id())
	assert.Equal(t, "担当太郎", user.ContactPersonName())
	assert.Equal(t, "社長次郎", user.PresidentName())
	assert.NotEqual(t, time.Time{}, user.CreatedAt())
	assert.NotEqual(t, time.Time{}, user.UpdatedAt())
}

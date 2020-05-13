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

func TestUser_SaveUserIndividual(t *testing.T) {
	// 事前準備
	dbMap, err := InitDb()
	defer dbMap.Db.Close()

	assert.NoError(t, err)
	repo := Repository{}
	user := user.NewUserIndividualEntity()
	user.SetName("個人太郎")

	// 実行
	err = repo.SaveUserIndividual(user, dbMap)
	assert.NoError(t, err)

	// idが0ではない。db登録されたidとcreatedAtとupdatedAtが入ってる
	assert.NotEqual(t, 0, user.Id())
	assert.NotEqual(t, time.Time{}, user.CreatedAt())
	assert.NotEqual(t, time.Time{}, user.UpdatedAt())
}

func TestUser_GetUserIndividualById(t *testing.T) {
	//　事前にデータ登録
	dbMap, err := InitDb()
	defer dbMap.Db.Close()

	assert.NoError(t, err)
	repo := &Repository{}
	user := user.NewUserIndividualEntity()
	user.SetName("個人太郎")
	err = repo.SaveUserIndividual(user, dbMap)
	assert.NoError(t, err)

	// idで取得する
	result, err := repo.GetUserIndividualById(user.Id(), dbMap)
	assert.NoError(t, err)
	assert.Equal(t, result.Id(), user.Id())
	assert.Equal(t, result.Name(), user.Name())
}

func TestUser_SaveUserCorporation(t *testing.T) {
	// entity作成
	user := user.NewUserCorporationEntity()
	user.SetContactPersonName("担当太郎")
	user.SetPresidentName("社長次郎")

	// db接続用意
	dbMap, err := InitDb()
	defer dbMap.Db.Close()
	assert.NoError(t, err)

	// repository用意
	repo := &Repository{}

	// 保存実行
	err = repo.SaveUserCorporation(user, dbMap)
	assert.NoError(t, err)

	// データ取得して内容確認する
	result, err := repo.getUserCorporationEntityById(user.Id(), dbMap)
	assert.NoError(t, err)

	assert.Equal(t, user.Id(), result.Id())
	assert.Equal(t, "担当太郎", result.ContactPersonName())
	assert.Equal(t, "社長次郎", result.PresidentName())
	assert.NotEqual(t, time.Time{}, result.CreatedAt())
	assert.NotEqual(t, time.Time{}, result.UpdatedAt())
}

func TestUser_getUserCorporationViewById(t *testing.T) {
	dbMap, err := InitDb()
	defer dbMap.Db.Close()
	assert.NoError(t, err)

	//　事前にデータ登録
	repo := &Repository{}
	user := user.NewUserCorporationEntity()
	user.SetContactPersonName("担当太郎")
	user.SetPresidentName("社長次郎")
	err = repo.SaveUserCorporation(user, dbMap)
	assert.NoError(t, err)

	// idで取得する
	result, err := repo.getUserCorporationEntityById(user.Id(), dbMap)
	assert.NoError(t, err)
	assert.Equal(t, result.Id(), user.Id())
	assert.Equal(t, "担当太郎", user.ContactPersonName())
	assert.Equal(t, "社長次郎", user.PresidentName())
	assert.NotEqual(t, time.Time{}, user.CreatedAt())
	assert.NotEqual(t, time.Time{}, user.UpdatedAt())
}

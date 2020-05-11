package user

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUser_InitDb(t *testing.T) {
	_, err := InitDb()
	assert.NoError(t, err)
}

func TestUser_Save(t *testing.T) {
	// 事前準備
	db, err := InitDb()
	assert.NoError(t, err)
	repo := Repository{}
	user := user.NewUserIndividualEntity()
	user.SetName("個人太郎")

	// 実行
	err = repo.Save(user, db)
	assert.NoError(t, err)

	// idが0ではない。db登録されたidとcreatedAtとupdatedAtが入ってる
	assert.NotEqual(t, 0, user.Id())
	assert.NotEqual(t, time.Time{}, user.CreatedAt())
	assert.NotEqual(t, time.Time{}, user.UpdatedAt())
}

func TestUser_GetUserIndividualById(t *testing.T) {
	//　事前にデータ登録
	db, err := InitDb()
	assert.NoError(t, err)
	repo := Repository{}
	user := user.NewUserIndividualEntity()
	user.SetName("個人太郎")
	err = repo.Save(user, db)
	assert.NoError(t, err)

	// idで取得する
	result, err := repo.GetUserIndividualById(user.Id(), db)
	assert.NoError(t, err)
	assert.Equal(t, result.Id(), user.Id())
	assert.Equal(t, result.Name(), user.Name())
}

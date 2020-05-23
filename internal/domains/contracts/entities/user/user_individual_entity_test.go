package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// UserIndividualのインスタンス化をテスト
func TestUserIndividual_NewIndividual(t *testing.T) {
	// インスタンス化
	userIndividual := NewUserIndividualEntity()
	userIndividual.SetName("顧客太郎")

	// テスト
	assert.Equal(t, "顧客太郎", userIndividual.Name())
}

// 個人顧客Entityの初期化と共にデータロードするやつ
func TestUserIndividual_Static_LoadUserIndividual(t *testing.T) {
	result := NewUserIndividualEntityWithData(
		1,
		"個人太郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)

	assert.Equal(t, 1, result.Id())
	assert.Equal(t, "個人太郎", result.Name())
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), result.CreatedAt())
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), result.UpdatedAt())
}

// インスタンス化された個人顧客Entityに対してデータロードするやつ
func TestUserIndividual_LoadUserIndividual(t *testing.T) {
	userIndividual := NewUserIndividualEntity()
	userIndividual.LoadData(
		1,
		"個人太郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)

	assert.Equal(t, 1, userIndividual.Id())
	assert.Equal(t, "個人太郎", userIndividual.Name())
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), userIndividual.CreatedAt())
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), userIndividual.UpdatedAt())
}

func TestUserIndividualEntity_NewName(t *testing.T) {
	name, err := NewName("個人顧客名")
	assert.NoError(t, err)
	assert.Equal(t, Name{"個人顧客名"}, name)
}

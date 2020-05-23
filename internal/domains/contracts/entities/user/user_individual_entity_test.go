package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// UserIndividualのインスタンス化をテスト
func TestUserIndividual_NewIndividual(t *testing.T) {
	// インスタンス化
	userIndividual, err := NewUserIndividualEntity("顧客太郎")
	assert.NoError(t, err)

	// テスト
	assert.Equal(t, "顧客太郎", userIndividual.Name())
}

// 個人顧客Entityの初期化と共にデータロードするやつ
func TestUserIndividual_Static_LoadUserIndividual(t *testing.T) {
	result, err := NewUserIndividualEntityWithData(
		1,
		"個人太郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	assert.Equal(t, 1, result.Id())
	assert.Equal(t, "個人太郎", result.Name())
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), result.CreatedAt())
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), result.UpdatedAt())
}

// インスタンス化された個人顧客Entityに対してデータロードするやつ
func TestUserIndividual_LoadUserIndividual(t *testing.T) {
	userIndividual, err := NewUserIndividualEntity("既存太郎")
	assert.NoError(t, err)
	err = userIndividual.LoadData(
		1,
		"個人太郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	assert.Equal(t, 1, userIndividual.Id())
	assert.Equal(t, "個人太郎", userIndividual.Name())
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), userIndividual.CreatedAt())
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), userIndividual.UpdatedAt())
}

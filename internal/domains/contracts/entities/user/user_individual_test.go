package user

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user/structures"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// UserIndividualのインスタンス化をテスト
func TestUserIndividual_NewIndividual(t *testing.T) {
	// インスタンス化
	userIndividual := NewUserIndividual()
	userIndividual.SetName("顧客太郎")

	// テスト
	assert.Equal(t, "顧客太郎", userIndividual.Name())
}

// 個人顧客Entityの初期化と共にデータロードするやつ
func TestUserIndividual_Static_LoadUserIndividual(t *testing.T) {
	data := &structures.UserIndividualView{}
	data.Id = 1
	data.Name = "個人太郎"
	data.CreatedAt = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	data.UpdatedAt = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	result := LoadUserIndividual(data)

	assert.Equal(t, 1, result.Id())
	assert.Equal(t, data.Name, result.Name())
	assert.Equal(t, data.CreatedAt, result.CreatedAt())
	assert.Equal(t, data.UpdatedAt, result.UpdatedAt())
}

// インスタンス化された個人顧客Entityに対してデータロードするやつ
func TestUserIndividual_LoadUserIndividual(t *testing.T) {
	data := &structures.UserIndividualView{}
	data.Id = 1
	data.Name = "個人太郎"
	data.CreatedAt = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	data.UpdatedAt = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	userIndividual := NewUserIndividual()
	userIndividual.LoadUserIndividual(data)

	assert.Equal(t, 1, userIndividual.Id())
	assert.Equal(t, data.Name, userIndividual.Name())
	assert.Equal(t, data.CreatedAt, userIndividual.CreatedAt())
	assert.Equal(t, data.UpdatedAt, userIndividual.UpdatedAt())
}

package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// UserIndividualのインスタンス化をテスト
func TestUserIndividual_NewIndividual(t *testing.T) {
	// インスタンス化
	userIndividual, errors := NewUserIndividualEntity("顧客太郎")
	assert.Len(t, errors, 0)

	// テスト
	assert.Equal(t, "顧客太郎", userIndividual.Name())
}

// 個人顧客Entityの初期化と共にデータロードするやつ
func TestUserIndividual_Static_LoadUserIndividual(t *testing.T) {
	result, errors := NewUserIndividualEntityWithData(
		1,
		"個人太郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.Len(t, errors, 0)

	assert.Equal(t, 1, result.Id())
	assert.Equal(t, "個人太郎", result.Name())
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), result.CreatedAt())
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), result.UpdatedAt())
}

// インスタンス化された個人顧客Entityに対してデータロードするやつ
func TestUserIndividual_LoadUserIndividual(t *testing.T) {
	userIndividual, errors := NewUserIndividualEntity("既存太郎")
	assert.Len(t, errors, 0)
	errors = userIndividual.LoadData(
		1,
		"個人太郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.Len(t, errors, 0)

	assert.Equal(t, 1, userIndividual.Id())
	assert.Equal(t, "個人太郎", userIndividual.Name())
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), userIndividual.CreatedAt())
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), userIndividual.UpdatedAt())
}

func TestUserIndividualEntity_NewName(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		name, errors := NewName("個人顧客名")
		assert.Len(t, errors, 0)
		assert.Equal(t, Name{"個人顧客名"}, name)
	})

	t.Run("名前が空文字だった時", func(t *testing.T) {
		name, errors := NewName("")
		assert.Len(t, errors, 1)
		assert.IsType(t, EmptyValidError{}, errors[0])
		assert.Equal(t, Name{}, name)
	})

	t.Run("名前が50文字を超えていた時", func(t *testing.T) {
		name, errors := NewName("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.Len(t, errors, 1)
		assert.IsType(t, OverLengthValidError{}, errors[0])
		assert.Equal(t, Name{}, name)
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		name, errors := NewName("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.Len(t, errors, 0)
		assert.Equal(t, Name{"0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789"}, name)
	})
}

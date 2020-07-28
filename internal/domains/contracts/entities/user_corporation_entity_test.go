package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// UserCorporationのインスタンス化をテスト
func TestUserCorporationEntity_Instantiate(t *testing.T) {
	// インスタンス化
	user := &UserCorporationEntity{}
	err := user.SetContactPersonName("担当太郎")
	assert.NoError(t, err)
	err = user.SetPresidentName("社長太郎")
	assert.NoError(t, err)
	err = user.SetCorporationName("会社名")
	assert.NoError(t, err)

	// テスト
	assert.Equal(t, "担当太郎", user.ContactPersonName())
	assert.Equal(t, "社長太郎", user.PresidentName())
	assert.Equal(t, "会社名", user.CorporationName())
}

func TestUserCorporationEntity_NewUserCorporationEntity(t *testing.T) {
	// インスタンス化
	user, err := NewUserCorporationEntity("イケてる会社", "担当太郎", "社長次郎")
	assert.NoError(t, err)

	// テスト
	assert.Equal(t, 0, user.Id())
	assert.Equal(t, "イケてる会社", user.CorporationName())
	assert.Equal(t, "担当太郎", user.ContactPersonName())
	assert.Equal(t, "社長次郎", user.PresidentName())
	assert.Equal(t, time.Time{}, user.UpdatedAt())
	assert.Equal(t, time.Time{}, user.CreatedAt())
}

func TestUserCorporationEntity_NewUserCorporationEntityWithData(t *testing.T) {
	// インスタンス化
	user, err := NewUserCorporationEntityWithData(
		1,
		"イケイケ会社",
		"担当太郎",
		"社長次郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	// テスト
	assert.Equal(t, 1, user.Id())
	assert.Equal(t, "イケイケ会社", user.CorporationName())
	assert.Equal(t, "担当太郎", user.ContactPersonName())
	assert.Equal(t, "社長次郎", user.PresidentName())
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), user.CreatedAt())
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), user.UpdatedAt())
}

func TestUserCorporationEntity_LoadData(t *testing.T) {
	// インスタンス化
	user := &UserCorporationEntity{}
	err := user.LoadData(
		1,
		"イケてる会社",
		"担当太郎",
		"社長次郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	// テスト
	assert.Equal(t, 1, user.Id())
	assert.Equal(t, "イケてる会社", user.CorporationName())
	assert.Equal(t, "担当太郎", user.ContactPersonName())
	assert.Equal(t, "社長次郎", user.PresidentName())
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), user.CreatedAt())
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), user.UpdatedAt())
}

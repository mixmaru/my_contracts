package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// UserCorporationのインスタンス化をテスト
func TestUserCorporationEntity_Instantiate(t *testing.T) {
	// インスタンス化
	user := &UserCorporationEntity{}
	user.SetContactPersonName("担当太郎")
	user.SetPresidentName("社長太郎")

	// テスト
	assert.Equal(t, "担当太郎", user.ContactPersonName())
	assert.Equal(t, "社長太郎", user.PresidentName())
}

func TestUserCorporationEntity_NewUserCorporationEntity(t *testing.T) {
	// インスタンス化
	user := NewUserCorporationEntity()
	user.SetContactPersonName("担当太郎")
	user.SetPresidentName("社長次郎")

	// テスト
	assert.Equal(t, 0, user.Id())
	assert.Equal(t, "担当太郎", user.ContactPersonName())
	assert.Equal(t, "社長次郎", user.PresidentName())
	assert.Equal(t, time.Time{}, user.UpdatedAt())
	assert.Equal(t, time.Time{}, user.CreatedAt())
}

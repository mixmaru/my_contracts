package user_corporation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// UserCorporationのインスタンス化をテスト
func TestUserCorporation_Instantiate(t *testing.T) {
	// インスタンス化
	user := &UserCorporation{}
	user.SetContactPersonName("担当太郎")
	user.SetPresidentName("社長太郎")

	// テスト
	assert.Equal(t, "担当太郎", user.ContactPersonName())
	assert.Equal(t, "社長太郎", user.PresidentName())
}

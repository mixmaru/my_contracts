package application_service

import (
	user_repository "github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// インスタンス化テスト
func TestUserApplication_NewUserApplicationService(t *testing.T) {
	userApp := NewUserApplicationService(&user_repository.Repository{})
	assert.IsType(t, &UserApplicationService{}, userApp)
}

// 個人顧客情報の登録とデータ取得のテスト
func TestUserApplicationService_RegisterUserIndividual(t *testing.T) {
	// インスタンス化
	userApp := NewUserApplicationService(&user_repository.Repository{})

	userId, err := userApp.RegisterUserIndividual("個人太郎")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, userId)

	// IDでuser情報を取得してチェックする
	user, err := userApp.GetUserIndividual(userId)
	assert.NoError(t, err)
	assert.Equal(t, userId, user.Id)
	assert.Equal(t, "個人太郎", user.Name)
	assert.NotEqual(t, time.Time{}, user.CreatedAt)
	assert.NotEqual(t, time.Time{}, user.UpdatedAt)

}

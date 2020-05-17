package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	user_repository "github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserApplication_NewUserApplicationService(t *testing.T) {
	// リポジトリモックを用意する
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)

	userApp := NewUserApplicationService(userRepositoryMock)
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

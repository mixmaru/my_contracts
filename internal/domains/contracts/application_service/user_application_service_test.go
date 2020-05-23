package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	user2 "github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	user_repository "github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
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
	// リポジトリのSaveUserIndividual()が受け取る引数を用意
	saveUserEntity, err := user2.NewUserIndividualEntity("個人太郎")
	assert.NoError(t, err)

	now := time.Now()
	returnUserEntity, err := user2.NewUserIndividualEntity("既存太郎")
	assert.NoError(t, err)
	err = returnUserEntity.LoadData(1, "個人太郎", now, now)
	assert.NoError(t, err)

	// モックリポジトリ作成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)
	userRepositoryMock.EXPECT().
		SaveUserIndividual(
			saveUserEntity,
			gomock.AssignableToTypeOf(&gorp.Transaction{}),
		).Return(returnUserEntity, nil).
		Times(1)

	// インスタンス化
	userApp := NewUserApplicationService(userRepositoryMock)

	registerdUser, err := userApp.RegisterUserIndividual("個人太郎")
	assert.NoError(t, err)
	assert.Equal(t, 1, registerdUser.Id)
	assert.Equal(t, "個人太郎", registerdUser.Name)
	assert.Equal(t, now, registerdUser.CreatedAt)
	assert.Equal(t, now, registerdUser.UpdatedAt)
}

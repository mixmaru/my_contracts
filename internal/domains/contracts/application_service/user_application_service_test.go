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
	saveUserEntity := user2.NewUserIndividualEntity()
	saveUserEntity.SetName("個人太郎")

	now := time.Now()
	returnUserEntity := user2.NewUserIndividualEntity()
	returnUserEntity.LoadData(1, "個人太郎", now, now)

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

	userId, err := userApp.RegisterUserIndividual("個人太郎")
	assert.NoError(t, err)
	assert.Equal(t, 1, userId)

	// IDでuser情報を取得してチェックする
	//user, err := userApp.GetUserIndividual(userId)
	//assert.NoError(t, err)
	//assert.Equal(t, userId, user.Id)
	//assert.Equal(t, "個人太郎", user.Name)
	//assert.NotEqual(t, time.Time{}, user.CreatedAt)
	//assert.NotEqual(t, time.Time{}, user.UpdatedAt)

}

//ctrl := gomock.NewController(t)
//
//// Assert that Bar() is invoked.
//defer ctrl.Finish()
//
//m := NewMockFoo(ctrl)
//
//// Asserts that the first and only call to Bar() is passed 99.
//// Anything else will fail.
//m.
//EXPECT().
//Bar(gomock.Eq(99)).
//Return(101)
//
//SUT(m)

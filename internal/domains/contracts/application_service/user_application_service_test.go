package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	user2 "github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user/values"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
	"time"
)

// 個人顧客情報の登録とデータ取得のテスト
func TestUserApplicationService_RegisterUserIndividual(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
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
		userApp := NewUserApplicationServiceWithMock(userRepositoryMock)

		registerdUser, validErrs, err := userApp.RegisterUserIndividual("個人太郎")
		assert.Len(t, validErrs, 0)
		assert.NoError(t, err)
		assert.Equal(t, 1, registerdUser.Id)
		assert.Equal(t, "個人太郎", registerdUser.Name)
		assert.Equal(t, now, registerdUser.CreatedAt)
		assert.Equal(t, now, registerdUser.UpdatedAt)
	})

	t.Run("バリデーションエラー　名前がから文字", func(t *testing.T) {
		// モックリポジトリ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)

		// インスタンス化
		userApp := NewUserApplicationServiceWithMock(userRepositoryMock)

		_, validErrs, err := userApp.RegisterUserIndividual("")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.IsType(t, values.EmptyValidError{}, validErrs[0])
	})

	t.Run("バリデーションエラー　名前が50文字以上", func(t *testing.T) {
		// モックリポジトリ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)

		// インスタンス化
		userApp := NewUserApplicationServiceWithMock(userRepositoryMock)

		_, validErrs, err := userApp.RegisterUserIndividual("000000000011111111112222222222333333333344444444445")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.IsType(t, values.OverLengthValidError{}, validErrs[0])
	})
}

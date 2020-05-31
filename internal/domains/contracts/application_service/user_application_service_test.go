package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
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

// 個人顧客情報の取得のテスト
func TestUserApplicationService_GetUserIndividual(t *testing.T) {
	// userリポジトリモック
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)

	t.Run("データがある時", func(t *testing.T) {
		// GetUserIndividualById()が返却するデータを定義
		returnUserEntity, err := user2.NewUserIndividualEntityWithData(
			1,
			"個人たろう",
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		)
		userRepositoryMock.EXPECT().
			GetUserIndividualById(1, gomock.Any()).
			Return(returnUserEntity, nil).
			Times(1)
		userAppService := NewUserApplicationServiceWithMock(userRepositoryMock)

		userData, err := userAppService.GetUserIndividual(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, userData.Id)
		assert.Equal(t, "個人たろう", userData.Name)
		assert.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), userData.CreatedAt)
		assert.Equal(t, time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC), userData.UpdatedAt)
	})

	t.Run("データがない時", func(t *testing.T) {
		// GetUserIndividualById()が返却するデータを定義
		userRepositoryMock.EXPECT().
			GetUserIndividualById(10000, gomock.Any()).
			Return(nil, nil). // データが無い時はnilが返る
			Times(1)
		userAppService := NewUserApplicationServiceWithMock(userRepositoryMock)
		userData, err := userAppService.GetUserIndividual(10000)
		assert.NoError(t, err)
		assert.Equal(t, data_transfer_objects.UserIndividualDto{}, userData)
	})
}

func TestUserApplicationService_RegisterUserCorporation(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		// リポジトリのSaveUserIndividual()が受け取る引数を用意
		//saveUserEntity, err := user2.NewUserIndividualEntity("個人太郎")
		//assert.NoError(t, err)
		//
		//now := time.Now()
		//returnUserEntity, err := user2.NewUserIndividualEntity("既存太郎")
		//assert.NoError(t, err)
		//err = returnUserEntity.LoadData(1, "個人太郎", now, now)
		//assert.NoError(t, err)
		//
		// モックリポジトリ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)
		//userRepositoryMock.EXPECT().
		//	SaveUserIndividual(
		//		saveUserEntity,
		//		gomock.AssignableToTypeOf(&gorp.Transaction{}),
		//	).Return(returnUserEntity, nil).
		//	Times(1)

		// インスタンス化
		userApp := NewUserApplicationServiceWithMock(userRepositoryMock)

		_, validErrs, err := userApp.RegisterUserCorporation("個人太郎")
		assert.Len(t, validErrs, 0)
		assert.NoError(t, err)
		//assert.Equal(t, 1, registeredUser.Id)
		//assert.Equal(t, "個人太郎", registeredUser.ContactPersonName)
		//assert.Equal(t, "社長太郎", registeredUser.PresidentName)
		//assert.Equal(t, now, registeredUser.CreatedAt)
		//assert.Equal(t, now, registeredUser.UpdatedAt)
	})

	t.Run("バリデーションエラー　名前がから文字", func(t *testing.T) {
		// モックリポジトリ作成
		//ctrl := gomock.NewController(t)
		//defer ctrl.Finish()
		//userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)
		//
		//// インスタンス化
		//userApp := NewUserApplicationServiceWithMock(userRepositoryMock)
		//
		//_, validErrs, err := userApp.RegisterUserIndividual("")
		//assert.NoError(t, err)
		//assert.Len(t, validErrs, 1)
		//assert.IsType(t, values.EmptyValidError{}, validErrs[0])
	})

	t.Run("バリデーションエラー　名前が50文字以上", func(t *testing.T) {
		//// モックリポジトリ作成
		//ctrl := gomock.NewController(t)
		//defer ctrl.Finish()
		//userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)
		//
		//// インスタンス化
		//userApp := NewUserApplicationServiceWithMock(userRepositoryMock)
		//
		//_, validErrs, err := userApp.RegisterUserIndividual("000000000011111111112222222222333333333344444444445")
		//assert.NoError(t, err)
		//assert.Len(t, validErrs, 1)
		//assert.IsType(t, values.OverLengthValidError{}, validErrs[0])
	})
}

package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
	"time"
)

// 個人顧客情報の登録とデータ取得のテスト
func TestUserApplicationService_RegisterUserIndividual(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		// リポジトリのSaveUserIndividual()が受け取る引数を用意
		saveUserEntity, err := entities.NewUserIndividualEntity("個人太郎")
		assert.NoError(t, err)

		now := time.Now()
		returnUserEntity, err := entities.NewUserIndividualEntity("既存太郎")
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
			).Return(1, nil).
			Times(1)
		userRepositoryMock.EXPECT().
			GetUserIndividualById(
				1,
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

		expectValidErrs := map[string][]string{
			"name": []string{
				"空です",
			},
		}

		assert.Equal(t, expectValidErrs, validErrs)
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
		expectValidErrs := map[string][]string{
			"name": []string{
				"50文字より多いです",
			},
		}

		assert.Equal(t, expectValidErrs, validErrs)
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
		returnUserEntity, err := entities.NewUserIndividualEntityWithData(
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
		now := time.Now()
		returnUserEntity, err := entities.NewUserCorporationEntityWithData(1, "担当太郎", "社長太郎", now, now)
		assert.NoError(t, err)

		// モックリポジトリ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)
		userRepositoryMock.EXPECT().
			SaveUserCorporation(
				gomock.Any(),
				gomock.AssignableToTypeOf(&gorp.Transaction{}),
			).Return(1, nil).
			Times(1)
		userRepositoryMock.EXPECT().
			GetUserCorporationById(
				1,
				gomock.AssignableToTypeOf(&gorp.Transaction{}),
			).Return(returnUserEntity, nil).
			Times(1)

		// インスタンス化
		userApp := NewUserApplicationServiceWithMock(userRepositoryMock)

		registeredUser, validErrs, err := userApp.RegisterUserCorporation("担当太郎", "社長太郎")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 0)
		assert.Equal(t, 1, registeredUser.Id)
		assert.Equal(t, "担当太郎", registeredUser.ContactPersonName)
		assert.Equal(t, "社長太郎", registeredUser.PresidentName)
		assert.Equal(t, now, registeredUser.CreatedAt)
		assert.Equal(t, now, registeredUser.UpdatedAt)
	})

	t.Run("バリデーションエラー　担当者名と社長名がから文字", func(t *testing.T) {
		// モックリポジトリ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)

		// インスタンス化
		userApp := NewUserApplicationServiceWithMock(userRepositoryMock)

		_, validErrs, err := userApp.RegisterUserCorporation("", "")
		assert.NoError(t, err)
		expectValidErrs := map[string][]string{
			"contact_person_name": []string{
				"空です",
			},
			"president_name": []string{
				"空です",
			},
		}

		assert.Equal(t, expectValidErrs, validErrs)
	})

	t.Run("バリデーションエラー　名前が50文字以上", func(t *testing.T) {
		// モックリポジトリ作成
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)

		// インスタンス化
		userApp := NewUserApplicationServiceWithMock(userRepositoryMock)

		_, validErrs, err := userApp.RegisterUserCorporation("000000000011111111112222222222333333333344444444445", "000000000011111111112222222222333333333344444444445")
		assert.NoError(t, err)
		expectValidErrs := map[string][]string{
			"contact_person_name": []string{
				"50文字より多いです",
			},
			"president_name": []string{
				"50文字より多いです",
			},
		}

		assert.Equal(t, expectValidErrs, validErrs)
	})
}
func TestUserApplicationService_GetUserCorporation(t *testing.T) {
	// userリポジトリモック
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)

	t.Run("データがある時", func(t *testing.T) {
		// GetUserIndividualById()が返却するデータを定義
		returnUserEntity, err := entities.NewUserCorporationEntityWithData(
			1,
			"担当たろう",
			"社長たろう",
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		)
		assert.NoError(t, err)

		userRepositoryMock.EXPECT().
			GetUserCorporationById(1, gomock.Any()).
			Return(returnUserEntity, nil).
			Times(1)
		userAppService := NewUserApplicationServiceWithMock(userRepositoryMock)

		userData, err := userAppService.GetUserCorporation(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, userData.Id)
		assert.Equal(t, "担当たろう", userData.ContactPersonName)
		assert.Equal(t, "社長たろう", userData.PresidentName)
		assert.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), userData.CreatedAt)
		assert.Equal(t, time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC), userData.UpdatedAt)
	})

	t.Run("データがない時", func(t *testing.T) {
		// GetUserIndividualById()が返却するデータを定義
		userRepositoryMock.EXPECT().
			GetUserCorporationById(10000, gomock.Any()).
			Return(nil, nil). // データが無い時はnilが返る
			Times(1)
		userAppService := NewUserApplicationServiceWithMock(userRepositoryMock)
		userData, err := userAppService.GetUserCorporation(10000)
		assert.NoError(t, err)
		assert.Equal(t, data_transfer_objects.UserCorporationDto{}, userData)
	})
}
func TestUserApplicationService_GetUserById(t *testing.T) {
	// 個人顧客と法人顧客データを登録
	app := NewUserApplicationService()
	individualDto, validErrors, err := app.RegisterUserIndividual("個人顧客取得テスト")
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)
	corporationDto, validErrors, err := app.RegisterUserCorporation("法人顧客取得テスト担当", "法人顧客取得テスト社長")
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)

	t.Run("個人顧客", func(t *testing.T) {
		t.Run("データがある時", func(t *testing.T) {
			user, err := app.GetUserById(individualDto.Id)
			assert.NoError(t, err)
			userDto, ok := user.(data_transfer_objects.UserIndividualDto)
			assert.True(t, ok)
			assert.NotZero(t, userDto.Id)
			assert.Equal(t, "個人顧客取得テスト", userDto.Name)
			assert.NotZero(t, userDto.CreatedAt)
			assert.NotZero(t, userDto.UpdatedAt)
		})

		t.Run("データが無いとき", func(t *testing.T) {
			user, err := app.GetUserById(-100)
			assert.NoError(t, err)
			assert.Nil(t, user)
		})
	})

	t.Run("法人顧客", func(t *testing.T) {
		t.Run("データがある時", func(t *testing.T) {
			user, err := app.GetUserById(corporationDto.Id)
			assert.NoError(t, err)
			userDto, ok := user.(data_transfer_objects.UserCorporationDto)
			assert.True(t, ok)
			assert.NotZero(t, userDto.Id)
			assert.Equal(t, "法人顧客取得テスト担当", userDto.ContactPersonName)
			assert.Equal(t, "法人顧客取得テスト社長", userDto.PresidentName)
			assert.NotZero(t, userDto.CreatedAt)
			assert.NotZero(t, userDto.UpdatedAt)

		})

		t.Run("データが無いとき", func(t *testing.T) {
			user, err := app.GetUserById(-100)
			assert.NoError(t, err)
			assert.Nil(t, user)
		})
	})
}

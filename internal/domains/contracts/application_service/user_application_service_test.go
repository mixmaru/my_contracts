package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 個人顧客情報の登録とデータ取得のテスト
func TestUserApplicationService_RegisterUserIndividual(t *testing.T) {
	userApp := NewUserApplicationService()
	t.Run("名前を渡すと個人顧客データを登録できる", func(t *testing.T) {
		registeredUser, validErrs, err := userApp.RegisterUserIndividual("個人太郎")
		assert.Len(t, validErrs, 0)
		assert.NoError(t, err)
		assert.NotZero(t, registeredUser.Id)
		assert.Equal(t, "個人太郎", registeredUser.Name)
		assert.NotZero(t, registeredUser.CreatedAt)
		assert.NotZero(t, registeredUser.UpdatedAt)
	})

	t.Run("バリデーションエラー　名前がから文字のときvalidErrsが返る", func(t *testing.T) {
		_, validErrs, err := userApp.RegisterUserIndividual("")
		assert.NoError(t, err)

		expectValidErrs := map[string][]string{
			"name": []string{
				"空です",
			},
		}

		assert.Equal(t, expectValidErrs, validErrs)
	})

	t.Run("バリデーションエラー　名前が50文字以上のときvalidErrorが返る", func(t *testing.T) {
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

func TestUserApplicationService_RegisterUserCorporation(t *testing.T) {
	userApp := NewUserApplicationService()
	t.Run("会社名と担当者名と社長名を渡すと法人顧客データが登録できる", func(t *testing.T) {
		registeredUser, validErrs, err := userApp.RegisterUserCorporation("イケてる会社", "担当太郎", "社長太郎")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 0)

		assert.NotZero(t, registeredUser.Id)
		assert.Equal(t, "イケてる会社", registeredUser.CorporationName)
		assert.Equal(t, "担当太郎", registeredUser.ContactPersonName)
		assert.Equal(t, "社長太郎", registeredUser.PresidentName)
		assert.NotZero(t, registeredUser.CreatedAt)
		assert.NotZero(t, registeredUser.UpdatedAt)
	})

	t.Run("バリデーションエラー　会社名と担当者名と社長名がから文字だったらバリデーションエラーが返る", func(t *testing.T) {
		_, validErrs, err := userApp.RegisterUserCorporation("", "", "")
		assert.NoError(t, err)
		expectValidErrs := map[string][]string{
			"corporation_name": []string{
				"空です",
			},
			"contact_person_name": []string{
				"空です",
			},
			"president_name": []string{
				"空です",
			},
		}

		assert.Equal(t, expectValidErrs, validErrs)
	})

	t.Run("バリデーションエラー　名前が50文字以上だったらバリデーションエラーメッセージが返る", func(t *testing.T) {
		_, validErrs, err := userApp.RegisterUserCorporation("000000000011111111112222222222333333333344444444445", "000000000011111111112222222222333333333344444444445", "000000000011111111112222222222333333333344444444445")
		assert.NoError(t, err)
		expectValidErrs := map[string][]string{
			"corporation_name": []string{
				"50文字より多いです",
			},
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

func TestUserApplicationService_GetUserById(t *testing.T) {
	// 個人顧客と法人顧客データを登録
	app := NewUserApplicationService()
	individualDto, validErrors, err := app.RegisterUserIndividual("個人顧客取得テスト")
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)
	corporationDto, validErrors, err := app.RegisterUserCorporation("法人顧客会社名", "法人顧客取得テスト担当", "法人顧客取得テスト社長")
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)

	t.Run("個人顧客", func(t *testing.T) {
		t.Run("データがある時はidでデータ取得ができる", func(t *testing.T) {
			user, err := app.GetUserById(individualDto.Id)
			assert.NoError(t, err)
			userDto, ok := user.(data_transfer_objects.UserIndividualDto)
			assert.True(t, ok)
			assert.NotZero(t, userDto.Id)
			assert.Equal(t, "個人顧客取得テスト", userDto.Name)
			assert.NotZero(t, userDto.CreatedAt)
			assert.NotZero(t, userDto.UpdatedAt)
		})

		t.Run("データが無いときはnilが返る", func(t *testing.T) {
			user, err := app.GetUserById(-100)
			assert.NoError(t, err)
			assert.Nil(t, user)
		})
	})

	t.Run("法人顧客", func(t *testing.T) {
		t.Run("データがある時はidでデータが取得できる", func(t *testing.T) {
			user, err := app.GetUserById(corporationDto.Id)
			assert.NoError(t, err)
			userDto, ok := user.(data_transfer_objects.UserCorporationDto)
			assert.True(t, ok)
			assert.NotZero(t, userDto.Id)
			assert.Equal(t, "法人顧客会社名", userDto.CorporationName)
			assert.Equal(t, "法人顧客取得テスト担当", userDto.ContactPersonName)
			assert.Equal(t, "法人顧客取得テスト社長", userDto.PresidentName)
			assert.NotZero(t, userDto.CreatedAt)
			assert.NotZero(t, userDto.UpdatedAt)
		})

		t.Run("データが無いときはnilが返る", func(t *testing.T) {
			user, err := app.GetUserById(-100)
			assert.NoError(t, err)
			assert.Nil(t, user)
		})
	})
}

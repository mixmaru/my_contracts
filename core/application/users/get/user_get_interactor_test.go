package get

import (
	"github.com/mixmaru/my_contracts/core/application/users"
	"github.com/mixmaru/my_contracts/core/application/users/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserGetInteractor_GetUserById(t *testing.T) {
	// 個人顧客と法人顧客データを登録
	userRep := db.NewUserRepository()
	createIntaractor := create.NewUserIndividualCreateInteractor(userRep)
	createResponse, err := createIntaractor.Handle(create.NewUserIndividualCreateUseCaseRequest("個人顧客取得テスト"))
	assert.NoError(t, err)
	assert.Len(t, createResponse.ValidationErrors, 0)
	//corporationDto, validErrors, err := createIntaractor.RegisterUserCorporation("法人顧客会社名", "法人顧客取得テスト担当", "法人顧客取得テスト社長")
	//assert.NoError(t, err)
	//assert.Len(t, validErrors, 0)

	t.Run("個人顧客", func(t *testing.T) {
		getIntaractor := NewUserGetInteractor(userRep)
		t.Run("データがある時はidでデータ取得ができる", func(t *testing.T) {
			response, err := getIntaractor.Handle(NewUserGetUseCaseRequest(createResponse.UserDto.Id))
			assert.NoError(t, err)
			userDto, ok := response.UserDto.(users.UserIndividualDto)
			assert.True(t, ok)
			assert.NotZero(t, userDto.Id)
			assert.Equal(t, "個人顧客取得テスト", userDto.Name)
			assert.NotZero(t, userDto.CreatedAt)
			assert.NotZero(t, userDto.UpdatedAt)
		})

		t.Run("データが無いときはnilが返る", func(t *testing.T) {
			response, err := getIntaractor.Handle(NewUserGetUseCaseRequest(-100))
			assert.NoError(t, err)
			assert.Nil(t, response.UserDto)
		})
	})

	//t.Run("法人顧客", func(t *testing.T) {
	//	t.Run("データがある時はidでデータが取得できる", func(t *testing.T) {
	//		user, err := createIntaractor.GetUserById(corporationDto.Id)
	//		assert.NoError(t, err)
	//		userDto, ok := user.(data_transfer_objects.UserCorporationDto)
	//		assert.True(t, ok)
	//		assert.NotZero(t, userDto.Id)
	//		assert.Equal(t, "法人顧客会社名", userDto.CorporationName)
	//		assert.Equal(t, "法人顧客取得テスト担当", userDto.ContactPersonName)
	//		assert.Equal(t, "法人顧客取得テスト社長", userDto.PresidentName)
	//		assert.NotZero(t, userDto.CreatedAt)
	//		assert.NotZero(t, userDto.UpdatedAt)
	//	})
	//
	//	t.Run("データが無いときはnilが返る", func(t *testing.T) {
	//		user, err := createIntaractor.GetUserById(-100)
	//		assert.NoError(t, err)
	//		assert.Nil(t, user)
	//	})
	//})
}

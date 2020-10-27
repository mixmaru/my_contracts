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
	createIndividualInteractor := create.NewUserIndividualCreateInteractor(userRep)
	createIndividualResponse, err := createIndividualInteractor.Handle(create.NewUserIndividualCreateUseCaseRequest("個人顧客取得テスト"))
	assert.NoError(t, err)
	assert.Len(t, createIndividualResponse.ValidationErrors, 0)
	createCorporationInteractor := create.NewUserCorporationCreateInteractor(userRep)
	createCorporationResponse, err := createCorporationInteractor.Handle(create.NewUserCorporationCreateUseCaseRequest("法人顧客会社名", "法人顧客取得テスト担当", "法人顧客取得テスト社長"))
	assert.NoError(t, err)
	assert.Len(t, createCorporationResponse.ValidationErrors, 0)

	t.Run("個人顧客", func(t *testing.T) {
		getInteractor := NewUserGetInteractor(userRep)
		t.Run("データがある時はidでデータ取得ができる", func(t *testing.T) {
			response, err := getInteractor.Handle(NewUserGetUseCaseRequest(createIndividualResponse.UserDto.Id))
			assert.NoError(t, err)
			userDto, ok := response.UserDto.(users.UserIndividualDto)
			assert.True(t, ok)
			assert.NotZero(t, userDto.Id)
			assert.Equal(t, "個人顧客取得テスト", userDto.Name)
			assert.NotZero(t, userDto.CreatedAt)
			assert.NotZero(t, userDto.UpdatedAt)
		})

		t.Run("データが無いときはnilが返る", func(t *testing.T) {
			response, err := getInteractor.Handle(NewUserGetUseCaseRequest(-100))
			assert.NoError(t, err)
			assert.Nil(t, response.UserDto)
		})
	})

	t.Run("法人顧客", func(t *testing.T) {
		getInteractor := NewUserGetInteractor(userRep)
		t.Run("データがある時はidでデータが取得できる", func(t *testing.T) {
			response, err := getInteractor.Handle(NewUserGetUseCaseRequest(createCorporationResponse.UserDto.Id))
			assert.NoError(t, err)
			userDto, ok := response.UserDto.(users.UserCorporationDto)
			assert.True(t, ok)
			assert.NotZero(t, userDto.Id)
			assert.Equal(t, "法人顧客会社名", userDto.CorporationName)
			assert.Equal(t, "法人顧客取得テスト担当", userDto.ContactPersonName)
			assert.Equal(t, "法人顧客取得テスト社長", userDto.PresidentName)
			assert.NotZero(t, userDto.CreatedAt)
			assert.NotZero(t, userDto.UpdatedAt)
		})

		t.Run("データが無いときはnilが返る", func(t *testing.T) {
			response, err := getInteractor.Handle(NewUserGetUseCaseRequest(-100))
			assert.NoError(t, err)
			assert.Nil(t, response.UserDto)
		})
	})
}

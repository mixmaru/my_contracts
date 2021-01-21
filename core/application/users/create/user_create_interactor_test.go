package create

import (
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 個人顧客情報の登録とデータ取得のテスト
func TestUserApplicationService_RegisterUserIndividual(t *testing.T) {
	interactor := NewUserIndividualCreateInteractor(db.NewUserRepository())
	t.Run("名前を渡すと個人顧客データを登録できる", func(t *testing.T) {
		response, err := interactor.Handle(NewUserIndividualCreateUseCaseRequest("個人太郎"))
		assert.NoError(t, err)
		assert.Len(t, response.ValidationErrors, 0)
		assert.NotZero(t, response.UserDto.Id)
		assert.Equal(t, "個人太郎", response.UserDto.Name)
		assert.NotZero(t, response.UserDto.CreatedAt)
		assert.NotZero(t, response.UserDto.UpdatedAt)
	})

	t.Run("バリデーションエラー　名前がから文字のときvalidErrsが返る", func(t *testing.T) {
		response, err := interactor.Handle(NewUserIndividualCreateUseCaseRequest(""))
		assert.NoError(t, err)

		expectValidErrs := map[string][]string{
			"name": []string{
				"空です",
			},
		}

		assert.Equal(t, expectValidErrs, response.ValidationErrors)
	})

	t.Run("バリデーションエラー　名前が50文字以上のときvalidErrorが返る", func(t *testing.T) {
		response, err := interactor.Handle(NewUserIndividualCreateUseCaseRequest("000000000011111111112222222222333333333344444444445"))
		assert.NoError(t, err)
		expectValidErrs := map[string][]string{
			"name": []string{
				"50文字より多いです",
			},
		}

		assert.Equal(t, expectValidErrs, response.ValidationErrors)
	})
}

func TestUserApplicationService_RegisterUserCorporation(t *testing.T) {
	interactor := NewUserCorporationCreateInteractor(db.NewUserRepository())
	t.Run("会社名と担当者名と社長名を渡すと法人顧客データが登録できる", func(t *testing.T) {
		response, err := interactor.Handle(NewUserCorporationCreateUseCaseRequest("イケてる会社", "担当太郎", "社長太郎"))
		assert.NoError(t, err)
		assert.Len(t, response.ValidationErrors, 0)

		dto := response.UserDto
		assert.NotZero(t, dto.Id)
		assert.Equal(t, "イケてる会社", dto.CorporationName)
		assert.Equal(t, "担当太郎", dto.ContactPersonName)
		assert.Equal(t, "社長太郎", dto.PresidentName)
		assert.NotZero(t, dto.CreatedAt)
		assert.NotZero(t, dto.UpdatedAt)
	})

	t.Run("バリデーションエラー　会社名と担当者名と社長名がから文字だったらバリデーションエラーが返る", func(t *testing.T) {
		response, err := interactor.Handle(NewUserCorporationCreateUseCaseRequest("", "", ""))
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

		assert.Equal(t, expectValidErrs, response.ValidationErrors)
	})

	t.Run("バリデーションエラー　名前が50文字以上だったらバリデーションエラーメッセージが返る", func(t *testing.T) {
		response, err := interactor.Handle(NewUserCorporationCreateUseCaseRequest("000000000011111111112222222222333333333344444444445", "000000000011111111112222222222333333333344444444445", "000000000011111111112222222222333333333344444444445"))
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

		assert.Equal(t, expectValidErrs, response.ValidationErrors)
	})
}

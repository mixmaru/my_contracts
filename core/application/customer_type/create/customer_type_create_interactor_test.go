package create

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestCustomerPropertyTypeCreateInteractor_Register(t *testing.T) {
	interactor := NewCustomerTypeCreateInteractor(db.NewCustomerTypeRepository(), db.NewCustomerPropertyTypeRepository())

	t.Run("名前とカスタマープロパティタイプIDを渡すとカスタマータイプが登録される", func(t *testing.T) {
		////// 準備
		// カスタマープロパティタイプの登録
		propertyIds, propertyDtos, err := preInsertCustomerProperty()
		assert.NoError(t, err)

		////// 実行
		timestampstr := utils.CreateTimestampString()
		request := NewCustomerTypeCreateUseCaseRequest("お得意様"+timestampstr, propertyIds)
		response, err := interactor.Handle(request)
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, response.ValidationErrors, 0)
		assert.NotZero(t, response.CustomerTypeDto.Id)
		assert.Equal(t, "お得意様"+timestampstr, response.CustomerTypeDto.Name)
		assert.Equal(t, propertyDtos, response.CustomerTypeDto.CustomerPropertyTypes)

		t.Run("既に登録されているカスタマー名だった場合はバリデーションエラーになる", func(t *testing.T) {
			request := NewCustomerTypeCreateUseCaseRequest("お得意様"+timestampstr, propertyIds)
			response, err := interactor.Handle(request)
			assert.NoError(t, err)
			expect := map[string][]string{
				"name": []string{
					"既に存在する名前です",
				},
			}

			assert.Len(t, response.ValidationErrors, 1)
			assert.Equal(t, expect, response.ValidationErrors)
			assert.Zero(t, response.CustomerTypeDto.Id)
		})

		t.Run("存在しないpropertyIdが指定されていた場合バリデーションエラーになる", func(t *testing.T) {
			request := NewCustomerTypeCreateUseCaseRequest("お得意様プロパティバリデーションエラー"+timestampstr, []int{-1000, -20000})
			response, err := interactor.Handle(request)
			assert.NoError(t, err)
			expect := map[string][]string{
				"customer_property_ids": []string{
					"-1000, -20000 は存在しないidです",
				},
			}

			assert.Len(t, response.ValidationErrors, 1)
			assert.Equal(t, expect, response.ValidationErrors)
			assert.Zero(t, response.CustomerTypeDto.Id)
		})
	})
}

func preInsertCustomerProperty() ([]int, []customer_property_type.CustomerPropertyTypeDto, error) {
	var retInts []int
	var retDto []customer_property_type.CustomerPropertyTypeDto

	timestampstr := utils.CreateTimestampString()
	propertyInteractor := create.NewCustomerPropertyTypeCreateInteractor(db.NewCustomerPropertyTypeRepository())

	for i := 0; i < 2; i++ {
		request := create.NewCustomerPropertyTypeCreateUseCaseRequest("好きな食べ物"+strconv.Itoa(i+1)+timestampstr, "string")
		response, err := propertyInteractor.Handle(request)
		if err != nil {
			return nil, nil, err
		}
		if len(response.ValidationErrors) > 0 {
			return nil, nil, errors.Errorf("バリデーションエラー。%+v", response.ValidationErrors)
		}
		retInts = append(retInts, response.CustomerPropertyTypeDto.Id)
		retDto = append(retDto, response.CustomerPropertyTypeDto)
	}

	return retInts, retDto, nil
}

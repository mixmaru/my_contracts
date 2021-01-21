package get_by_id

import (
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/test_utils"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomerGetByIdInteractor_Handle(t *testing.T) {
	////// 事前準備
	// カスタマータイプ作成
	timestampStr := utils.CreateTimestampString()
	customerTypeDto, err := test_utils.PreCreateCustomerPropertyTypeAndCustomerType("お得意様"+timestampStr, []test_utils.PropertyParam{
		{
			PropertyTypeName: "性別" + timestampStr,
			PropertyType:     test_utils.PROPERTY_TYPE_STRING,
		},
		{
			PropertyTypeName: "年齢" + timestampStr,
			PropertyType:     test_utils.PROPERTY_TYPE_NUMERIC,
		},
	})
	assert.NoError(t, err)
	// カスタマー登録
	customerDto, err := test_utils.PreCreateCustomer("山田邦明", customerTypeDto.Id, map[int]interface{}{
		customerTypeDto.CustomerPropertyTypes[0].Id: "男",
		customerTypeDto.CustomerPropertyTypes[1].Id: 32.,
	})
	assert.NoError(t, err)
	assert.NotZero(t, customerDto)

	t.Run("登録されているCustomerのIdを渡すとデータが取得できる", func(t *testing.T) {
		////// 実行
		request := NewCustomerGetByIdUseCaseRequest(customerDto.Id)
		interactor := NewCustomerGetByIdInteractor(db.NewCustomerRepository())
		response, err := interactor.Handle(request)
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, response.ValidationErrors, 0)
		assert.Equal(t, customerDto, response.CustomerDto)
	})

	t.Run("存在しないCustomerIdを渡すとゼロ値が返る", func(t *testing.T) {

	})
}

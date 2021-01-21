package get_by_id

import (
	"github.com/mixmaru/my_contracts/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomerGetByIdInteractor_Handle(t *testing.T) {
	////// 事前準備
	// カスタマー登録
	customerDto, err := test_utils.PreCreateCustomer()
	assert.NoError(t, err)
	assert.NotZero(t, customerDto)

	t.Run("登録されているCustomerのIdを渡すとデータが取得できる", func(t *testing.T) {

	})

	t.Run("存在しないCustomerIdを渡すとゼロ値が返る", func(t *testing.T) {

	})
}

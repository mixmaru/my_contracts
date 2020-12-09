package customer

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCustomerType_CustomerParamTypes(t *testing.T) {
	t.Run("すべてのCustomerParamTypeが返ってくる", func(t *testing.T) {
		c := &CustomerTypeEntity{
			id:   1,
			name: "優良顧客",
			customerParamTypes: map[int]*CustomerPropertyTypeEntity{
				1: {
					id:        1,
					name:      "性別",
					paramType: PARAM_TYPE_STRING,
				},
			},
		}
		actual := c.CustomerParamTypes()
		expected := map[int]*CustomerPropertyTypeEntity{
			1: {
				id:        1,
				name:      "性別",
				paramType: PARAM_TYPE_STRING,
			},
		}
		assert.True(t, reflect.DeepEqual(expected, actual))
	})

	t.Run("返っt北CustomerParamTypeを変更しても本体に影響しない", func(t *testing.T) {
		c := &CustomerTypeEntity{
			id:   1,
			name: "優良顧客",
			customerParamTypes: map[int]*CustomerPropertyTypeEntity{
				1: {
					id:        1,
					name:      "性別",
					paramType: PARAM_TYPE_STRING,
				},
			},
		}

		// 取得した値を変更する
		temp := c.CustomerParamTypes()
		temp[1].name = "変更した"

		// 再取得したものに影響しないか確認する
		actual := c.CustomerParamTypes()
		expected := map[int]*CustomerPropertyTypeEntity{
			1: {
				id:        1,
				name:      "性別",
				paramType: PARAM_TYPE_STRING,
			},
		}
		assert.True(t, reflect.DeepEqual(expected, actual))
	})
}

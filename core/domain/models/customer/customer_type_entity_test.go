package customer

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestCustomerType_CustomerParamTypes(t *testing.T) {
	t.Run("すべてのCustomerParamTypeが返ってくる", func(t *testing.T) {
		c := &CustomerType{
			id:   1,
			name: "優良顧客",
			customerParamTypes: map[int]*CustomerParamType{
				1: {
					id:        1,
					name:      "性別",
					paramType: PARAM_TYPE_STRING,
					createdAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			createdAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			updatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		actual := c.CustomerParamTypes()
		expected := map[int]*CustomerParamType{
			1: {
				id:        1,
				name:      "性別",
				paramType: PARAM_TYPE_STRING,
				createdAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				updatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		assert.True(t, reflect.DeepEqual(expected, actual))
	})

	t.Run("返っt北CustomerParamTypeを変更しても本体に影響しない", func(t *testing.T) {
		c := &CustomerType{
			id:   1,
			name: "優良顧客",
			customerParamTypes: map[int]*CustomerParamType{
				1: {
					id:        1,
					name:      "性別",
					paramType: PARAM_TYPE_STRING,
					createdAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			createdAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			updatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		// 取得した値を変更する
		temp := c.CustomerParamTypes()
		temp[1].name = "変更した"

		// 再取得したものに影響しないか確認する
		actual := c.CustomerParamTypes()
		expected := map[int]*CustomerParamType{
			1: {
				id:        1,
				name:      "性別",
				paramType: PARAM_TYPE_STRING,
				createdAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				updatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		assert.True(t, reflect.DeepEqual(expected, actual))
	})
}

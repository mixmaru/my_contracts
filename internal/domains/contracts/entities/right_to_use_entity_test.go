package entities

import (
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRightToUseEntity_NewRightToUseEntity(t *testing.T) {
	t.Run("契約Entityと利用権開始日時を渡すと利用権Entityがインスタンス化できる", func(t *testing.T) {
		rightToUseEntity := NewRightToUseEntity(
			1,
			utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 2, 2, 0, 0, 0, 0),
		)
		assert.Zero(t, rightToUseEntity.Id())
		assert.Equal(t, 1, rightToUseEntity.ContractId())
		assert.True(t, rightToUseEntity.ValidFrom().Equal(utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0)))
		assert.True(t, rightToUseEntity.ValidTo().Equal(utils.CreateJstTime(2020, 2, 2, 0, 0, 0, 0)))
		assert.Zero(t, rightToUseEntity.CreatedAt())
		assert.Zero(t, rightToUseEntity.UpdatedAt())
	})
}

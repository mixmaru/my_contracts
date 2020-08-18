package data_mappers

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRightToUser_NewRightToUserFromEntity(t *testing.T) {
	t.Run("権利エンティティからマッパーを作成する", func(t *testing.T) {
		// 準備
		entity := entities.NewRightToUseWithData(
			1,
			2,
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
		)

		// 実行
		mapper := NewRightToUseMapperFromEntity(entity)

		// 検証
		assert.Equal(t, 1, mapper.Id)
		assert.Equal(t, 2, mapper.ContractId)
		assert.True(t, mapper.ValidFrom.Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
		assert.True(t, mapper.ValidTo.Equal(time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)))
		assert.True(t, mapper.CreatedAt.Equal(time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)))
		assert.True(t, mapper.UpdatedAt.Equal(time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC)))
	})
}

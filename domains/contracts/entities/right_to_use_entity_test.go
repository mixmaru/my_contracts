package entities

import (
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRightToUseEntity_NewRightToUseEntity(t *testing.T) {
	t.Run("契約Entityと利用権開始日時を渡すと利用権Entityがインスタンス化できる", func(t *testing.T) {
		rightToUseEntity := NewRightToUseEntity(
			utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 2, 2, 0, 0, 0, 0),
		)
		assert.Zero(t, rightToUseEntity.Id())
		assert.True(t, rightToUseEntity.ValidFrom().Equal(utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0)))
		assert.True(t, rightToUseEntity.ValidTo().Equal(utils.CreateJstTime(2020, 2, 2, 0, 0, 0, 0)))
		assert.Zero(t, rightToUseEntity.CreatedAt())
		assert.Zero(t, rightToUseEntity.UpdatedAt())
	})
}

func TestRightToUseEntity_NewRightToUseEntityWithData(t *testing.T) {
	t.Run("すべての要素データを読み込ませてインスタンス化する", func(t *testing.T) {
		// 実行
		entity := NewRightToUseEntityWithData(
			1,
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
		)

		// 検証
		assert.Equal(t, 1, entity.Id())
		assert.True(t, entity.ValidFrom().Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
		assert.True(t, entity.ValidTo().Equal(time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)))
		assert.True(t, entity.CreatedAt().Equal(time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)))
		assert.True(t, entity.UpdatedAt().Equal(time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC)))
	})
}

func TestRightToUseEntity_LoadData(t *testing.T) {
	t.Run("要素データを読み込ませて中身を上書きできる", func(t *testing.T) {
		// 準備
		entity := NewRightToUseEntity(time.Time{}, time.Time{})

		// 実行
		entity.LoadData(
			1,
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
		)

		// 検証
		assert.Equal(t, 1, entity.Id())
		assert.True(t, entity.ValidFrom().Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
		assert.True(t, entity.ValidTo().Equal(time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)))
		assert.True(t, entity.CreatedAt().Equal(time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)))
		assert.True(t, entity.UpdatedAt().Equal(time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC)))
	})
}

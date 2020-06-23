package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProductRecord_SetDataToEntity(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		productRecord := ProductRecord{
			Id:    1,
			Name:  "名前",
			Price: decimal.NewFromFloat(1000),
			CreatedAtUpdatedAt: CreatedAtUpdatedAt{
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		entity := &entities.ProductEntity{}

		err := productRecord.SetDataToEntity(entity)
		assert.NoError(t, err)
		assert.Equal(t, 1, entity.Id())
		assert.Equal(t, "名前", entity.Name())
		assert.Equal(t, decimal.NewFromFloat(1000), entity.Price())
		assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), entity.CreatedAt())
		assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), entity.UpdatedAt())
	})

	t.Run("違うentityが渡されたとき", func(t *testing.T) {
		productRecord := ProductRecord{
			Id:    1,
			Name:  "名前",
			Price: decimal.NewFromFloat(1000),
			CreatedAtUpdatedAt: CreatedAtUpdatedAt{
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		entity := &entities.UserEntity{}

		err := productRecord.SetDataToEntity(entity)
		assert.Error(t, err)
	})
}

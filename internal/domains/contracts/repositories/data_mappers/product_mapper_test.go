package data_mappers

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProductRecord_SetDataToEntity(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		productRecord := ProductMapper{
			Id:   1,
			Name: "名前",
			CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		entity := &entities.ProductEntity{}

		err := productRecord.SetDataToEntity(entity)
		assert.NoError(t, err)
		assert.Equal(t, 1, entity.Id())
		assert.Equal(t, "名前", entity.Name())

		price, exist := entity.MonthlyPrice()
		assert.True(t, exist)
		assert.Equal(t, "1000", price.String())
		assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), entity.CreatedAt())
		assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), entity.UpdatedAt())
	})

	t.Run("違うentityが渡されたとき", func(t *testing.T) {
		productRecord := ProductMapper{
			Id:   1,
			Name: "名前",
			CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		entity := &entities.UserEntity{}

		err := productRecord.SetDataToEntity(entity)
		assert.Error(t, err)
	})
}

package data_mappers

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestContractRecord_SetDataToEntity(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		contractMapper := ContractMapper{
			Id:        1,
			UserId:    2,
			ProductId: 3,
			CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		entity := &entities.ContractEntity{}

		err := contractMapper.SetDataToEntity(entity)
		assert.NoError(t, err)
		assert.Equal(t, 1, entity.Id())
		assert.Equal(t, 2, entity.UserId())
		assert.Equal(t, 3, entity.ProductId())

		assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), entity.CreatedAt())
		assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), entity.UpdatedAt())
	})

	t.Run("違うentityが渡されたとき", func(t *testing.T) {
		contractMapper := ContractMapper{
			Id:        1,
			UserId:    2,
			ProductId: 3,
			CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
				CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		entity := &entities.UserEntity{}

		err := contractMapper.SetDataToEntity(entity)
		assert.Error(t, err)
	})
}

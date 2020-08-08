package data_mappers

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestContractRecord_SetDataToEntity(t *testing.T) {
	t.Run("ContractEntityを渡すと_Entityにmapperが持っているデータをセットする", func(t *testing.T) {
		contractMapper := ContractMapper{
			Id:               1,
			UserId:           2,
			ContractDate:     utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			BillingStartDate: utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0),
			ProductId:        3,

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
		assert.EqualValues(t, utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), entity.ContractDate())
		assert.EqualValues(t, utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0), entity.BillingStartDate())

		assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), entity.CreatedAt())
		assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), entity.UpdatedAt())
	})

	t.Run("違う型のentityが渡されたときは_エラーとなる", func(t *testing.T) {
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

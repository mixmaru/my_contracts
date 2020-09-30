package entities

import (
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// UserIndividualのインスタンス化をテスト
func TestContractEntity_NewContractEntity(t *testing.T) {
	// インスタンス化
	entity := NewContractEntity(
		1,
		2,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
	)

	// テスト
	assert.Equal(t, 1, entity.UserId())
	assert.Equal(t, 2, entity.ProductId())
	assert.EqualValues(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), entity.ContractDate())
	assert.EqualValues(t, time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), entity.BillingStartDate())
}

func TestContractEntity_NewContractEntityWithData(t *testing.T) {
	// インスタンス化
	contractDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	billingStartDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	entity, err := NewContractEntityWithData(1, 2, 3, contractDate, billingStartDate, createdAt, updatedAt)
	assert.NoError(t, err)

	assert.Equal(t, 1, entity.Id())
	assert.Equal(t, 2, entity.UserId())
	assert.Equal(t, 3, entity.ProductId())
	assert.True(t, contractDate.Equal(entity.ContractDate()))
	assert.True(t, billingStartDate.Equal(entity.BillingStartDate()))
	assert.True(t, createdAt.Equal(entity.CreatedAt()))
	assert.True(t, updatedAt.Equal(entity.UpdatedAt()))
}

func TestContractEntity_LoadData(t *testing.T) {
	t.Run("プライベートプロパティに値をセットすることができる", func(t *testing.T) {
		contractEntity := &ContractEntity{}

		contractDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
		billingStartDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		updateAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		err := contractEntity.LoadData(
			1,
			2,
			3,
			contractDate,
			billingStartDate,
			createdAt,
			updateAt,
		)
		assert.NoError(t, err)

		assert.Equal(t, 1, contractEntity.Id())
		assert.Equal(t, 2, contractEntity.UserId())
		assert.Equal(t, 3, contractEntity.ProductId())
		assert.EqualValues(t, contractDate, contractEntity.ContractDate())
		assert.EqualValues(t, billingStartDate, contractEntity.BillingStartDate())
		assert.EqualValues(t, createdAt, contractEntity.CreatedAt())
		assert.EqualValues(t, updateAt, contractEntity.UpdatedAt())
	})
}

func TestContractEntity_LastBillingStartDate(t *testing.T) {
	contract, err := NewContractEntityWithData(
		1,
		2,
		3,
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
	)
	assert.NoError(t, err)

	t.Run("2020_01_02を渡すと直近の課金開始日_2020_01_02が返る", func(t *testing.T) {
		actual := contract.LastBillingStartDate(utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0))
		assert.True(t, actual.Equal(utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0)))
	})
	t.Run("2020_02_01を渡すと直近の課金開始日_2020_01_02が返る", func(t *testing.T) {
		actual := contract.LastBillingStartDate(utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0))
		assert.True(t, actual.Equal(utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0)))
	})
	t.Run("2020_02_02を渡すと直近の課金開始日_2020_02_02が返る", func(t *testing.T) {
		actual := contract.LastBillingStartDate(utils.CreateJstTime(2020, 2, 2, 0, 0, 0, 0))
		assert.True(t, actual.Equal(utils.CreateJstTime(2020, 2, 2, 0, 0, 0, 0)))
	})
	t.Run("2020_03_01を渡すと直近の課金開始日_2020_02_02が返る", func(t *testing.T) {
		actual := contract.LastBillingStartDate(utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0))
		assert.True(t, actual.Equal(utils.CreateJstTime(2020, 2, 2, 0, 0, 0, 0)))
	})
}
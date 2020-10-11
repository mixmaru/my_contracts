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
	rightToUses := make([]*RightToUseEntity, 0, 2)
	rightToUses = append(rightToUses, NewRightToUseEntity(
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
	))
	rightToUses = append(rightToUses, NewRightToUseEntity(
		time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
	))
	entity := NewContractEntity(
		1,
		2,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
		rightToUses,
	)

	// テスト
	assert.Equal(t, 1, entity.UserId())
	assert.Equal(t, 2, entity.ProductId())
	assert.EqualValues(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), entity.ContractDate())
	assert.EqualValues(t, time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), entity.BillingStartDate())
	// 使用権
	actualRightToUses := entity.RightToUses()
	assert.Len(t, actualRightToUses, 2)
	assert.True(t, actualRightToUses[0].validFrom.Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
	assert.True(t, actualRightToUses[0].validTo.Equal(time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)))
	assert.True(t, actualRightToUses[1].validFrom.Equal(time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)))
	assert.True(t, actualRightToUses[1].validTo.Equal(time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)))
	actualRightToUses[0].validFrom = time.Time{}
}

func TestContractEntity_NewContractEntityWithData(t *testing.T) {
	// インスタンス化
	contractDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	billingStartDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	entity, err := NewContractEntityWithData(1, 2, 3, contractDate, billingStartDate, createdAt, updatedAt, []*RightToUseEntity{})
	assert.NoError(t, err)

	assert.Equal(t, 1, entity.Id())
	assert.Equal(t, 2, entity.UserId())
	assert.Equal(t, 3, entity.ProductId())
	assert.True(t, contractDate.Equal(entity.ContractDate()))
	assert.True(t, billingStartDate.Equal(entity.BillingStartDate()))
	assert.Len(t, entity.RightToUses(), 0)
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
			[]*RightToUseEntity{},
		)
		assert.NoError(t, err)

		assert.Equal(t, 1, contractEntity.Id())
		assert.Equal(t, 2, contractEntity.UserId())
		assert.Equal(t, 3, contractEntity.ProductId())
		assert.EqualValues(t, contractDate, contractEntity.ContractDate())
		assert.EqualValues(t, billingStartDate, contractEntity.BillingStartDate())
		assert.Len(t, contractEntity.RightToUses(), 0)
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
		[]*RightToUseEntity{},
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

func TestContractEntity_ArchiveRightToUseById(t *testing.T) {
	////// 準備
	entity, err := NewContractEntityWithData(
		1, 2, 3, utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		[]*RightToUseEntity{
			NewRightToUseEntityWithData(
				4,
				utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
				6,
				utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			),
			NewRightToUseEntityWithData(
				5,
				utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0),
				0,
				utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
			),
		},
	)
	assert.NoError(t, err)

	t.Run("その契約集約が持っている使用権Idを渡すと、その使用権がアーカイブ（リポジトリで更新をかけたときにhistoryテーブル行き）される", func(t *testing.T) {
		////// 実行
		err = entity.ArchiveRightToUseById(4)
		assert.NoError(t, err)

		////// 検証
		rightToUses := entity.RightToUses()
		assert.Len(t, rightToUses, 1)
		assert.Equal(t, 5, rightToUses[0].Id())
		assert.Equal(t, []int{4}, entity.GetToArchiveRightToUseIds())
	})

	t.Run("その契約集約が持っていない使用権Idを渡すとエラーが発生する", func(t *testing.T) {
		////// 実行
		err = entity.ArchiveRightToUseById(100)
		////// 検証
		assert.Error(t, err)
	})
}

func TestContractEntity_ArchiveRightToUseByValidTo(t *testing.T) {
	t.Run("使用期限日（ValidTo）を渡すと、それ以前に使用期限日時がある使用権がアーカイブ（リポジトリで更新をかけたときにhistoryテーブル行き）される", func(t *testing.T) {
		////// 準備
		entity, err := NewContractEntityWithData(
			1, 2, 3, utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			[]*RightToUseEntity{
				NewRightToUseEntityWithData(
					4,
					utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
					6,
					utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
				),
				NewRightToUseEntityWithData(
					5,
					utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0),
					0,
					utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
				),
			},
		)
		assert.NoError(t, err)

		////// 実行
		err = entity.ArchiveRightToUseByValidTo(utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0))
		assert.NoError(t, err)

		////// 検証
		rightToUses := entity.RightToUses()
		assert.Len(t, rightToUses, 0)
		assert.Equal(t, []int{4, 5}, entity.GetToArchiveRightToUseIds())
	})
}

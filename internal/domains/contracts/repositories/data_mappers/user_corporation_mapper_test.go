package data_mappers

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserCorporationMapper_NewUserCorporationMapperFromUserCorporationEntity(t *testing.T) {
	// entity用意
	entity, err := entities.NewUserCorporationEntityWithData(
		1,
		"担当太郎",
		"社長次郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	// 実行
	resultData := NewUserCorporationMapperFromUserCorporationEntity(entity)

	expect := &UserCorporationMapper{
		UserId:            1,
		ContactParsonName: "担当太郎",
		PresidentName:     "社長次郎",
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	// テスト
	assert.Equal(t, expect, resultData)
}

func TestUserCorporationMapper_PreInsert(t *testing.T) {
	// entity用意
	newUser := UserCorporationMapper{
		UserId:            0,
		ContactParsonName: "担当太郎",
		PresidentName:     "社長次郎",
	}

	err := newUser.PreInsert(nil)
	assert.NoError(t, err)

	assert.NotEqual(t, time.Time{}, newUser.CreatedAt)
	assert.NotEqual(t, time.Time{}, newUser.UpdatedAt)
}

func TestUserCorporationMapper_PreUpdate(t *testing.T) {
	// entity用意
	newUser := UserCorporationMapper{
		UserId:            0,
		ContactParsonName: "担当太郎",
		PresidentName:     "社長次郎",
	}

	err := newUser.PreUpdate(nil)
	assert.NoError(t, err)

	assert.Equal(t, time.Time{}, newUser.CreatedAt)
	assert.NotEqual(t, time.Time{}, newUser.UpdatedAt)
}

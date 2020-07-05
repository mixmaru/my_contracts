package data_mappers

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserIndividualMapper_NewUserIndividualMapperFromUserIndividualEntity(t *testing.T) {
	// entity用意
	entity, err := entities.NewUserIndividualEntityWithData(
		1,
		"担当太郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	// 実行
	resultData := NewUserIndividualMapperFromUserIndividualEntity(entity)

	expect := &UserIndividualMapper{
		UserId: 1,
		Name:   "担当太郎",
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	// テスト
	assert.Equal(t, expect, resultData)
}

func TestUserIndividualMapper_PreInsert(t *testing.T) {
	// entity用意
	newUser := UserIndividualMapper{
		UserId: 0,
		Name:   "担当太郎",
	}

	err := newUser.PreInsert(nil)
	assert.NoError(t, err)

	assert.NotEqual(t, time.Time{}, newUser.CreatedAt)
	assert.NotEqual(t, time.Time{}, newUser.UpdatedAt)
}

func TestUserIndividualMapper_PreUpdate(t *testing.T) {
	// entity用意
	newUser := UserIndividualMapper{
		UserId: 0,
		Name:   "担当太郎",
	}

	err := newUser.PreUpdate(nil)
	assert.NoError(t, err)

	assert.Equal(t, time.Time{}, newUser.CreatedAt)
	assert.NotEqual(t, time.Time{}, newUser.UpdatedAt)
}

package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUser_NewUserRecordFromUserIndividualEntity(t *testing.T) {
	userEntity, err := entities.NewUserIndividualEntityWithData(
		1,
		"個人たろう",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	user := NewUserRecordFromUserIndividualEntity(userEntity)
	expect := &UserRecord{
		Id: 1,
		CreatedAtUpdatedAt: CreatedAtUpdatedAt{
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	assert.Equal(t, expect, user)
}

func TestUser_NewUserRecordFromUserCorporationEntity(t *testing.T) {
	userEntity, err := entities.NewUserCorporationEntityWithData(
		1,
		"担当たろう",
		"社長じろう",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	user := NewUserRecordFromUserCorporationEntity(userEntity)
	expect := &UserRecord{
		Id: 1,
		CreatedAtUpdatedAt: CreatedAtUpdatedAt{
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	assert.Equal(t, expect, user)
}

func TestUserRecord_PreInsert(t *testing.T) {
	// entity用意
	newUser := UserRecord{
		Id: 0,
	}

	err := newUser.PreInsert(nil)
	assert.NoError(t, err)

	assert.NotEqual(t, time.Time{}, newUser.CreatedAt)
	assert.NotEqual(t, time.Time{}, newUser.UpdatedAt)
}

func TestUserRecord_PreUpdate(t *testing.T) {
	// entity用意
	newUser := UserRecord{
		Id: 0,
	}

	err := newUser.PreUpdate(nil)
	assert.NoError(t, err)

	assert.Equal(t, time.Time{}, newUser.CreatedAt)
	assert.NotEqual(t, time.Time{}, newUser.UpdatedAt)
}

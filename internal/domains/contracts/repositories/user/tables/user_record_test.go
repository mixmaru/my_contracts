package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUser_NewUserRecordFromUserIndividualEntity(t *testing.T) {
	userEntity := user.NewUserIndividualEntityWithData(
		1,
		"個人たろう",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)

	user := NewUserRecordFromUserIndividualEntity(userEntity)
	expect := &UserRecord{
		Id:        1,
		CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	assert.Equal(t, expect, user)
}

func TestUser_NewUserRecordFromUserCorporationEntity(t *testing.T) {
	userEntity := user.NewUserCorporationEntityWithData(
		1,
		"担当たろう",
		"社長じろう",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)

	user := NewUserRecordFromUserCorporationEntity(userEntity)
	expect := &UserRecord{
		Id:        1,
		CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	assert.Equal(t, expect, user)
}

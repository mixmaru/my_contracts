package db_maps

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUser_LoadUserIndividual(t *testing.T) {
	userEntity := user.NewUserIndividualEntityWithData(
		1,
		"個人たろう",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)

	user := NewUserFromUserIndividualEntity(userEntity)
	expect := &UserDbMap{
		Id:        1,
		CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	assert.Equal(t, expect, user)
}

package data_mappers

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUser_NewUserMapperFromUserIndividualEntity(t *testing.T) {
	userEntity, err := entities.NewUserIndividualEntityWithData(
		1,
		"個人たろう",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	user := NewUserMapperFromUserIndividualEntity(userEntity)
	expect := &UserMapper{
		Id: 1,
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	assert.Equal(t, expect, user)
}

func TestUser_NewUserMapperFromUserCorporationEntity(t *testing.T) {
	userEntity, err := entities.NewUserCorporationEntityWithData(
		1,
		"イケてる会社",
		"担当たろう",
		"社長じろう",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	user := NewUserMapperFromUserCorporationEntity(userEntity)
	expect := &UserMapper{
		Id: 1,
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	assert.Equal(t, expect, user)
}

func TestUserMapper_PreInsert(t *testing.T) {
	// entity用意
	newUser := UserMapper{
		Id: 0,
	}

	err := newUser.PreInsert(nil)
	assert.NoError(t, err)

	assert.NotEqual(t, time.Time{}, newUser.CreatedAt)
	assert.NotEqual(t, time.Time{}, newUser.UpdatedAt)
}

func TestUserMapper_PreUpdate(t *testing.T) {
	// entity用意
	newUser := UserMapper{
		Id: 0,
	}

	err := newUser.PreUpdate(nil)
	assert.NoError(t, err)

	assert.Equal(t, time.Time{}, newUser.CreatedAt)
	assert.NotEqual(t, time.Time{}, newUser.UpdatedAt)
}

func TestUserMapper_SetDataToEntity(t *testing.T) {
	userMapper := UserMapper{}
	userMapper.Id = 1
	userMapper.CreatedAt = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	userMapper.UpdatedAt = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	var userEntity entities.UserEntity
	err := userMapper.SetDataToEntity(&userEntity)
	assert.NoError(t, err)

	assert.Equal(t, 1, userEntity.Id())
	assert.True(t, userEntity.CreatedAt().Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
	assert.True(t, userEntity.UpdatedAt().Equal(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
}

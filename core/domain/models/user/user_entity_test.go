package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// インスタンス化された顧客Entityに対してデータロードするやつ
func TestUserEntity_LoadData(t *testing.T) {
	user := &UserEntity{}
	user.LoadData(
		1,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)

	assert.Equal(t, 1, user.Id())
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), user.CreatedAt())
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), user.UpdatedAt())
}

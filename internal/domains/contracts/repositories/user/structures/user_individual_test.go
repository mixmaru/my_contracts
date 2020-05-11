package structures

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserIndividual_LoadUserIndividual(t *testing.T) {
	// entity用意
	entity := user.NewUserIndividualWithData(
		1,
		"担当太郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	)

	// 実行
	resultData := NewUserIndividualFromUserIndividual(entity)

	expect := &UserIndividual{
		UserId:    1,
		Name:      "担当太郎",
		CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// テスト
	assert.Equal(t, expect, resultData)
}

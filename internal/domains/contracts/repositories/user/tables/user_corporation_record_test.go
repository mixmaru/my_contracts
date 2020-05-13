package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserCorporationRecord_NewUserCorporationRecordFromUserCorporationEntity(t *testing.T) {
	// entity用意
	entity := user.NewUserCorporationEntityWithData(
		1,
		"担当太郎",
		"社長次郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)

	// 実行
	resultData := NewUserCorporationRecordFromUserCorporationEntity(entity)

	expect := &UserCorporationRecord{
		UserId:            1,
		ContactParsonName: "担当太郎",
		PresidentName:     "社長次郎",
		CreatedAt:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// テスト
	assert.Equal(t, expect, resultData)
}

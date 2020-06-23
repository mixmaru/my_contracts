package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserCorporationView_SetDataToEntity(t *testing.T) {
	record := UserCorporationView{}
	record.Id = 1
	record.ContactPersonName = "担当太郎"
	record.PresidentName = "社長太郎"
	record.CreatedAt = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	record.UpdatedAt = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	entity := entities.NewUserCorporationEntity()
	err := record.SetDataToEntity(entity)
	assert.NoError(t, err)

	assert.Equal(t, 1, entity.Id())
	assert.Equal(t, "担当太郎", entity.ContactPersonName())
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), entity.CreatedAt())
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), entity.UpdatedAt())
}

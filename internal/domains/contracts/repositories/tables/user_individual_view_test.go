package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserIndividualView_SetDataToEntity(t *testing.T) {
	record := UserIndividualView{}
	record.Id = 1
	record.Name = "個人たろう"
	record.CreatedAt = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	record.UpdatedAt = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	entity := entities.UserIndividualEntity{}
	err := record.SetDataToEntity(&entity)
	assert.NoError(t, err)

	assert.Equal(t, 1, entity.Id())
	assert.Equal(t, "個人たろう", entity.Name())
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), entity.CreatedAt())
	assert.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), entity.UpdatedAt())
}

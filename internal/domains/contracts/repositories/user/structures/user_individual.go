package structures

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"time"
)

type UserIndividual struct {
	UserId    int       `db:"user_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func LoadUserIndividual(entity *user.UserIndividual) *UserIndividual {
	return &UserIndividual{
		UserId:    entity.Id(),
		Name:      entity.Name(),
		CreatedAt: entity.CreatedAt(),
		UpdatedAt: entity.UpdatedAt(),
	}
}

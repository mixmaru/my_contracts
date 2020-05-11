package db_maps

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"time"
)

type UserDbMap struct {
	Id        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// UserIndividualEntity Entityからデータを読み込んでUser(DBマッピング用)を作成する
func NewUserFromUserIndividualEntity(userIndividual *user.UserIndividualEntity) *UserDbMap {
	return &UserDbMap{
		Id:        userIndividual.Id(),
		CreatedAt: userIndividual.CreatedAt(),
		UpdatedAt: userIndividual.UpdatedAt(),
	}
}

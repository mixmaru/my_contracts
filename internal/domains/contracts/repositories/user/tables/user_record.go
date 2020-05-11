package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"time"
)

type UserRecord struct {
	Id        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// UserIndividualEntity Entityからデータを読み込んでUser(DBマッピング用)を作成する
func NewUserRecordFromUserIndividualEntity(userIndividual *user.UserIndividualEntity) *UserRecord {
	return &UserRecord{
		Id:        userIndividual.Id(),
		CreatedAt: userIndividual.CreatedAt(),
		UpdatedAt: userIndividual.UpdatedAt(),
	}
}

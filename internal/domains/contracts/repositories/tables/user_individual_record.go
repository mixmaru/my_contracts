package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
)

type UserIndividualRecord struct {
	UserId int    `db:"user_id"`
	Name   string `db:"name"`
	CreatedAtUpdatedAt
}

// UserIndividualEntity Entityからデータを読み込んでUserIndividual(DBマッピング用)を作成する
func NewUserIndividualRecordFromUserIndividualEntity(entity *user.UserIndividualEntity) *UserIndividualRecord {
	return &UserIndividualRecord{
		UserId: entity.Id(),
		Name:   entity.Name(),
		CreatedAtUpdatedAt: CreatedAtUpdatedAt{
			CreatedAt: entity.CreatedAt(),
			UpdatedAt: entity.UpdatedAt(),
		},
	}
}

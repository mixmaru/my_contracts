package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"gopkg.in/gorp.v2"
	"time"
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

// insert時に時刻をセットするhook
func (u *UserIndividualRecord) PreInsert(s gorp.SqlExecutor) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	return nil
}

// updateに時刻をセットするhook
func (u *UserIndividualRecord) PreUpdate(s gorp.SqlExecutor) error {
	u.UpdatedAt = time.Now()
	return nil
}

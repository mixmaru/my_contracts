package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"gopkg.in/gorp.v2"
	"time"
)

type UserCorporationRecord struct {
	UserId            int       `db:"user_id"`
	ContactParsonName string    `db:"contact_person_name"`
	PresidentName     string    `db:"president_name"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// UserIndividualEntity Entityからデータを読み込んでUserIndividual(DBマッピング用)を作成する
func NewUserCorporationRecordFromUserCorporationEntity(entity *user.UserCorporationEntity) *UserCorporationRecord {
	return &UserCorporationRecord{
		UserId:            entity.Id(),
		ContactParsonName: entity.ContactPersonName(),
		PresidentName:     entity.PresidentName(),
		CreatedAt:         entity.CreatedAt(),
		UpdatedAt:         entity.UpdatedAt(),
	}
}

// insert時に時刻をセットするhook
func (u *UserCorporationRecord) PreInsert(s gorp.SqlExecutor) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	return nil
}

// updateに時刻をセットするhook
func (u *UserCorporationRecord) PreUpdate(s gorp.SqlExecutor) error {
	u.UpdatedAt = time.Now()
	return nil
}

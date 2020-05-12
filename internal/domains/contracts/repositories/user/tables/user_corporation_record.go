package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"time"
)

type UserCorporationRecord struct {
	UserId            int       `db:"user_id"`
	ContactParsonName string    `db:"contact_parson_name"`
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

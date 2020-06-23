package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
)

type UserCorporationRecord struct {
	UserId            int    `db:"user_id"`
	ContactParsonName string `db:"contact_person_name"`
	PresidentName     string `db:"president_name"`
	CreatedAtUpdatedAt
}

// UserIndividualEntity Entityからデータを読み込んでUserIndividual(DBマッピング用)を作成する
func NewUserCorporationRecordFromUserCorporationEntity(entity *entities.UserCorporationEntity) *UserCorporationRecord {
	return &UserCorporationRecord{
		UserId:            entity.Id(),
		ContactParsonName: entity.ContactPersonName(),
		PresidentName:     entity.PresidentName(),
		CreatedAtUpdatedAt: CreatedAtUpdatedAt{
			CreatedAt: entity.CreatedAt(),
			UpdatedAt: entity.UpdatedAt(),
		},
	}
}

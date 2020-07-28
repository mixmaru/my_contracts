package data_mappers

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
)

type UserCorporationMapper struct {
	UserId            int    `db:"user_id"`
	ContactParsonName string `db:"contact_person_name"`
	PresidentName     string `db:"president_name"`
	CompanyName       string `db:"company_name"`
	CreatedAtUpdatedAtMapper
}

// UserIndividualEntity Entityからデータを読み込んでUserIndividual(DBマッピング用)を作成する
func NewUserCorporationMapperFromUserCorporationEntity(entity *entities.UserCorporationEntity) *UserCorporationMapper {
	return &UserCorporationMapper{
		UserId:            entity.Id(),
		CompanyName:       entity.CorporationName(),
		ContactParsonName: entity.ContactPersonName(),
		PresidentName:     entity.PresidentName(),
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: entity.CreatedAt(),
			UpdatedAt: entity.UpdatedAt(),
		},
	}
}

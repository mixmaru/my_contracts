package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
)

type UserRecord struct {
	Id int `db:"id"`
	CreatedAtUpdatedAt
}

// UserIndividualEntityからデータを読み込んでUser(DBマッピング用)を作成する
func NewUserRecordFromUserIndividualEntity(userIndividual *user.UserIndividualEntity) *UserRecord {
	return &UserRecord{
		Id: userIndividual.Id(),
		CreatedAtUpdatedAt: CreatedAtUpdatedAt{
			CreatedAt: userIndividual.CreatedAt(),
			UpdatedAt: userIndividual.UpdatedAt(),
		},
	}
}

// UserCorporationEntityからデータを読み込んでUser(DBマッピング用)を作成する
func NewUserRecordFromUserCorporationEntity(userCorporation *user.UserCorporationEntity) *UserRecord {
	return &UserRecord{
		Id: userCorporation.Id(),
		CreatedAtUpdatedAt: CreatedAtUpdatedAt{
			CreatedAt: userCorporation.CreatedAt(),
			UpdatedAt: userCorporation.UpdatedAt(),
		},
	}
}

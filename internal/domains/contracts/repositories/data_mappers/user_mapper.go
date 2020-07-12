package data_mappers

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
)

type UserMapper struct {
	Id int `db:"id"`
	CreatedAtUpdatedAtMapper
}

// UserIndividualEntityからデータを読み込んでUser(DBマッピング用)を作成する
func NewUserMapperFromUserIndividualEntity(userIndividual *entities.UserIndividualEntity) *UserMapper {
	return &UserMapper{
		Id: userIndividual.Id(),
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: userIndividual.CreatedAt(),
			UpdatedAt: userIndividual.UpdatedAt(),
		},
	}
}

// UserCorporationEntityからデータを読み込んでUser(DBマッピング用)を作成する
func NewUserMapperFromUserCorporationEntity(userCorporation *entities.UserCorporationEntity) *UserMapper {
	return &UserMapper{
		Id: userCorporation.Id(),
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: userCorporation.CreatedAt(),
			UpdatedAt: userCorporation.UpdatedAt(),
		},
	}
}
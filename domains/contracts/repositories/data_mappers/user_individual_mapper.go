package data_mappers

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
)

type UserIndividualMapper struct {
	UserId int    `db:"user_id"`
	Name   string `db:"name"`
	CreatedAtUpdatedAtMapper
}

// UserIndividualEntity Entityからデータを読み込んでUserIndividual(DBマッピング用)を作成する
func NewUserIndividualMapperFromUserIndividualEntity(entity *entities.UserIndividualEntity) *UserIndividualMapper {
	return &UserIndividualMapper{
		UserId: entity.Id(),
		Name:   entity.Name(),
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: entity.CreatedAt(),
			UpdatedAt: entity.UpdatedAt(),
		},
	}
}

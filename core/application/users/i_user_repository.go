package users

import (
	entity "github.com/mixmaru/my_contracts/core/domain/models/user"
	"gopkg.in/gorp.v2"
)

type IUserRepository interface {
	GetUserById(id int, executor gorp.SqlExecutor) (interface{}, error)

	SaveUserIndividual(userEntity *entity.UserIndividualEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetUserIndividualById(id int, executor gorp.SqlExecutor) (*entity.UserIndividualEntity, error)

	SaveUserCorporation(userEntity *entity.UserCorporationEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetUserCorporationById(id int, executor gorp.SqlExecutor) (*entity.UserCorporationEntity, error)
}

package interfaces

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"gopkg.in/gorp.v2"
)

type IUserRepository interface {
	GetUserById(id int, executor gorp.SqlExecutor) (interface{}, error)

	SaveUserIndividual(userEntity *entities.UserIndividualEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetUserIndividualById(id int, executor gorp.SqlExecutor) (*entities.UserIndividualEntity, error)

	SaveUserCorporation(userEntity *entities.UserCorporationEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetUserCorporationById(id int, executor gorp.SqlExecutor) (*entities.UserCorporationEntity, error)
}

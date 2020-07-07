package interfaces

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"gopkg.in/gorp.v2"
)

type IUserRepository interface {
	SaveUserIndividual(userEntity *entities.UserIndividualEntity, executor gorp.SqlExecutor) (*entities.UserIndividualEntity, error)
	GetUserIndividualById(id int, executor gorp.SqlExecutor) (*entities.UserIndividualEntity, error)

	SaveUserCorporation(userEntity *entities.UserCorporationEntity, executor gorp.SqlExecutor) (*entities.UserCorporationEntity, error)
	GetUserCorporationById(id int, executor gorp.SqlExecutor) (*entities.UserCorporationEntity, error)
}

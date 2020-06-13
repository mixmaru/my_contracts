package interfaces

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"gopkg.in/gorp.v2"
)

type IUserRepository interface {
	SaveUserIndividual(userEntity *entities.UserIndividualEntity, transaction *gorp.Transaction) (*entities.UserIndividualEntity, error)
	GetUserIndividualById(id int, transaction *gorp.Transaction) (*entities.UserIndividualEntity, error)

	SaveUserCorporation(userEntity *entities.UserCorporationEntity, transaction *gorp.Transaction) (*entities.UserCorporationEntity, error)
	GetUserCorporationById(id int, transaction *gorp.Transaction) (*entities.UserCorporationEntity, error)
}

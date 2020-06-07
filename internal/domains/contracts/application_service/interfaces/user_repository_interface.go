package interfaces

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"gopkg.in/gorp.v2"
)

type IUserRepository interface {
	SaveUserIndividual(userEntity *user.UserIndividualEntity, transaction *gorp.Transaction) (*user.UserIndividualEntity, error)
	GetUserIndividualById(id int, transaction *gorp.Transaction) (*user.UserIndividualEntity, error)

	SaveUserCorporation(userEntity *user.UserCorporationEntity, transaction *gorp.Transaction) (*user.UserCorporationEntity, error)
	GetUserCorporationById(id int, transaction *gorp.Transaction) (*user.UserCorporationEntity, error)
}

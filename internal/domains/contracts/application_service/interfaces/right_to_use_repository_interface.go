package interfaces

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"gopkg.in/gorp.v2"
)

type IRightToUseRepository interface {
	Create(rightToUseEntity *entities.RightToUseEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetById(id int, executor gorp.SqlExecutor) (*entities.RightToUseEntity, error)
}

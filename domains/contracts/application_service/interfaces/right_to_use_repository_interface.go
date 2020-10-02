package interfaces

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"gopkg.in/gorp.v2"
	"time"
)

type IRightToUseRepository interface {
	GetRecurTargets(executeDate time.Time, executor gorp.SqlExecutor) ([]*entities.RightToUseEntity, error)
}

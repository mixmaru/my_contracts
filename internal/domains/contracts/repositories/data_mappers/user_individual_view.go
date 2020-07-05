package data_mappers

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/pkg/errors"
)

type UserIndividualView struct {
	UserRecord
	Name string
}

func (u *UserIndividualView) SetDataToEntity(entity interface{}) error {
	value, ok := entity.(*entities.UserIndividualEntity)
	if !ok {
		return errors.New(fmt.Sprintf("*entities.UserIndividualEntity型ではないものが渡ってきた。 %v", entity))
	}

	err := value.LoadData(u.Id, u.Name, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

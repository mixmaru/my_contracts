package tables

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/pkg/errors"
	"reflect"
)

type UserCorporationView struct {
	UserRecord
	ContactPersonName string `db:"contact_person_name"`
	PresidentName     string `db:"president_name"`
}

func (u *UserCorporationView) SetDataToEntity(entity interface{}) error {
	userCorprationEntity, ok := entity.(*entities.UserCorporationEntity)
	if !ok {
		return errors.New(fmt.Sprintf("entityが*entities.UserCorporationEntityではない。%v", reflect.TypeOf(entity)))
	}
	err := userCorprationEntity.LoadData(u.Id, u.ContactPersonName, u.PresidentName, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

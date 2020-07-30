package data_mappers

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/pkg/errors"
	"reflect"
)

type UserCorporationView struct {
	UserMapper
	CorporationName   string `db:"company_name"`
	ContactPersonName string `db:"contact_person_name"`
	PresidentName     string `db:"president_name"`
}

func (u *UserCorporationView) SetDataToEntity(entity interface{}) error {
	userCorporationEntity, ok := entity.(*entities.UserCorporationEntity)
	if !ok {
		return errors.New(fmt.Sprintf("entityが*entities.UserCorporationEntityではない。%v", reflect.TypeOf(entity)))
	}
	err := userCorporationEntity.LoadData(u.Id, u.CorporationName, u.ContactPersonName, u.PresidentName, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

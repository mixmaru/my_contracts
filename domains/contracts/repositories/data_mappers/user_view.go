package data_mappers

import (
	"database/sql"
)

type UserView struct {
	UserId                  int            `db:"user_id"`
	UserType                string         `db:"user_type"`
	UserIndividualName      sql.NullString `db:"user_individual_name"`
	UserIndividualCreatedAt sql.NullTime   `db:"user_individual_created_at"`
	UserIndividualUpdatedAt sql.NullTime   `db:"user_individual_updated_at"`

	UserCorporationCorporationName    sql.NullString `db:"user_corporation_corporation_name"`
	UserCorporationContractPersonName sql.NullString `db:"user_corporation_contact_person_name"`
	UserCorporationPresidentName      sql.NullString `db:"user_corporation_president_name"`
	UserCorporationCreatedAt          sql.NullTime   `db:"user_corporation_created_at"`
	UserCorporationUpdatedAt          sql.NullTime   `db:"user_corporation_updated_at"`
}

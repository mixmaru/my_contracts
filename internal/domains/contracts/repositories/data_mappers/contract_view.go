package data_mappers

import (
	"database/sql"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"time"
)

type ContractView struct {
	Id        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	ProductId        int             `db:"product_id"`
	ProductName      string          `db:"product_name"`
	ProductPrice     decimal.Decimal `db:"product_price"`
	ProductCreatedAt time.Time       `db:"product_created_at"`
	ProductUpdatedAt time.Time       `db:"product_updated_at"`

	UserId                  int            `db:"user_id"`
	UserType                string         `db:"user_type"`
	UserIndividualName      sql.NullString `db:"user_individual_name"`
	UserIndividualCreatedAt sql.NullTime   `db:"user_individual_created_at"`
	UserIndividualUpdatedAt sql.NullTime   `db:"user_individual_updated_at"`

	UserCorporationContractPersonName sql.NullString `db:"user_corporation_contact_person_name"`
	UserCorporationPresidentName      sql.NullString `db:"user_corporation_president_name"`
	UserCorporationCreatedAt          sql.NullTime   `db:"user_corporation_created_at"`
	UserCorporationUpdatedAt          sql.NullTime   `db:"user_corporation_updated_at"`
}

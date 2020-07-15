package data_mappers

import (
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"time"
)

type ContractView struct {
	Id        int       `db:"id"`
	UserId    int       `db:"user_id"`
	UserType  string    `db:"user_type"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	UserIndividualName      string    `db:"user_individual_name"`
	UserIndividualCreatedAt time.Time `db:"user_individual_created_at"`
	UserIndividualUpdatedAt time.Time `db:"user_individual_updated_at"`

	UserCorporationContractPersonName string    `db:"user_corporation_contract_person_name"`
	UserCorporationPresidentName      string    `db:"user_corporation_president_name"`
	UserCorporationCreatedAt          time.Time `db:"user_corporation_created_at"`
	UserCorporationUpdatedAt          time.Time `db:"user_corporation_updated_at"`

	ProductId        int             `db:"product_id"`
	ProductName      string          `db:"product_name"`
	ProductPrice     decimal.Decimal `db:"product_price"`
	ProductCreatedAt time.Time       `db:"product_created_at"`
	ProductUpdatedAt time.Time       `db:"product_updated_at"`
}

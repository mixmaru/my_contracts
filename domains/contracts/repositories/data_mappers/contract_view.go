package data_mappers

import (
	"github.com/mixmaru/my_contracts/lib/decimal"
	"time"
)

type ContractView struct {
	Id               int       `db:"id"`
	ContractDate     time.Time `db:"contract_date""`
	BillingStartDate time.Time `db:"billing_start_date"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`

	ProductId        int             `db:"product_id"`
	ProductName      string          `db:"product_name"`
	ProductPrice     decimal.Decimal `db:"product_price"`
	ProductCreatedAt time.Time       `db:"product_created_at"`
	ProductUpdatedAt time.Time       `db:"product_updated_at"`

	UserView
}

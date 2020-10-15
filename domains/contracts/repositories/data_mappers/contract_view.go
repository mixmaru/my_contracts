package data_mappers

import (
	"database/sql"
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

	RightToUseId        sql.NullInt64 `db:"right_to_use_id"`
	RightToUseValidFrom sql.NullTime  `db:"right_to_use_valid_from"`
	RightToUseValidTo   sql.NullTime  `db:"right_to_use_valid_to"`
	RightToUseCreatedAt sql.NullTime  `db:"right_to_use_active_created_at"`
	RightToUseUpdatedAt sql.NullTime  `db:"right_to_use_active_updated_at"`

	BillDetailId int `db:"bill_detail_id"`

	UserView
}

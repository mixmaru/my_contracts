package data_mappers

import (
	"database/sql"
	"time"
)

type BillMapper struct {
	Id                 int          `db:"id"`
	BillingDate        time.Time    `db:"billing_date"`
	UserId             int          `db:"user_id"`
	PaymentConfirmedAt sql.NullTime `db:"payment_confirmed_at"`
	CreatedAtUpdatedAtMapper
}

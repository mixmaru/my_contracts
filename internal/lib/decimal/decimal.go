package decimal

import (
	"database/sql/driver"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Decimal struct {
	decimal decimal.Decimal
}

func New(value int64, exp int32) Decimal {
	return Decimal{decimal.New(value, exp)}
}

func NewFromInt(value int64) Decimal {
	return Decimal{decimal.NewFromInt(value)}
}

func NewFromFloat(value float64) Decimal {
	return Decimal{decimal.NewFromFloat(value)}
}

func NewFromString(value string) (Decimal, error) {
	decimal, err := decimal.NewFromString(value)
	if err != nil {
		return Decimal{}, errors.WithStack(err)
	}

	return Decimal{decimal}, nil
}

func (d *Decimal) String() string {
	return d.decimal.String()
}

func (d *Decimal) Add(decimal Decimal) Decimal {
	return Decimal{d.decimal.Add(decimal.decimal)}
}

func (d *Decimal) Sub(decimal Decimal) Decimal {
	return Decimal{d.decimal.Sub(decimal.decimal)}
}

func (d *Decimal) Mul(decimal Decimal) Decimal {
	return Decimal{d.decimal.Mul(decimal.decimal)}
}

func (d *Decimal) Div(decimal Decimal) Decimal {
	return Decimal{d.decimal.Div(decimal.decimal)}
}

func (d *Decimal) Equal(decimal Decimal) bool {
	return d.decimal.Equal(decimal.decimal)
}

func (d *Decimal) IsPositive() bool {
	return d.decimal.IsPositive()
}

func (d *Decimal) IsNegative() bool {
	return d.decimal.IsNegative()
}

func (d *Decimal) IsZero() bool {
	return d.decimal.IsZero()
}

// Value implements the driver.Valuer interface for database serialization.
func (d Decimal) Value() (driver.Value, error) {
	return d.decimal.Value()
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Decimal) Scan(value interface{}) error {
	return d.decimal.Scan(value)
}

package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/pkg/errors"
)

type ProductRecord struct {
	Id    int             `db:"id"`
	Name  string          `db:"name"`
	Price decimal.Decimal `db:"price"`
	CreatedAtUpdatedAt
}

func (p *ProductRecord) SetDataToEntity(entity interface{}) error {
	value, ok := entity.(*entities.ProductEntity)
	if !ok {
		return errors.New("*entities.ProductEntityではないものが渡された")
	}
	err := value.LoadData(
		p.Id,
		p.Name,
		p.Price.String(),
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

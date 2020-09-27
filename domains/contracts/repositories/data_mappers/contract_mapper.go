package data_mappers

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/pkg/errors"
	"time"
)

type ContractMapper struct {
	Id               int       `db:"id"`
	UserId           int       `db:"user_id"`
	ProductId        int       `db:"product_id"`
	ContractDate     time.Time `db:"contract_date"`
	BillingStartDate time.Time `db:"billing_start_date"`
	CreatedAtUpdatedAtMapper
}

func (c *ContractMapper) SetDataToEntity(entity interface{}) error {
	value, ok := entity.(*entities.ContractEntity)
	if !ok {
		return errors.Errorf("*entities.ContractEntityではないものが渡された。entity: %t", entity)
	}
	err := value.LoadData(
		c.Id,
		c.UserId,
		c.ProductId,
		c.ContractDate,
		c.BillingStartDate,
		c.CreatedAt,
		c.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

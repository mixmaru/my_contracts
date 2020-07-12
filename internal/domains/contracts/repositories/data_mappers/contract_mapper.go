package data_mappers

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/pkg/errors"
)

type ContractMapper struct {
	Id        int `db:"id"`
	UserId    int `db:"user_id"`
	ProductId int `db:"product_id"`
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
		c.CreatedAt,
		c.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

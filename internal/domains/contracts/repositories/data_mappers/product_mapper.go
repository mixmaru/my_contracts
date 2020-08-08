package data_mappers

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/pkg/errors"
)

type ProductMapper struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	CreatedAtUpdatedAtMapper
}

func (p *ProductMapper) SetDataToEntity(entity interface{}) error {
	value, ok := entity.(*entities.ProductEntity)
	if !ok {
		return errors.New("*entities.ProductEntityではないものが渡された")
	}
	err := value.LoadData(
		p.Id,
		p.Name,
		"0", // todo: あとで調整の必要
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

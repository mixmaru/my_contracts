package data_mappers

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"time"
)

type RightToUseMapper struct {
	Id         int       `db:"id"`
	ContractId int       `db:"contract_id"`
	ValidFrom  time.Time `db:"valid_from"`
	ValidTo    time.Time `db:"valid_to"`
	CreatedAtUpdatedAtMapper
}

func NewRightToUseMapperFromEntity(entity *entities.RightToUseEntity) RightToUseMapper {
	mapper := RightToUseMapper{}
	mapper.Id = entity.Id()
	mapper.ValidFrom = entity.ValidFrom()
	mapper.ValidTo = entity.ValidTo()
	mapper.CreatedAt = entity.CreatedAt()
	mapper.UpdatedAt = entity.UpdatedAt()
	return mapper
}

type RightToUseActiveMapper struct {
	RightToUseId int `db:"right_to_use_id"`
	CreatedAtUpdatedAtMapper
}

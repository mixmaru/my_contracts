package data_mappers

import (
	"time"
)

type RightToUseMapper struct {
	Id         int       `db:"id"`
	ContractId int       `db:"contract_id"`
	ValidFrom  time.Time `db:"valid_from"`
	ValidTo    time.Time `db:"valid_to"`
	CreatedAtUpdatedAtMapper
}

type RightToUseActiveMapper struct {
	RightToUseId int `db:"right_to_use_id"`
	CreatedAtUpdatedAtMapper
}

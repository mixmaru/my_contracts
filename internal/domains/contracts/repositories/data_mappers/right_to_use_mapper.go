package data_mappers

import "time"

type RightToUseMapper struct {
	Id         int       `db:"id"`
	ContractId string    `db:"contract_id"`
	ValidFrom  time.Time `db:"valid_from"`
	ValidTo    time.Time `db:"valid_to"`
	CreatedAtUpdatedAtMapper
}

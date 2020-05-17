package data_transfer_objects

import "time"

type UserIndividualDto struct {
	Id        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

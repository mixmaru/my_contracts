package data_transfer_objects

import "time"

type BaseDto struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

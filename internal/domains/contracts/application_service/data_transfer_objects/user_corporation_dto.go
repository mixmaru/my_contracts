package data_transfer_objects

import "time"

type UserCorporationDto struct {
	Id                int
	ContactPersonName string
	PresidentName     string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

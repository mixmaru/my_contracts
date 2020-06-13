package entities

import "time"

type ICreatedAtUpdatedAt interface {
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type CreatedAtUpdatedAt struct {
	createdAt time.Time
	updatedAt time.Time
}

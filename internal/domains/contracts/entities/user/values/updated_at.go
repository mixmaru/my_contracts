package values

import "time"

type UpdatedAt struct {
	value time.Time
}

func NewUpdatedAt(time time.Time) UpdatedAt {
	return UpdatedAt{
		value: time,
	}
}

func (c *UpdatedAt) Value() time.Time {
	return c.value
}

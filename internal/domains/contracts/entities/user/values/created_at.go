package values

import "time"

type CreatedAt struct {
	value time.Time
}

func NewCreatedAt(time time.Time) CreatedAt {
	return CreatedAt{
		value: time,
	}
}

func (c *CreatedAt) Value() time.Time {
	return c.value
}

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

func (c *CreatedAt) GetValue() time.Time {
	return c.value
}

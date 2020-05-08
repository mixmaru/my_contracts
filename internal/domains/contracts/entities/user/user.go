package user

import (
	"time"
)

type IUser interface {
	Id() IdInt
}

//////// エンティティ定義 ///////////
type User struct {
	id        IdInt
	createdAt CreatedAt
	updatedAt UpdatedAt
}

func (u *User) Id() IdInt {
	return u.id
}

//////// 値オブジェクト定義 ///////////
type IdInt struct {
	value int
}

func NewIdInt(id int) IdInt {
	return IdInt{
		value: id,
	}
}
func (i *IdInt) Value() int {
	return i.value
}

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

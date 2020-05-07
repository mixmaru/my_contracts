package values

type IdInt struct {
	id int
}

func NewIdInt(Id int) IdInt {
	return IdInt{
		id: Id,
	}
}

func (i *IdInt) Id() int {
	return i.id
}

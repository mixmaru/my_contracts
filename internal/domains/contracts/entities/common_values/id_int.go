package common_values

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

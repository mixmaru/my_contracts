package values

type Id struct {
	value int
}

func NewIdInt(id int) Id {
	return Id{
		value: id,
	}
}
func (i *Id) Value() int {
	return i.value
}

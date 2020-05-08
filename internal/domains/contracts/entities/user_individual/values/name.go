package values

type Name struct {
	value string
}

func NewName(name string) Name {
	return Name{
		value: name,
	}
}

func (n *Name) Value() string {
	return n.value
}

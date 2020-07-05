package tables

type EntitySetter interface {
	SetDataToEntity(entity interface{}) error
}

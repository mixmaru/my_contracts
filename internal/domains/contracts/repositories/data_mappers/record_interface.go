package data_mappers

type EntitySetter interface {
	SetDataToEntity(entity interface{}) error
}

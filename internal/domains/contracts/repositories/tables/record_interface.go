package tables

type IRecord interface {
	SetDataToEntity(entity interface{}) error
}

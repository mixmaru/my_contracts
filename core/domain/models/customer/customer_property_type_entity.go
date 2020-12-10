package customer

type PropertyType int8

const (
	PROPERTY_TYPE_STRING  PropertyType = 0
	PROPERTY_TYPE_NUMERIC PropertyType = 1
)

/*
自由に設定できるカスタマーの属性。
name = "性別", propertyType = PROPERTY_TYPE_STRING で、"女"とかを設定できるようになる
*/
type CustomerPropertyTypeEntity struct {
	id           int
	name         string
	propertyType PropertyType
}

func NewCustomerParamTypeEntity(name string, propertyType PropertyType) *CustomerPropertyTypeEntity {
	return &CustomerPropertyTypeEntity{name: name, propertyType: propertyType}
}

func NewCustomerParamTypeEntityWithData(id int, name string, propertyType PropertyType) *CustomerPropertyTypeEntity {
	entity := CustomerPropertyTypeEntity{}
	entity.LoadData(id, name, propertyType)
	return &entity
}

func (c *CustomerPropertyTypeEntity) Id() int {
	return c.id
}

func (c *CustomerPropertyTypeEntity) Name() string {
	return c.name
}

func (c *CustomerPropertyTypeEntity) PropertyType() PropertyType {
	return c.propertyType
}

func (c *CustomerPropertyTypeEntity) LoadData(id int, name string, propertyType PropertyType) {
	c.id = id
	c.name = name
	c.propertyType = propertyType
}

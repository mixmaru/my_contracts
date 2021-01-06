package customer

type CustomerEntity struct {
	id         int
	name       string
	properties map[string]interface{}
}

func NewCustomerEntity(name string, properties map[string]interface{}) *CustomerEntity {
	return &CustomerEntity{name: name, properties: properties}
}

func NewCustomerEntityWithData(id int, name string, properties map[string]interface{}) *CustomerEntity {
	entity := CustomerEntity{}
	entity.LoadData(id, name, properties)
	return &entity
}

func (c *CustomerEntity) LoadData(id int, name string, properties map[string]interface{}) {
	c.id = id
	c.name = name
	c.properties = properties
}

func (c *CustomerEntity) Id() int {
	return c.id
}

func (c *CustomerEntity) Name() string {
	return c.name
}

func (c *CustomerEntity) Properties() map[string]interface{} {
	return c.properties
}

package customer

type CustomerEntity struct {
	id             int
	name           string
	customerTypeId int
	properties     map[int]interface{}
}

func NewCustomerEntity(name string, customerTypeId int, properties map[int]interface{}) *CustomerEntity {
	return &CustomerEntity{
		name:           name,
		customerTypeId: customerTypeId,
		properties:     properties,
	}
}

func NewCustomerEntityWithData(id int, name string, customerTypeId int, properties map[int]interface{}) *CustomerEntity {
	entity := CustomerEntity{}
	entity.LoadData(id, name, customerTypeId, properties)
	return &entity
}

func (c *CustomerEntity) LoadData(id int, name string, customerTypeId int, properties map[int]interface{}) {
	c.id = id
	c.name = name
	c.customerTypeId = customerTypeId
	c.properties = properties
}

func (c *CustomerEntity) Id() int {
	return c.id
}

func (c *CustomerEntity) Name() string {
	return c.name
}

func (c *CustomerEntity) CustomerTypeId() int {
	return c.customerTypeId
}

func (c *CustomerEntity) Properties() map[int]interface{} {
	return c.properties
}

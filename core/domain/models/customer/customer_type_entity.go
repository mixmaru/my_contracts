package customer

type CustomerTypeEntity struct {
	id                   int
	name                 string
	customerParamTypeIds []int
}

func NewCustomerTypeEntity(name string, customerParamTypeIds []int) *CustomerTypeEntity {
	return &CustomerTypeEntity{name: name, customerParamTypeIds: customerParamTypeIds}
}

func NewCustomerTypeEntityWithData(id int, name string, customerParamTypeIds []int) *CustomerTypeEntity {
	entity := CustomerTypeEntity{}
	entity.LoadData(id, name, customerParamTypeIds)
	return &entity
}

func (c *CustomerTypeEntity) Id() int {
	return c.id
}

func (c *CustomerTypeEntity) Name() string {
	return c.name
}

func (c *CustomerTypeEntity) CustomerPropertyTypeIds() []int {
	return c.customerParamTypeIds
}

func (c *CustomerTypeEntity) LoadData(id int, name string, customerParamTypeIds []int) {
	c.id = id
	c.name = name
	c.customerParamTypeIds = customerParamTypeIds
}

package customer

type CustomerTypeEntity struct {
	id                   int
	name                 string
	customerParamTypeIds []int
}

func NewCustomerTypeEntity(name string, customerParamTypeIds []int) *CustomerTypeEntity {
	return &CustomerTypeEntity{name: name, customerParamTypeIds: customerParamTypeIds}
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

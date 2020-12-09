package customer

type CustomerTypeEntity struct {
	id                 int
	name               string
	customerParamTypes map[int]*CustomerParamTypeEntity
}

func NewCustomerTypeEntity(name string, customerParamTypes map[int]*CustomerParamTypeEntity) *CustomerTypeEntity {
	return &CustomerTypeEntity{name: name, customerParamTypes: customerParamTypes}
}

func (c *CustomerTypeEntity) Id() int {
	return c.id
}

func (c *CustomerTypeEntity) Name() string {
	return c.name
}

//// 外部からいじられないようにデータコピーして渡す
func (c *CustomerTypeEntity) CustomerParamTypes() map[int]*CustomerParamTypeEntity {
	retParamTypes := make(map[int]*CustomerParamTypeEntity, len(c.customerParamTypes))
	for _, paramType := range c.customerParamTypes {
		entity := *paramType
		retParamTypes[entity.id] = &entity
	}
	return retParamTypes
}

const (
	PARAM_TYPE_STRING  string = "string"
	PARAM_TYPE_NUMERIC string = "numeric"
)

/*
自由に設定できるカスタマーの属性。
name = "性別", paramType = PARAM_TYPE_STRING で、"女"とかを設定できるようになる
*/
type CustomerParamTypeEntity struct {
	id        int
	name      string
	paramType string
}

func NewCustomerParamTypeEntity(name string, paramType string) *CustomerParamTypeEntity {
	return &CustomerParamTypeEntity{name: name, paramType: paramType}
}

func (c CustomerParamTypeEntity) Id() int {
	return c.id
}

func (c CustomerParamTypeEntity) Name() string {
	return c.name
}

func (c CustomerParamTypeEntity) ParamType() string {
	return c.paramType
}

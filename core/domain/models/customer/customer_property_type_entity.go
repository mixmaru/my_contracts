package customer

const (
	PROPERTY_TYPE_STRING  string = "string"
	PROPERTY_TYPE_NUMERIC string = "numeric"
)

/*
自由に設定できるカスタマーの属性。
name = "性別", paramType = PROPERTY_TYPE_STRING で、"女"とかを設定できるようになる
*/
type CustomerPropertyTypeEntity struct {
	id        int
	name      string
	paramType string
}

func NewCustomerParamTypeEntity(name string, paramType string) *CustomerPropertyTypeEntity {
	return &CustomerPropertyTypeEntity{name: name, paramType: paramType}
}

func (c CustomerPropertyTypeEntity) Id() int {
	return c.id
}

func (c CustomerPropertyTypeEntity) Name() string {
	return c.name
}

func (c CustomerPropertyTypeEntity) ParamType() string {
	return c.paramType
}

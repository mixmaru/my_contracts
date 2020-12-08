package customer

import (
	"time"
)

type CustomerType struct {
	id                 int
	name               string
	customerParamTypes map[int]*CustomerParamType
	createdAt          time.Time
	updatedAt          time.Time
}

func NewCustomerType(id int, name string, customerParamTypes map[int]*CustomerParamType) *CustomerType {
	return &CustomerType{id: id, name: name, customerParamTypes: customerParamTypes}
}

func (c *CustomerType) Id() int {
	return c.id
}

func (c *CustomerType) Name() string {
	return c.name
}

//// 外部からいじられないようにデータコピーして渡す
func (c *CustomerType) CustomerParamTypes() map[int]*CustomerParamType {
	retParamTypes := make(map[int]*CustomerParamType, len(c.customerParamTypes))
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
type CustomerParamType struct {
	id        int
	name      string
	paramType string
	createdAt time.Time
	updatedAt time.Time
}

func NewCustomerParamType(id int, name string, paramType string) *CustomerParamType {
	return &CustomerParamType{id: id, name: name, paramType: paramType}
}

func (c CustomerParamType) Id() int {
	return c.id
}

func (c CustomerParamType) Name() string {
	return c.name
}

func (c CustomerParamType) ParamType() string {
	return c.paramType
}

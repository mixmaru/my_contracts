package customer

type CustomerDto struct {
	Id             int
	CustomerTypeId int
	Name           string
	Properties     PropertyDto
}

type PropertyDto map[string]string

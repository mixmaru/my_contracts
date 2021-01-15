package customer

import "github.com/mixmaru/my_contracts/core/domain/models/customer"

type CustomerDto struct {
	Id             int
	CustomerTypeId int
	Name           string
	Properties     PropertyDto
}

type PropertyDto map[int]interface{}

func NewCustomerDtoFromEntity(entity *customer.CustomerEntity) CustomerDto {
	dto := CustomerDto{}
	dto.Id = entity.Id()
	dto.Name = entity.Name()
	dto.CustomerTypeId = entity.CustomerTypeId()
	dto.Properties = entity.Properties()
	return dto
}

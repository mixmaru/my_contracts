package data_transfer_objects

import "github.com/mixmaru/my_contracts/internal/domains/contracts/entities"

const userIndividualDtoType = "individual"

type UserIndividualDto struct {
	Name string
	Type string
	BaseDto
}

func NewUserIndividualDtoFromEntity(entity *entities.UserIndividualEntity) UserIndividualDto {
	dto := UserIndividualDto{}
	dto.Type = userIndividualDtoType
	dto.Id = entity.Id()
	dto.Name = entity.Name()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	return dto
}

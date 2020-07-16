package data_transfer_objects

import "github.com/mixmaru/my_contracts/internal/domains/contracts/entities"

type UserIndividualDto struct {
	Name string
	BaseDto
}

func NewUserIndividualDtoFromEntity(entity *entities.UserIndividualEntity) UserIndividualDto {
	dto := UserIndividualDto{}
	dto.Id = entity.Id()
	dto.Name = entity.Name()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	return dto
}

package data_transfer_objects

import "github.com/mixmaru/my_contracts/internal/domains/contracts/entities"

type UserCorporationDto struct {
	ContactPersonName string
	PresidentName     string
	BaseDto
}

func NewUserCorporationDtoFromEntity(entity *entities.UserCorporationEntity) UserCorporationDto {
	dto := UserCorporationDto{}
	dto.Id = entity.Id()
	dto.ContactPersonName = entity.ContactPersonName()
	dto.PresidentName = entity.PresidentName()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	return dto
}

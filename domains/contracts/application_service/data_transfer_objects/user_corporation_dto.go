package data_transfer_objects

import "github.com/mixmaru/my_contracts/domains/contracts/entities"

const userCorporationDtoType = "corporation"

type UserCorporationDto struct {
	CorporationName   string
	ContactPersonName string
	PresidentName     string
	Type              string
	BaseDto
}

func NewUserCorporationDtoFromEntity(entity *entities.UserCorporationEntity) UserCorporationDto {
	dto := UserCorporationDto{}
	dto.Type = userCorporationDtoType
	dto.Id = entity.Id()
	dto.CorporationName = entity.CorporationName()
	dto.ContactPersonName = entity.ContactPersonName()
	dto.PresidentName = entity.PresidentName()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	return dto
}

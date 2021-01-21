package users

import (
	"github.com/mixmaru/my_contracts/core/domain/models/user"
	"time"
)

const userCorporationDtoType = "corporation"

type UserCorporationDto struct {
	Id                int
	CorporationName   string
	ContactPersonName string
	PresidentName     string
	Type              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewUserCorporationDtoFromEntity(entity *user.UserCorporationEntity) UserCorporationDto {
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

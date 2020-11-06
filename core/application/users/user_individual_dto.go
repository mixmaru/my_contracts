package users

import (
	entity "github.com/mixmaru/my_contracts/core/domain/models/user"
	"time"
)

const userIndividualDtoType = "individual"

type UserIndividualDto struct {
	Id        int
	Name      string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserIndividualDtoFromEntity(entity *entity.UserIndividualEntity) UserIndividualDto {
	dto := UserIndividualDto{}
	dto.Type = userIndividualDtoType
	dto.Id = entity.Id()
	dto.Name = entity.Name()
	dto.CreatedAt = entity.CreatedAt()
	dto.UpdatedAt = entity.UpdatedAt()
	return dto
}

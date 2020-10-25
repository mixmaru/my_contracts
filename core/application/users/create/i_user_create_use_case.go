package create

import (
	"github.com/mixmaru/my_contracts/core/application/users"
)

////// 個人顧客
type IUserIndividualCreateUseCase interface {
	Handle(request *UserIndividualCreateUseCaseRequest) (*UserIndividualCreateUseCaseResponse, error)
}

type UserIndividualCreateUseCaseRequest struct {
	Name string
}

func NewUserIndividualCreateUseCaseRequest(name string) *UserIndividualCreateUseCaseRequest {
	return &UserIndividualCreateUseCaseRequest{Name: name}
}

type UserIndividualCreateUseCaseResponse struct {
	UserDto          users.UserIndividualDto
	ValidationErrors map[string][]string
}

func NewUserIndividualCreateUseCaseResponse(userDto users.UserIndividualDto, validationErrors map[string][]string) *UserIndividualCreateUseCaseResponse {
	return &UserIndividualCreateUseCaseResponse{
		UserDto:          userDto,
		ValidationErrors: validationErrors,
	}
}

////// 法人顧客
type IUserCorporationCreateUseCase interface {
	Handle(request *UserCorporationCreateUseCaseRequest) (*UserCorporationCreateUseCaseResponse, error)
}

type UserCorporationCreateUseCaseRequest struct {
	CorporationName   string
	ContactPersonName string
	PresidentName     string
}

func NewUserCorporationCreateUseCaseRequest(corporationName, contractPersonName, presidentName string) *UserCorporationCreateUseCaseRequest {
	return &UserCorporationCreateUseCaseRequest{
		CorporationName:   corporationName,
		ContactPersonName: contractPersonName,
		PresidentName:     presidentName,
	}
}

type UserCorporationCreateUseCaseResponse struct {
	UserDto          users.UserCorporationDto
	ValidationErrors map[string][]string
}

func NewUserCorporationCreateUseCaseResponse(userDto users.UserCorporationDto, validationErrors map[string][]string) *UserCorporationCreateUseCaseResponse {
	return &UserCorporationCreateUseCaseResponse{
		UserDto:          userDto,
		ValidationErrors: validationErrors,
	}
}

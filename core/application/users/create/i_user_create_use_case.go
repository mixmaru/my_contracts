package create

import "github.com/mixmaru/my_contracts/core/application/users"

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

package create

import "github.com/mixmaru/my_contracts/core/application/users"

type IUserIndividualCreateUseCase interface {
	Handle(request *UserIndividualCreateUseCaseRequest) (*UserIndividualCreateUseCaseResponse, error)
}

type UserIndividualCreateUseCaseRequest struct {
	Name string
}

func NewUserCreateUseCaseRequest(name string) *UserIndividualCreateUseCaseRequest {
	return &UserIndividualCreateUseCaseRequest{Name: name}
}

type UserIndividualCreateUseCaseResponse struct {
	UserDto         users.UserIndividualDto
	ValidationError map[string][]string
}

func NewUserCreateUseCaseResponse(userDto users.UserIndividualDto, validationError map[string][]string) *UserIndividualCreateUseCaseResponse {
	return &UserIndividualCreateUseCaseResponse{
		UserDto:         userDto,
		ValidationError: validationError,
	}
}

package get

type IUserGetUseCase interface {
	Handle(request *UserGetUseCaseRequest) (*UserGetUseCaseResponse, error)
}

type UserGetUseCaseRequest struct {
	UserId int
}

func NewUserGetUseCaseRequest(userId int) *UserGetUseCaseRequest {
	return &UserGetUseCaseRequest{
		UserId: userId,
	}
}

type UserGetUseCaseResponse struct {
	UserDto interface{}
}

func NewUserGetUseCaseResponse(userDto interface{}) *UserGetUseCaseResponse {
	return &UserGetUseCaseResponse{
		UserDto: userDto,
	}
}

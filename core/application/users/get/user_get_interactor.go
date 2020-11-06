package get

import (
	"github.com/mixmaru/my_contracts/core/application/users"
	entities2 "github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/pkg/errors"
)

type UserGetInteractor struct {
	userRepository users.IUserRepository
}

func NewUserGetInteractor(userRepository users.IUserRepository) *UserGetInteractor {
	return &UserGetInteractor{userRepository: userRepository}
}

func (u *UserGetInteractor) Handle(request *UserGetUseCaseRequest) (*UserGetUseCaseResponse, error) {
	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	response := &UserGetUseCaseResponse{}

	gotUser, err := u.userRepository.GetUserById(request.UserId, conn)
	if err != nil {
		return nil, err
	}

	if gotUser == nil {
		// データがない場合
		return response, nil
	}

	switch gotUser.(type) {
	case *entities2.UserIndividualEntity:
		response.UserDto = users.NewUserIndividualDtoFromEntity(gotUser.(*entities2.UserIndividualEntity))
		return response, nil
	case *entities2.UserCorporationEntity:
		response.UserDto = users.NewUserCorporationDtoFromEntity(gotUser.(*entities2.UserCorporationEntity))
		return response, nil
	default:
		return nil, errors.Errorf("考慮していないtypeが来た。type: %T, userId: %v", gotUser, request.UserId)
	}
}

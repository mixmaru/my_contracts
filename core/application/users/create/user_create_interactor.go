package create

import (
	"fmt"
	"github.com/mixmaru/my_contracts/core/application/users"
	entities "github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/mixmaru/my_contracts/core/domain/validators"
	"github.com/mixmaru/my_contracts/domains/contracts/entities/values"
	"github.com/pkg/errors"
)

type UserIndividualCreateInteractor struct {
	userRepository users.IUserRepository
}

func NewUserIndividualCreateInteractor(userRepository users.IUserRepository) *UserIndividualCreateInteractor {
	return &UserIndividualCreateInteractor{userRepository: userRepository}
}

func (u *UserIndividualCreateInteractor) Handle(request *UserIndividualCreateUseCaseRequest) (*UserIndividualCreateUseCaseResponse, error) {
	// 入力値バリデーション
	validationErrors, err := userIndividualValidation(request.Name)
	if err != nil {
		return NewUserCreateUseCaseResponse(users.UserIndividualDto{}, nil), err
	}
	if len(validationErrors) > 0 {
		return NewUserCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), nil
	}

	// エンティティ作成
	userEntity, err := entities.NewUserIndividualEntity(request.Name)
	if err != nil {
		return NewUserCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), err
	}

	// トランザクション開始
	dbMap, err := db.GetConnection()
	if err != nil {
		return NewUserCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), err
	}
	defer dbMap.Db.Close()
	tran, err := dbMap.Begin()
	if err != nil {
		return NewUserCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), errors.Wrapf(err, "トランザクション開始失敗")
	}

	// リポジトリ使って保存
	savedId, err := u.userRepository.SaveUserIndividual(userEntity, tran)
	if err != nil {
		return NewUserCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), err
	}

	// 保存データ再読込
	savedUserEntity, err := u.userRepository.GetUserIndividualById(savedId, tran)
	if err != nil {
		return NewUserCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), err
	}

	// コミット
	err = tran.Commit()
	if err != nil {
		return NewUserCreateUseCaseResponse(users.UserIndividualDto{}, validationErrors), errors.Wrapf(err, "コミット失敗")
	}
	userDto := users.NewUserIndividualDtoFromEntity(savedUserEntity)
	return NewUserCreateUseCaseResponse(userDto, validationErrors), nil
}

func userIndividualValidation(name string) (validationErrors map[string][]string, err error) {
	validationErrors = map[string][]string{}

	// 個人顧客名バリデーション
	nameValidErrors := values.NameValidate(name)
	if len(nameValidErrors) > 0 {
		validationErrors["name"] = []string{}
	}
	for _, validErr := range nameValidErrors {
		var errorMessage string
		switch validErr {
		case validators.EmptyStringValidError:
			errorMessage = "空です"
		case validators.OverLengthStringValidError:
			errorMessage = fmt.Sprintf("%v文字より多いです", values.NameMaxLength)
		default:
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。name: %v, validErrorText: %v", name, validators.ValidErrorText(validErr)))
		}
		validationErrors["name"] = append(validationErrors["name"], errorMessage)
	}

	return validationErrors, nil
}

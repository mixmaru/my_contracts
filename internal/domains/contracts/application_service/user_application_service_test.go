package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	user_repository "github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserApplication_NewUserApplicationService(t *testing.T) {
	// リポジトリモックを用意する
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)

	userApp := NewUserApplicationService(userRepositoryMock)
	assert.IsType(t, &UserApplicationService{}, userApp)
}

func TestUserApplicationService_RegisterUserIndividual(t *testing.T) {
	// インスタンス化
	userApp := NewUserApplicationService(&user_repository.Repository{})

	userId, err := userApp.RegisterUserIndividual("個人太郎")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, userId)

	// todo: IDでuser情報を取得してチェックする
}

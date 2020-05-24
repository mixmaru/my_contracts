package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserApplicationService_NewUserApplicationService(t *testing.T) {
	// インスタンス化テスト
	userApp := NewUserApplicationService()
	assert.IsType(t, &UserApplicationService{}, userApp)
}

func TestUserApplicationService_NewUserApplicationServiceWithMock(t *testing.T) {
	// mock作成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)

	// インスタンス化テスト
	userApp := NewUserApplicationServiceWithMock(userRepositoryMock)
	assert.IsType(t, &UserApplicationService{}, userApp)
}

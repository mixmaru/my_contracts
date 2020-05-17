package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
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
	// リポジトリモックを用意する
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)

	//// チェック用にメソッドを設定。
	//// 来る想定の引数でメソッドが実行されるかチェックする
	// [個人太郎]がセットされたuserEntityとトランザクション(nilではない)が渡されて1回だけ実行されるはず
	userEntity := user.NewUserIndividualEntity()
	userEntity.SetName("個人太郎")
	userRepositoryMock.EXPECT().SaveUserIndividual(gomock.Eq(userEntity), gomock.Not(nil)).Return(nil).Times(1)

	// モックを渡してインスタンス化
	userApp := NewUserApplicationService(userRepositoryMock)

	userId, err := userApp.RegisterUserIndividual("個人太郎")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, userId)

}

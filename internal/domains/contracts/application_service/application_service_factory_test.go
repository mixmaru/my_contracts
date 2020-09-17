package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	"github.com/stretchr/testify/assert"
	"testing"
)

////// UserApplicationService
func TestApplicationServiceFactory_NewUserApplicationService(t *testing.T) {
	// インスタンス化テスト
	userApp := NewUserApplicationService()
	assert.IsType(t, &UserApplicationService{}, userApp)
}

func TestApplicationServiceFactory_NewUserApplicationServiceWithMock(t *testing.T) {
	// mock作成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepositoryMock := mock_interfaces.NewMockIUserRepository(ctrl)

	// インスタンス化テスト
	userApp := NewUserApplicationServiceWithMock(userRepositoryMock)
	assert.IsType(t, &UserApplicationService{}, userApp)
}

////// ProductApplicationService
func TestApplicationServiceFactory_NewProductApplicationService(t *testing.T) {
	// インスタンス化テスト
	productApp := NewProductApplicationService()
	assert.IsType(t, &ProductApplicationService{}, productApp)
}

func TestApplicationServiceFactory_NewProductApplicationServiceWithMock(t *testing.T) {
	// mock作成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	productRepositoryMock := mock_interfaces.NewMockIProductRepository(ctrl)

	// インスタンス化テスト
	productApp := NewProductApplicationServiceWithMock(productRepositoryMock)
	assert.IsType(t, &ProductApplicationService{}, productApp)

}

////// ContractApplicationService
func TestApplicationServiceFactory_NewContractApplicationService(t *testing.T) {
	// インスタンス化テスト
	app := NewContractApplicationService()
	assert.IsType(t, &ContractApplicationService{}, app)
}

func TestApplicationServiceFactory_NewContractApplicationServiceWithMock(t *testing.T) {
	// mock作成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repositoryMock := mock_interfaces.NewMockIContractRepository(ctrl)

	// インスタンス化テスト
	app := NewContractApplicationServiceWithMock(repositoryMock)
	assert.IsType(t, &ContractApplicationService{}, app)
}

////// BillApplicationService
func TestApplicationServiceFactory_NewBillApplicationService(t *testing.T) {
	// インスタンス化テスト
	app := NewBillApplicationService()
	assert.IsType(t, &BillApplicationService{}, app)
}

func TestApplicationServiceFactory_NewBillApplicationServiceWithMock(t *testing.T) {
	// mock作成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	product := mock_interfaces.NewMockIProductRepository(ctrl)
	contract := mock_interfaces.NewMockIContractRepository(ctrl)
	rightToUse := mock_interfaces.NewMockIRightToUseRepository(ctrl)
	bill := mock_interfaces.NewMockIBillRepository(ctrl)

	// インスタンス化テスト
	app := NewBillApplicationServiceWithMock(product, contract, rightToUse, bill)
	assert.IsType(t, &BillApplicationService{}, app)
}

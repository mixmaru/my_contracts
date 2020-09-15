package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/domain_service"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/pkg/errors"
	"time"
)

type BillApplicationService struct {
	productRepository    interfaces.IProductRepository
	contractRepository   interfaces.IContractRepository
	rightToUseRepository interfaces.IRightToUseRepository
	billRepository       interfaces.IBillRepository
}

// 渡した指定日を実行日として請求の実行をする
func (b *BillApplicationService) ExecuteBilling(executeDate time.Time) error {
	db, err := db_connection.GetConnection()
	if err != nil {
		return err
	}
	defer db.Db.Close()
	tran, err := db.Begin()
	billDomain := domain_service.NewBillingCalculatorDomainService(b.productRepository, b.contractRepository, b.rightToUseRepository, b.billRepository)
	err = billDomain.ExecuteBilling(executeDate, tran)
	if err != nil {
		tran.Rollback()
		return err
	}
	err = tran.Commit()
	if err != nil {
		return errors.Wrapf(err, "コミットに失敗しました。")
	}
	return nil
}

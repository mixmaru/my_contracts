package application_service

import (
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/domains/contracts/domain_service"
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/db_connection"
	"time"
)

type BillApplicationService struct {
	productRepository  interfaces.IProductRepository
	contractRepository interfaces.IContractRepository
	billRepository     interfaces.IBillRepository
}

// 渡した指定日を実行日として請求の実行をする
func (b *BillApplicationService) ExecuteBilling(executeDate time.Time) ([]data_transfer_objects.BillDto, error) {
	db, err := db_connection.GetConnection()
	if err != nil {
		return nil, err
	}
	defer db.Db.Close()

	// 請求対象使用権をもつ契約を取得する
	contracts, err := b.contractRepository.GetBillingTargetByBillingDate(executeDate, db)
	if err != nil {
		return nil, err
	}

	// for で回しながら請求を実行
	billDomain := domain_service.NewBillingCalculatorDomainService(b.productRepository, b.contractRepository, b.billRepository)
	retBillDtos := make([]data_transfer_objects.BillDto, 0, len(contracts))
	for _, contract := range contracts {
		billAgg := entities.NewBillingAggregation(executeDate, contract.UserId())
		for _, rightToUse := range contract.RightToUses() {
			// 未請求かつ、validFromと課金開始日がexecuteDate以前の物を請求実行
			if domain_service.IsBillingTarget(executeDate, contract.BillingStartDate(), rightToUse) {
				amount, err := billDomain.BillingAmount(rightToUse, contract.BillingStartDate(), db)
				if err != nil {
					return nil, err
				}
				err = billAgg.AddBillDetail(entities.NewBillingDetailEntity(rightToUse.Id(), amount))
				if err != nil {
					return nil, err
				}
			}
		}
		// 保存
		savedBillId, err := b.billRepository.Create(billAgg, db)
		if err != nil {
			return nil, err
		}
		// 再取得してdtoにする
		savedBill, err := b.billRepository.GetById(savedBillId, db)
		billDto, err := data_transfer_objects.NewBillDtoFromEntity(savedBill)
		if err != nil {
			return retBillDtos, err
		}
		retBillDtos = append(retBillDtos, billDto)
	}

	// 請求データを返す
	return retBillDtos, nil
}

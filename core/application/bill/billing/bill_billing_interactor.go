package billing

import (
	"github.com/mixmaru/my_contracts/core/application/bill"
	"github.com/mixmaru/my_contracts/core/application/contracts"
	"github.com/mixmaru/my_contracts/core/application/products"
	bill_entity "github.com/mixmaru/my_contracts/core/domain/models/bill"
	"github.com/mixmaru/my_contracts/core/domain/services"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/pkg/errors"
)

type BillBillingInteractor struct {
	productRepository  products.IProductRepository
	contractRepository contracts.IContractRepository
	billRepository     bill.IBillRepository
}

func NewBillBillingInteractor(productRepository products.IProductRepository, contractRepository contracts.IContractRepository, billRepository bill.IBillRepository) *BillBillingInteractor {
	return &BillBillingInteractor{productRepository: productRepository, contractRepository: contractRepository, billRepository: billRepository}
}

/*
渡した指定日を実行日として請求の実行をする

処理途中でエラーが発生しても、処理完了した部分までは処理された状態になり、データが返却される
*/
func (b *BillBillingInteractor) Handle(request *BillBillingUseCaseRequest) (*BillBillingUseCaseResponse, error) {
	response := &BillBillingUseCaseResponse{}

	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	// 請求対象使用権をもつ契約を取得する
	contracts, err := b.contractRepository.GetBillingTargetByBillingDate(request.ExecuteDate, conn)
	if err != nil {
		return nil, err
	}

	// for で回しながら請求を実行
	billDomain := services.NewBillService(b.productRepository, b.contractRepository, b.billRepository)
	response.BillDtos = make([]bill.BillDto, 0, len(contracts))
	for _, contract := range contracts {
		billEntity := bill_entity.NewBillEntity(request.ExecuteDate, contract.UserId())
		// トランザクション
		tran, err := conn.Begin()
		if err != nil {
			return response, errors.Wrapf(err, "トランザクション開始失敗")
		}
		for _, rightToUse := range contract.RightToUses() {
			// 未請求かつ、validFromと課金開始日がexecuteDate以前の物を請求実行
			if services.IsBillingTarget(request.ExecuteDate, contract.BillingStartDate(), rightToUse) {
				amount, err := billDomain.BillingAmount(rightToUse, contract.BillingStartDate(), tran)
				if err != nil {
					return response, err
				}
				err = billEntity.AddBillDetail(bill_entity.NewBillDetailEntity(rightToUse.Id(), amount))
				if err != nil {
					return response, err
				}
			}
		}
		// 保存
		savedBillId, err := b.billRepository.Create(billEntity, tran)
		if err != nil {
			return response, err
		}
		// コミットする
		err = tran.Commit()
		if err != nil {
			return response, errors.Wrapf(err, "トランザクションコミット失敗")
		}

		// 再取得してdtoにする
		savedBill, err := b.billRepository.GetById(savedBillId, conn)
		billDto, err := bill.NewBillDtoFromEntity(savedBill)
		if err != nil {
			return response, err
		}
		response.BillDtos = append(response.BillDtos, billDto)
	}

	// 請求データを返す
	return response, nil
}

package entities

import (
	"time"
)

type RightToUseEntity struct {
	BaseEntity
	contractId int
	validFrom  time.Time
	validTo    time.Time
}

func NewRightToUseEntity(contractId int, validFrom, validTo time.Time) *RightToUseEntity {
	return &RightToUseEntity{
		contractId: contractId,
		validFrom:  validFrom,
		validTo:    validTo,
	}
}

func (r *RightToUseEntity) ContractId() int {
	return r.contractId
}

func (r *RightToUseEntity) ValidFrom() time.Time {
	return r.validFrom
}

func (r *RightToUseEntity) ValidTo() time.Time {
	return r.validTo
}

//func NewContractEntityWithData(id, userId, productId int, contractDate, billingStartDate, createdAt, updatedAt time.Time) (*ContractEntity, error) {
//	entity := ContractEntity{}
//	err := entity.LoadData(id, userId, productId, contractDate, billingStartDate, createdAt, updatedAt)
//	if err != nil {
//		return nil, err
//	}
//	return &entity, nil
//}
//
//func (c *ContractEntity) UserId() int {
//	return c.userId
//}
//
//func (c *ContractEntity) ProductId() int {
//	return c.productId
//}
//
//func (c *ContractEntity) ContractDate() time.Time {
//	return c.contractDate
//}
//
//func (c *ContractEntity) BillingStartDate() time.Time {
//	return c.billingStartDate
//}
//
///*
//対象日以下の最大の課金開始日（直近の課金開始日）を返す
//*/
//func (c *ContractEntity) LastBillingStartDate(targetDate time.Time) time.Time {
//	billingDate := c.billingStartDate
//	for true {
//		nextBillingStartDate := billingDate.AddDate(0, 1, 0)
//		if nextBillingStartDate.After(targetDate) {
//			return billingDate
//		} else {
//			billingDate = nextBillingStartDate
//		}
//	}
//	return billingDate // ここにはこないはず
//}
//
////// 保持データをセットし直す
//func (c *ContractEntity) LoadData(id, userId, productId int, contractDate, billingStartDate, createdAt, updatedAt time.Time) error {
//	c.id = id
//	c.userId = userId
//	c.productId = productId
//	c.contractDate = contractDate
//	c.billingStartDate = billingStartDate
//
//	c.createdAt = createdAt
//	c.updatedAt = updatedAt
//	return nil
//}

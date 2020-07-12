package entities

import (
	"time"
)

type ContractEntity struct {
	BaseEntity
	userId    int
	productId int
}

func NewContractEntity(userId int, productId int) (*ContractEntity, error) {
	return &ContractEntity{
		userId:    userId,
		productId: productId,
	}, nil
}

//func NewProductEntityWithData(id int, name string, price string, createdAt, updatedAt time.Time) (*ProductEntity, error) {
//	productEntity := ProductEntity{}
//	err := productEntity.LoadData(id, name, price, createdAt, updatedAt)
//	if err != nil {
//		return nil, err
//	}
//	return &productEntity, nil
//}

func (c *ContractEntity) UserId() int {
	return c.userId
}

func (c *ContractEntity) ProductId() int {
	return c.productId
}

//// 保持データをセットし直す
func (c *ContractEntity) LoadData(id int, userId int, productId int, createdAt time.Time, updatedAt time.Time) error {
	c.id = id
	c.userId = userId
	c.productId = productId
	c.createdAt = createdAt
	c.updatedAt = updatedAt
	return nil
}

package entities

import (
	"time"
)

type ContractEntity struct {
	BaseEntity
	userId    int
	productId int
}

func NewContractEntity(userId int, productId int) *ContractEntity {
	return &ContractEntity{
		userId:    userId,
		productId: productId,
	}
}

func NewContractEntityWithData(id int, userId int, productId int, createdAt, updatedAt time.Time) (*ContractEntity, error) {
	entity := ContractEntity{}
	err := entity.LoadData(id, userId, productId, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

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

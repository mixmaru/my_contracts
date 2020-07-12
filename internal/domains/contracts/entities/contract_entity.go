package entities

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
//func (p *ProductEntity) LoadData(id int, name string, price string, createdAt time.Time, updatedAt time.Time) error {
//	nameValue, err := values.NewProductNameValue(name)
//	if err != nil {
//		return err
//	}
//	priceValue, err := values.NewProductPriceValue(price)
//
//	p.id = id
//	p.name = nameValue
//	p.price = priceValue
//	p.createdAt = createdAt
//	p.updatedAt = updatedAt
//	return nil
//}

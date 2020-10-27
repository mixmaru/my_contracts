package contract

import (
	"github.com/pkg/errors"
	"time"
)

type ContractEntity struct {
	id                 int
	userId             int
	productId          int
	contractDate       time.Time
	billingStartDate   time.Time
	rightToUseEntities []*rightToUseEntity // アクティブな使用権
	toArchive          []*rightToUseEntity // ArchiveRightToUseByIdで指定されたやつ。リポジトリで処理されるのを期待
	createdAt          time.Time
	updatedAt          time.Time
}

func NewContractEntity(userId int, productId int, contractDate, billingStartDate time.Time, rightToUses []*rightToUseEntity) *ContractEntity {
	return &ContractEntity{
		userId:             userId,
		productId:          productId,
		contractDate:       contractDate,
		billingStartDate:   billingStartDate,
		rightToUseEntities: rightToUses,
	}
}

func NewContractEntityWithData(id, userId, productId int, contractDate, billingStartDate, createdAt, updatedAt time.Time, rightToUses []*rightToUseEntity) (*ContractEntity, error) {
	entity := ContractEntity{}
	err := entity.LoadData(id, userId, productId, contractDate, billingStartDate, createdAt, updatedAt, rightToUses)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (c *ContractEntity) Id() int {
	return c.id
}

func (c *ContractEntity) UserId() int {
	return c.userId
}

func (c *ContractEntity) ProductId() int {
	return c.productId
}

func (c *ContractEntity) ContractDate() time.Time {
	return c.contractDate
}

func (c *ContractEntity) BillingStartDate() time.Time {
	return c.billingStartDate
}

func (c *ContractEntity) CreatedAt() time.Time {
	return c.createdAt
}

func (c *ContractEntity) UpdatedAt() time.Time {
	return c.updatedAt
}

// 外部からいじられないようにデータコピーして渡す
func (c *ContractEntity) RightToUses() []*rightToUseEntity {
	retEntities := make([]*rightToUseEntity, 0, len(c.rightToUseEntities))
	for _, rightToUse := range c.rightToUseEntities {
		entity := *rightToUse
		retEntities = append(retEntities, &entity)
	}
	return retEntities
}

// 次期使用権の追加
func (c *ContractEntity) AddNextTermRightToUses(rightToUse *rightToUseEntity) {
	c.rightToUseEntities = append(c.rightToUseEntities, rightToUse)
}

/*
対象日以下の最大の課金開始日（直近の課金開始日）を返す
*/
func (c *ContractEntity) LastBillingStartDate(targetDate time.Time) time.Time {
	billingDate := c.billingStartDate
	for true {
		nextBillingStartDate := billingDate.AddDate(0, 1, 0)
		if nextBillingStartDate.After(targetDate) {
			return billingDate
		} else {
			billingDate = nextBillingStartDate
		}
	}
	return billingDate // ここにはこないはず
}

//// 保持データをセットし直す
func (c *ContractEntity) LoadData(id, userId, productId int, contractDate, billingStartDate, createdAt, updatedAt time.Time, rightToUses []*rightToUseEntity) error {
	c.id = id
	c.userId = userId
	c.productId = productId
	c.contractDate = contractDate
	c.billingStartDate = billingStartDate
	c.rightToUseEntities = rightToUses

	c.createdAt = createdAt
	c.updatedAt = updatedAt
	return nil
}

/*
指定idの使用権をアーカイブ行きにする
*/
func (c *ContractEntity) ArchiveRightToUseById(rightToUseId int) error {
	for i, rightToUse := range c.rightToUseEntities {
		if rightToUse.Id() == rightToUseId {
			c.toArchive = append(c.toArchive, rightToUse)
			c.rightToUseEntities = remove(c.rightToUseEntities, i)
			return nil
		}
	}
	return errors.Errorf("指定Idの使用権が存在しない。rightToUseId: %v", rightToUseId)
}

func remove(slice []*rightToUseEntity, s int) []*rightToUseEntity {
	return append(slice[:s], slice[s+1:]...)
}

/*
指定のValidTo(有効期限終了日時)以前に期限を迎えている使用権をアーカイブ行きにする
*/
func (c *ContractEntity) ArchiveRightToUseByValidTo(ValidTo time.Time) {
	target := c.rightToUseEntities
	c.rightToUseEntities = []*rightToUseEntity{}
	for _, rightToUse := range target {
		if rightToUse.validTo.Equal(ValidTo) || rightToUse.validTo.Before(ValidTo) {
			c.toArchive = append(c.toArchive, rightToUse)
		} else {
			c.rightToUseEntities = append(c.rightToUseEntities, rightToUse)
		}
	}
	return
}

/*
アーカイブ行き指定された使用権Idのスライスを返す
*/
func (c *ContractEntity) GetToArchiveRightToUseIds() []int {
	retIds := make([]int, 0, len(c.toArchive))
	for _, entity := range c.toArchive {
		retIds = append(retIds, entity.Id())
	}
	return retIds
}

/*
アーカイブ行き指定された使用権のスライスを返す
*/
func (c *ContractEntity) GetToArchiveRightToUses() []*rightToUseEntity {
	retEntities := make([]*rightToUseEntity, 0, len(c.toArchive))
	for _, rightToUse := range c.toArchive {
		entity := *rightToUse // 外からいじられないようにコピーしてわたす
		retEntities = append(retEntities, &entity)
	}
	return retEntities
}

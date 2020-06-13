package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
)

type ProductApplicationService struct {
	productRepository interfaces.IProductRepository
}

func (p *ProductApplicationService) Register(name string, price decimal.Decimal) (interface{}, ValidationErrors, error) {
	// バリデーション実行
	// entityを作成
	// リポジトリで保存
	// dtoに詰める
	// 返却
	return nil, nil, nil
}

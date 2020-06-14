package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/pkg/errors"
)

type ProductApplicationService struct {
	productRepository interfaces.IProductRepository
}

func (p *ProductApplicationService) Register(name string, price decimal.Decimal) (data_transfer_objects.ProductDto, ValidationErrors, error) {
	// バリデーション実行

	// entityを作成
	entity := entities.NewProductEntity(name, price)

	// トランザクション開始
	conn, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.ProductDto{}, nil, err
	}
	tran, err := conn.Begin()
	if err != nil {
		return data_transfer_objects.ProductDto{}, nil, errors.WithStack(err)
	}

	// リポジトリで保存
	savedEntity, err := p.productRepository.Save(entity, tran)
	if err != nil {
		return data_transfer_objects.ProductDto{}, nil, err
	}

	err = tran.Commit()
	if err != nil {
		return data_transfer_objects.ProductDto{}, nil, errors.WithStack(err)
	}

	// dtoに詰める
	dto := data_transfer_objects.ProductDto{
		Name:  savedEntity.Name(),
		Price: savedEntity.Price(),
		BaseDto: data_transfer_objects.BaseDto{
			Id:        savedEntity.Id(),
			CreatedAt: savedEntity.CreatedAt(),
			UpdatedAt: savedEntity.UpdatedAt(),
		},
	}
	// 返却
	return dto, nil, nil
}

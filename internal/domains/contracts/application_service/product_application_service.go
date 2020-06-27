package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/pkg/errors"
)

type ProductApplicationService struct {
	productRepository interfaces.IProductRepository
}

func (p *ProductApplicationService) Register(name string, price string) (data_transfer_objects.ProductDto, ValidationErrors, error) {
	// 入力値バリデーション
	validationErrors, err := p.registerValidation(name, price)
	if err != nil {
		return data_transfer_objects.ProductDto{}, nil, err
	}
	if len(validationErrors) > 0 {
		return data_transfer_objects.ProductDto{}, validationErrors, nil
	}

	// entityを作成
	priceDecimal, err := decimal.NewFromString(price)
	entity, err := entities.NewProductEntity(name, priceDecimal)
	if err != nil {
		return data_transfer_objects.ProductDto{}, nil, err
	}

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
	dto := data_transfer_objects.NewProductDtoFromEntity(savedEntity)

	// 返却
	return dto, nil, nil
}

func (p *ProductApplicationService) registerValidation(name string, price string) (ValidationErrors, error) {
	validationErrors := ValidationErrors{}

	// 商品名バリデーション
	// 文字数50文字以上
	productNameValidErrors := values.ProductNameValidate(name)
	if len(productNameValidErrors) > 0 {
		validationErrors["name"] = []string{}
	}
	for _, validError := range productNameValidErrors {
		validationErrors["name"] = append(validationErrors["name"], validError.Error())
	}
	// 重複チェック
	productEntity, err := p.productRepository.GetByName(name, nil)
	if err != nil {
		return validationErrors, err
	}
	if productEntity != nil {
		validationErrors["name"] = append(validationErrors["name"], "すでに存在します")
	}

	// 価格バリデーション
	// decimalに変換可能か？
	// マイナスでないか？

	//// 担当者名バリデーション
	//contactPersonNameValidErrors := values.ContactPersonNameValidate(contactPersonName)
	//if len(contactPersonNameValidErrors) > 0 {
	//	validationErrors["contact_person_name"] = []string{}
	//}
	//for _, validError := range contactPersonNameValidErrors {
	//	validationErrors["contact_person_name"] = append(validationErrors["contact_person_name"], validError.Error())
	//}
	//
	//// 社長名バリデーション
	//presidentNameValidErrors := values.PresidentNameValidate(presidentName)
	//if len(presidentNameValidErrors) > 0 {
	//	validationErrors["president_name"] = []string{}
	//}
	//for _, validError := range presidentNameValidErrors {
	//	validationErrors["president_name"] = append(validationErrors["president_name"], validError.Error())
	//}

	return validationErrors, nil

}

func (p *ProductApplicationService) Get(id int) (data_transfer_objects.ProductDto, error) {
	// リポジトリつかってデータ取得
	entity, err := p.productRepository.GetById(id, nil)
	if err != nil {
		return data_transfer_objects.ProductDto{}, err
	}
	if entity == nil {
		// データがない
		return data_transfer_objects.ProductDto{}, nil
	}

	// dtoにつめる
	dto := data_transfer_objects.NewProductDtoFromEntity(entity)

	// 返却
	return dto, nil
}

package application_service

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values/validators"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
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
	entity, err := entities.NewProductEntity(name, price)
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
	productNameValidErrors := values.ProductNameValidate(name)
	if len(productNameValidErrors) > 0 {
		validationErrors["name"] = []string{}
	}
	for _, validError := range productNameValidErrors {
		var errorMessage string
		switch validError {
		case validators.EmptyStringValidError:
			errorMessage = "空です"
		case validators.OverLengthStringValidError:
			errorMessage = fmt.Sprintf("%v文字より多いです", values.ProductNameMaxLength)
		default:
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。name: %v, validErrorText: %v", name, validators.ValidErrorTest(validError)))
		}
		validationErrors["name"] = append(validationErrors["name"], errorMessage)
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
	productPriceValidErrors, err := values.ProductPriceValidate(price)
	if err != nil {
		return validationErrors, err
	}
	if len(productPriceValidErrors) > 0 {
		validationErrors["price"] = []string{}
	}
	for _, validError := range productPriceValidErrors {
		//validationErrors["price"] = append(validationErrors["price"], validError.Error())
		var errorMessage string
		switch validError {
		case validators.EmptyStringValidError:
			errorMessage = "空です"
		case validators.NumericStringValidError:
			errorMessage = "数値ではありません"
		case validators.NegativeValidError:
			errorMessage = "マイナス値です"
		default:
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。price: %v, validErrorText: %v", price, validators.ValidErrorTest(validError)))
		}
		validationErrors["price"] = append(validationErrors["price"], errorMessage)
	}

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

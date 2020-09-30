package application_service

import (
	"fmt"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/entities/values"
	"github.com/mixmaru/my_contracts/domains/contracts/entities/values/validators"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/db_connection"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type ProductApplicationService struct {
	productRepository interfaces.IProductRepository
}

func (p *ProductApplicationService) Register(name string, price string) (productDto data_transfer_objects.ProductDto, validationErrors map[string][]string, err error) {
	// トランザクション開始
	conn, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.ProductDto{}, nil, err
	}
	defer conn.Db.Close()
	tran, err := conn.Begin()

	if err != nil {
		return data_transfer_objects.ProductDto{}, nil, errors.WithStack(err)
	}
	// 入力値バリデーション
	validationErrors, err = p.registerValidation(name, price, tran)
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

	// リポジトリで保存
	savedId, err := p.productRepository.Save(entity, tran)
	if err != nil {
		return data_transfer_objects.ProductDto{}, nil, err
	}

	// 再読込
	savedEntity, err := p.productRepository.GetById(savedId, tran)
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

func (p *ProductApplicationService) registerValidation(name string, price string, executor gorp.SqlExecutor) (validationErrors map[string][]string, err error) {
	validationErrors = map[string][]string{}

	// 商品名バリデーション
	productNameValidErrors, err := values.ProductNameValue{}.Validate(name)
	if err != nil {
		return nil, err
	}

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
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。name: %v, validErrorText: %v", name, validators.ValidErrorText(validError)))
		}
		validationErrors["name"] = append(validationErrors["name"], errorMessage)
	}

	// 価格バリデーション
	productPriceValidErrors, err := values.ProductPriceValue{}.Validate(price)
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
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。price: %v, validErrorText: %v", price, validators.ValidErrorText(validError)))
		}
		validationErrors["price"] = append(validationErrors["price"], errorMessage)
	}

	return validationErrors, nil

}

func (p *ProductApplicationService) Get(id int) (data_transfer_objects.ProductDto, error) {
	conn, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.ProductDto{}, err
	}
	defer conn.Db.Close()

	// リポジトリつかってデータ取得
	entity, err := p.productRepository.GetById(id, conn)
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
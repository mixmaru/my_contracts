package create

import (
	"fmt"
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/mixmaru/my_contracts/core/domain/validators"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/pkg/errors"
)

type ProductCreateInteractor struct {
	productRepository products.IProductRepository
}

func NewProductCreateInteractor(productRepository products.IProductRepository) *ProductCreateInteractor {
	return &ProductCreateInteractor{
		productRepository: productRepository,
	}
}

func (i *ProductCreateInteractor) Handle(request *ProductCreateUseCaseRequest) (*ProductCreateUseCaseResponse, error) {
	// トランザクション開始
	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()
	tran, err := conn.Begin()

	if err != nil {
		return nil, errors.WithStack(err)
	}
	// 入力値バリデーション
	validationErrors, err := createValidation(request.Name, request.Price)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		response := NewProductCreateUseCaseResponse(products.ProductDto{}, validationErrors)
		return response, nil
	}

	// entityを作成
	entity, err := product.NewProductEntity(request.Name, request.Price)
	if err != nil {
		return nil, err
	}

	// リポジトリで保存
	savedId, err := i.productRepository.Save(entity, tran)
	if err != nil {
		return nil, err
	}

	// 再読込
	savedEntity, err := i.productRepository.GetById(savedId, tran)
	if err != nil {
		return nil, err
	}

	err = tran.Commit()
	if err != nil {
		return nil, errors.Wrapf(err, "コミット失敗。savedEntity: %+v", savedEntity)
	}

	// dtoに詰める
	dto := products.NewProductDtoFromEntity(savedEntity)

	response := NewProductCreateUseCaseResponse(dto, validationErrors)
	return response, nil
}

func createValidation(name string, price string) (validationErrors map[string][]string, err error) {
	validationErrors = map[string][]string{}

	// 商品名バリデーション
	productNameValidErrors, err := product.ProductNameValue{}.Validate(name)
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
			errorMessage = fmt.Sprintf("%v文字より多いです", product.ProductNameMaxLength)
		default:
			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。name: %v, validErrorText: %v", name, validators.ValidErrorText(validError)))
		}
		validationErrors["name"] = append(validationErrors["name"], errorMessage)
	}

	// 価格バリデーション
	productPriceValidErrors, err := product.ProductPriceValue{}.Validate(price)
	if err != nil {
		return validationErrors, err
	}
	if len(productPriceValidErrors) > 0 {
		validationErrors["price"] = []string{}
	}
	for _, validError := range productPriceValidErrors {
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

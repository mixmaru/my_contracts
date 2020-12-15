package create

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type CustomerPropertyTypeCreateInteractor struct {
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository
}

func NewCustomerPropertyTypeCreateInteractor(
	customerPropertyTypeRepository customer_property_type.ICustomerPropertyTypeRepository,
) *CustomerPropertyTypeCreateInteractor {
	return &CustomerPropertyTypeCreateInteractor{
		customerPropertyTypeRepository: customerPropertyTypeRepository,
	}
}

const MAX_RETRY_NUM int = 2

func (i *CustomerPropertyTypeCreateInteractor) Handle(
	request *CustomerPropertyTypeCreateUseCaseRequest,
) (*CustomerPropertyTypeCreateUseCaseResponse, error) {
	response := CustomerPropertyTypeCreateUseCaseResponse{}

	var err error
	// db接続
	conn, err := db.GetConnection()
	if err != nil {
		return nil, errors.Wrapf(err, "db接続に失敗しました")
	}
	defer conn.Db.Close()

	// 事前にnameの重複チェックを行っても別トランザクションから先に挿入されると重複する可能性がある
	// nameの重複が起こったら1回まで再実行する
	execCount := 0
	var savedIds []int
	for {
		execCount++
		if execCount > MAX_RETRY_NUM {
			return nil, errors.Errorf("実行回数がMAX_RETRY_NUMを超えました。execCount: %v, MAX_RETRY_NUM: %v, request: %+v", execCount, MAX_RETRY_NUM, request)
		}

		tran, err := conn.Begin()
		if err != nil {
			return nil, errors.Wrapf(err, "トランザクション開始に失敗しました")
		}

		// バリデーション
		response.ValidationError, err = i.validation(request, tran)
		if err != nil {
			return nil, err
		}

		if len(response.ValidationError) > 0 {
			return &response, nil
		}

		// エンティティを作る
		var propertyType customer.PropertyType
		switch request.Type {
		case "string":
			propertyType = customer.PROPERTY_TYPE_STRING
		case "numeric":
			propertyType = customer.PROPERTY_TYPE_NUMERIC
		default:
			return nil, errors.Errorf("想定外")
		}
		entity := customer.NewCustomerPropertyTypeEntity(request.Name, propertyType)

		// 登録実行する
		savedIds, err = i.customerPropertyTypeRepository.Create([]*customer.CustomerPropertyTypeEntity{entity}, tran)
		if err != nil {
			if execCount < MAX_RETRY_NUM {
				continue
			}
			tran.Rollback()
			return nil, errors.Wrapf(err, "保存実行に失敗しました。entity: %v", entity)
		}
		if err := tran.Commit(); err != nil {
			tran.Rollback()
			return nil, errors.Wrapf(err, "コミットに失敗しました")
		}
		break
	}

	// 再読込する
	reloadedEntities, err := i.customerPropertyTypeRepository.GetByIds(savedIds, conn)
	if err != nil {
		return nil, errors.Wrapf(err, "再読込に失敗しました。savedIds: %v", savedIds)
	}

	// dtoに詰める
	response.CustomerPropertyTypeDto, err = customer_property_type.NewCustomerPropertyTypeDtoFromEntity(reloadedEntities[0])
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (i *CustomerPropertyTypeCreateInteractor) validation(request *CustomerPropertyTypeCreateUseCaseRequest, executor gorp.SqlExecutor) (map[string][]string, error) {
	validationErrors := map[string][]string{}

	// 同名チェック
	entity, err := i.customerPropertyTypeRepository.GetByName(request.Name, executor)
	if err != nil {
		return nil, err
	}
	if entity != nil {
		validationErrors["name"] = []string{"既に存在する名前です"}
	}

	// タイプチェック
	if request.Type != "string" && request.Type != "numeric" {
		validationErrors["type"] = []string{"stringでもnumericでもありません"}
	}
	return validationErrors, nil
}

//
//func (i *ProductCreateInteractor) Handle(request *ProductCreateUseCaseRequest) (*ProductCreateUseCaseResponse, error) {
//	// トランザクション開始
//	conn, err := db.GetConnection()
//	if err != nil {
//		return nil, err
//	}
//	defer conn.Db.Close()
//	tran, err := conn.Begin()
//
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	// 入力値バリデーション
//	validationErrors, err := createValidation(request.Name, request.Price)
//	if err != nil {
//		return nil, err
//	}
//	if len(validationErrors) > 0 {
//		response := NewProductCreateUseCaseResponse(products.ProductDto{}, validationErrors)
//		return response, nil
//	}
//
//	// entityを作成
//	entity, err := product.NewProductEntity(request.Name, request.Price)
//	if err != nil {
//		return nil, err
//	}
//
//	// リポジトリで保存
//	savedId, err := i.productRepository.Save(entity, tran)
//	if err != nil {
//		return nil, err
//	}
//
//	// 再読込
//	savedEntity, err := i.productRepository.GetById(savedId, tran)
//	if err != nil {
//		return nil, err
//	}
//
//	err = tran.Commit()
//	if err != nil {
//		return nil, errors.Wrapf(err, "コミット失敗。savedEntity: %+v", savedEntity)
//	}
//
//	// dtoに詰める
//	dto := products.NewProductDtoFromEntity(savedEntity)
//
//	response := NewProductCreateUseCaseResponse(dto, validationErrors)
//	return response, nil
//}
//
//func createValidation(name string, price string) (validationErrors map[string][]string, err error) {
//	validationErrors = map[string][]string{}
//
//	// 商品名バリデーション
//	productNameValidErrors, err := product.ProductNameValue{}.Validate(name)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(productNameValidErrors) > 0 {
//		validationErrors["name"] = []string{}
//	}
//	for _, validError := range productNameValidErrors {
//		var errorMessage string
//		switch validError {
//		case validators.EmptyStringValidError:
//			errorMessage = "空です"
//		case validators.OverLengthStringValidError:
//			errorMessage = fmt.Sprintf("%v文字より多いです", product.ProductNameMaxLength)
//		default:
//			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。name: %v, validErrorText: %v", name, validators.ValidErrorText(validError)))
//		}
//		validationErrors["name"] = append(validationErrors["name"], errorMessage)
//	}
//
//	// 価格バリデーション
//	productPriceValidErrors, err := product.ProductPriceValue{}.Validate(price)
//	if err != nil {
//		return validationErrors, err
//	}
//	if len(productPriceValidErrors) > 0 {
//		validationErrors["price"] = []string{}
//	}
//	for _, validError := range productPriceValidErrors {
//		var errorMessage string
//		switch validError {
//		case validators.EmptyStringValidError:
//			errorMessage = "空です"
//		case validators.NumericStringValidError:
//			errorMessage = "数値ではありません"
//		case validators.NegativeValidError:
//			errorMessage = "マイナス値です"
//		default:
//			return validationErrors, errors.New(fmt.Sprintf("想定外エラー。price: %v, validErrorText: %v", price, validators.ValidErrorText(validError)))
//		}
//		validationErrors["price"] = append(validationErrors["price"], errorMessage)
//	}
//
//	return validationErrors, nil
//}

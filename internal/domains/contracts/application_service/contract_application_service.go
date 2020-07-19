package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type ContractApplicationService struct {
	ContractRepository interfaces.IContractRepository
	UserRepository     interfaces.IUserRepository
}

func (c *ContractApplicationService) Register(userId int, productId int) (productDto data_transfer_objects.ContractDto, validationErrors map[string][]string, err error) {
	// トランザクション開始
	conn, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, err
	}
	defer conn.Db.Close()
	tran, err := conn.Begin()
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, errors.WithStack(err)
	}

	// 入力値バリデーション
	validationErrors, err = c.registerValidation(userId, productId, tran)
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, err
	}
	if len(validationErrors) > 0 {
		return data_transfer_objects.ContractDto{}, validationErrors, nil
	}

	// entityを作成
	entity := entities.NewContractEntity(userId, productId)

	// リポジトリで保存
	savedId, err := c.ContractRepository.Create(entity, tran)
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, err
	}

	// 再読込
	savedEntity, _, _, err := c.ContractRepository.GetById(savedId, tran)
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, err
	}

	err = tran.Commit()
	if err != nil {
		return data_transfer_objects.ContractDto{}, nil, errors.WithStack(err)
	}

	// dtoに詰める
	dto := data_transfer_objects.NewContractDtoFromEntity(savedEntity)

	// 返却
	return dto, nil, nil
}

func (c *ContractApplicationService) registerValidation(userId int, productId int, executor gorp.SqlExecutor) (validationErrors map[string][]string, err error) {
	validationErrors = map[string][]string{}

	// userの存在チェック
	user, err := c.UserRepository.GetUserById(userId, executor)
	if err != nil {
		return nil, err
	}
	if user == nil {
		validationErrors["user_id"] = []string{
			"存在しません",
		}
	}

	return validationErrors, nil

	// 商品名バリデーション
	//productNameValidErrors, err := values.ProductNameValue{}.Validate(name)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if len(productNameValidErrors) > 0 {
	//	validationErrors["name"] = []string{}
	//}
	//for _, validError := range productNameValidErrors {
	//	var errorMessage string
	//	switch validError {
	//	case validators.EmptyStringValidError:
	//		errorMessage = "空です"
	//	case validators.OverLengthStringValidError:
	//		errorMessage = fmt.Sprintf("%v文字より多いです", values.ProductNameMaxLength)
	//	default:
	//		return validationErrors, errors.New(fmt.Sprintf("想定外エラー。name: %v, validErrorText: %v", name, validators.ValidErrorTest(validError)))
	//	}
	//	validationErrors["name"] = append(validationErrors["name"], errorMessage)
	//}

	// 重複チェック
	//productEntity, err := p.productRepository.GetByName(name, executor)
	//if err != nil {
	//	return validationErrors, err
	//}
	//if productEntity != nil {
	//	validationErrors["name"] = append(validationErrors["name"], "すでに存在します")
	//}
	//
	//// 価格バリデーション
	//productPriceValidErrors, err := values.ProductPriceValue{}.Validate(price)
	//if err != nil {
	//	return validationErrors, err
	//}
	//if len(productPriceValidErrors) > 0 {
	//	validationErrors["price"] = []string{}
	//}
	//for _, validError := range productPriceValidErrors {
	//	//validationErrors["price"] = append(validationErrors["price"], validError.Error())
	//	var errorMessage string
	//	switch validError {
	//	case validators.EmptyStringValidError:
	//		errorMessage = "空です"
	//	case validators.NumericStringValidError:
	//		errorMessage = "数値ではありません"
	//	case validators.NegativeValidError:
	//		errorMessage = "マイナス値です"
	//	default:
	//		return validationErrors, errors.New(fmt.Sprintf("想定外エラー。price: %v, validErrorText: %v", price, validators.ValidErrorTest(validError)))
	//	}
	//	validationErrors["price"] = append(validationErrors["price"], errorMessage)
	//}
	//
	//return validationErrors, nil
	//
}

func (c *ContractApplicationService) GetById(id int) (contractDto data_transfer_objects.ContractDto, productDto data_transfer_objects.ProductDto, userDto interface{}, err error) {
	conn, err := db_connection.GetConnection()
	if err != nil {
		return data_transfer_objects.ContractDto{}, data_transfer_objects.ProductDto{}, nil, err
	}
	defer conn.Db.Close()

	// リポジトリつかってデータ取得
	contractEntity, productEntity, userEntity, err := c.ContractRepository.GetById(id, conn)
	if err != nil {
		return data_transfer_objects.ContractDto{}, data_transfer_objects.ProductDto{}, nil, err
	}
	if contractEntity == nil {
		// データがない
		return data_transfer_objects.ContractDto{}, data_transfer_objects.ProductDto{}, nil, nil
	}

	// dtoにつめる
	contractDto = data_transfer_objects.NewContractDtoFromEntity(contractEntity)
	productDto = data_transfer_objects.NewProductDtoFromEntity(productEntity)
	switch userEntity.(type) {
	case *entities.UserIndividualEntity:
		userDto = data_transfer_objects.NewUserIndividualDtoFromEntity(userEntity.(*entities.UserIndividualEntity))
	case *entities.UserCorporationEntity:
		userDto = data_transfer_objects.NewUserCorporationDtoFromEntity(userEntity.(*entities.UserCorporationEntity))
	default:
		return data_transfer_objects.ContractDto{}, data_transfer_objects.ProductDto{}, nil, errors.Errorf("意図しないUser型が来た。userEntity: %t", userEntity)
	}

	// 返却
	return contractDto, productDto, userDto, nil
}

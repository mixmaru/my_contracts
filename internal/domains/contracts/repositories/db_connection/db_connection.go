package db_connection

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"os"
)

// gorpのdbMapを作成する
func GetConnection() (*gorp.DbMap, error) {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	executeMode, err := utils.GetExecuteMode()
	if err != nil {
		return nil, err
	}
	connectionStr, err := getConnectionString(executeMode)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(data_mappers.UserMapper{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(data_mappers.UserIndividualMapper{}, "users_individual")
	dbmap.AddTableWithName(data_mappers.UserCorporationMapper{}, "users_corporation")

	dbmap.AddTableWithName(data_mappers.ProductMapper{}, "products").SetKeys(true, "Id")
	dbmap.AddTableWithName(data_mappers.ProductPriceMonthlyMapper{}, "product_price_monthlies")

	dbmap.AddTableWithName(data_mappers.ContractMapper{}, "contracts").SetKeys(true, "Id")

	return dbmap, nil
}

// 実行モード（test, development, production）を渡すと、適したdb接続文字列を返す
func getConnectionString(executeMode string) (string, error) {
	switch executeMode {
	case utils.Test:
		return fmt.Sprintf("host=%v user=%v dbname=%v password=%v sslmode=disable", os.Getenv("DB_TEST_HOST"), os.Getenv("DB_TEST_USER"), os.Getenv("DB_TEST_NAME"), os.Getenv("DB_TEST_PASSWORD")), nil
	case utils.Development:
		return fmt.Sprintf("host=%v user=%v dbname=%v password=%v sslmode=disable", os.Getenv("DB_DEVELOPMENT_HOST"), os.Getenv("DB_DEVELOPMENT_USER"), os.Getenv("DB_DEVELOPMENT_NAME"), os.Getenv("DB_DEVELOPMENT_PASSWORD")), nil
	case utils.Production:
		return fmt.Sprintf("host=%v user=%v dbname=%v password=%v sslmode=disable", os.Getenv("DB_PRODUCTION_HOST"), os.Getenv("DB_PRODUCTION_USER"), os.Getenv("DB_PRODUCTION_NAME"), os.Getenv("DB_PRODUCTION_PASSWORD")), nil
	default:
		return "", errors.New(fmt.Sprintf("考慮されてない値が渡されました。executeMode: %+v", executeMode))
	}
}

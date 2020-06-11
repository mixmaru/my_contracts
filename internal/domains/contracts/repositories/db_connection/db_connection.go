package db_connection

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	productTables "github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/product/tables"
	userTables "github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user/tables"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"os"
	"reflect"
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
	dbmap.AddTableWithName(userTables.UserRecord{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(userTables.UserIndividualRecord{}, "users_individual")
	dbmap.AddTableWithName(userTables.UserCorporationRecord{}, "users_corporation")

	dbmap.AddTableWithName(productTables.ProductRecord{}, "products").SetKeys(true, "Id")

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

// executorがトランザクションだったらそれをそのまま返す。
// executorがnilだったら、dbConnectionを返す。
// ※repository内で、いちいちそれがトランザクションなのか、dbConnectionを取得しないと行けないのかの条件分岐を書く必要性を無くすために用意した
func GetConnectionIfNotTransaction(executor gorp.SqlExecutor) (gorp.SqlExecutor, error) {
	if executor == nil || reflect.ValueOf(executor).IsNil() {
		return GetConnection()
	}

	tran, ok := executor.(*gorp.Transaction)
	if ok {
		return tran, nil
	}

	return nil, errors.New(fmt.Sprintf("GetConnectionに失敗しました。executorが考慮外 executor: %+v", executor))
}

// executorがDbMapだったらCloseする
// executorがトランザクションだったらなにもしない
// ※repository内で、いちいちそれがトランザクションなのか、DbMapなのかを判断してClose処理のための条件分岐を書く必要性を無くすために用意した
func CloseConnectionIfNotTransaction(executor gorp.SqlExecutor) error {
	// dbMapが渡されたらそれをcloseする
	dbMap, ok := executor.(*gorp.DbMap)
	if ok {
		dbMap.Db.Close()
		return nil
	}

	// transactionが渡されたら、なにもしない。（トランザクション管理は上位でおこなっているため）
	_, ok = executor.(*gorp.Transaction)
	if ok {
		return nil
	}

	return errors.New(fmt.Sprintf("CloseConnectionに失敗しました。executorが考慮外 executor: %+v", executor))
}

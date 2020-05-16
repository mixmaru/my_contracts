package db_connection

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user/tables"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"reflect"
)

// gorpのdbMapを作成する
func GetConnection() (*gorp.DbMap, error) {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("postgres", "user=postgres dbname=my_contracts_development password=password sslmode=disable")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(tables.UserRecord{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(tables.UserIndividualRecord{}, "users_individual")
	dbmap.AddTableWithName(tables.UserCorporationRecord{}, "users_corporation")

	return dbmap, nil
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

package db_connection

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user/tables"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

func GetConnection(executor gorp.SqlExecutor) (gorp.SqlExecutor, error) {
	if executor == nil {
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

	tran, ok := executor.(*gorp.Transaction)
	if ok {
		return tran, nil
	}

	return nil, errors.New(fmt.Sprintf("GetConnectionに失敗しました。executorが考慮外 executor: %+v", executor))
}

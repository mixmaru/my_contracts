package user

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user/structures"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"time"
)

type Repository struct {
}

func InitDb() (*gorp.DbMap, error) {
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
	dbmap.AddTableWithName(structures.User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(structures.UserIndividual{}, "users_individual")

	return dbmap, nil
}

func (r *Repository) Save(individual *user.UserIndividual, sqlExecutor gorp.SqlExecutor) error {
	// エンティティからリポジトリ用構造体に値をセットし直す
	// もしくはエンティティが吐き出すようにしてもいいかも。あとで考える
	now := time.Now()

	user := &structures.User{}
	user.CreatedAt = now
	user.UpdatedAt = now

	err := sqlExecutor.Insert(user)
	if err != nil {
		return errors.WithStack(err)
	}

	// individualを保存
	u_individual := &structures.UserIndividual{}
	u_individual.UserId = user.Id
	u_individual.Name = individual.Name()
	u_individual.UpdatedAt = now
	u_individual.CreatedAt = now

	err = sqlExecutor.Insert(u_individual)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

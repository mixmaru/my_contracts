package user

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user_individual"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_InitDb(t *testing.T) {
	_, err := InitDb()
	assert.NoError(t, err)
}

func TestUser_Save(t *testing.T) {
	// 事前準備
	db, err := InitDb()
	assert.NoError(t, err)
	repo := Repository{}
	user := user_individual.UserIndividual{}
	user.SetName("個人太郎")

	// 実行
	err = repo.Save(user, db)

	// test
	// とりあえずエラーでなければokとする
	assert.NoError(t, err)
}

//func TestUser_test(t *testing.T) {
//	db := initDb()
//	user := &Repository{}
//	user.CreatedAt = time.Now()
//	user.UpdatedAt = time.Now()
//	err := db.Insert(user)
//	if err != nil {
//		assert.Failf(t, "%+v", err.Error())
//	}
//
//	loadUser := &Repository{}
//	err = db.SelectOne(loadUser, "SELECT * FROM users WHERE id = $1", user.Id)
//	if err != nil {
//		assert.Failf(t, "%+v", err.Error())
//	}
//
//}
//
//func initDb() *gorp.DbMap {
//	// connect to db using standard Go database/sql API
//	// use whatever database/sql driver you wish
//	db, err := sql.Open("postgres", "user=postgres dbname=my_contracts_development password=password sslmode=disable")
//	checkErr(err, "sql.Open failed")
//
//	// construct a gorp DbMap
//	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
//
//	// add a table, setting the table name to 'posts' and
//	// specifying that the Id property is an auto incrementing PK
//	dbmap.AddTableWithName(Repository{}, "users").SetKeys(true, "Id")
//
//	// create the table. in a production system you'd generally
//	// use a migration tool, or create the tables via scripts
//	//err = dbmap.CreateTablesIfNotExists()
//	//checkErr(err, "Create tables failed")
//
//	return dbmap
//}
//
//func checkErr(err error, msg string) {
//	if err != nil {
//		log.Fatalln(msg, err)
//	}
//}

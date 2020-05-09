package user

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestUser_test(t *testing.T) {
	db := initDb()
	user := &User{}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err := db.Insert(user)
	if err != nil {
		assert.Failf(t, "%+v", err.Error())
	}
}

//var err error
//Db, err = sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
//if err != nil {
//panic(err)
//}
type User struct {
	Id        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("postgres", "user=postgres dbname=my_contracts_development password=password sslmode=disable")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	//err = dbmap.CreateTablesIfNotExists()
	//checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

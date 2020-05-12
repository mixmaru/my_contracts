package user

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user/tables"
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
	dbmap.AddTableWithName(tables.UserRecord{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(tables.UserIndividualRecord{}, "users_individual")
	dbmap.AddTableWithName(tables.UserCorporationRecord{}, "users_corporation")

	return dbmap, nil
}

// 個人顧客エンティティを保存する
func (r *Repository) SaveUserIndividual(userEntity *user.UserIndividualEntity, executor gorp.SqlExecutor) error {
	// エンティティからリポジトリ用構造体に値をセットし直す
	// もしくはエンティティが吐き出すようにしてもいいかも。あとで考える
	// db用構造体オブジェクトがentityを読み込む用にする。
	now := time.Now()

	user := tables.NewUserRecordFromUserIndividualEntity(userEntity)
	user.CreatedAt = now
	user.UpdatedAt = now

	err := executor.Insert(user)
	if err != nil {
		return errors.WithStack(err)
	}

	// individualを保存
	userIndividualDbMap := tables.NewUserIndividualRecordFromUserIndividualEntity(userEntity)
	userIndividualDbMap.UserId = user.Id
	userIndividualDbMap.CreatedAt = now
	userIndividualDbMap.UpdatedAt = now

	err = executor.Insert(userIndividualDbMap)
	if err != nil {
		return errors.WithStack(err)
	}

	// dbから再読込してentityに詰め直す
	userDbData, err := r.getUserIndividualViewById(user.Id, executor)
	if err != nil {
		return err
	}
	userEntity.LoadData(
		userDbData.Id,
		userDbData.Name,
		userDbData.CreatedAt,
		userDbData.UpdatedAt,
	)

	return nil
}

func (r *Repository) GetUserIndividualById(id int, sqlExecutor gorp.SqlExecutor) (*user.UserIndividualEntity, error) {
	// dbからデータ取得
	userData, err := r.getUserIndividualViewById(id, sqlExecutor)
	if err != nil {
		return nil, err
	}

	// entityに詰める
	userEntity := user.NewUserIndividualEntityWithData(
		userData.Id,
		userData.Name,
		userData.CreatedAt,
		userData.UpdatedAt,
	)
	return userEntity, nil
}

// dbからid指定で個人顧客情報を取得する
func (r *Repository) getUserIndividualViewById(id int, executor gorp.SqlExecutor) (*tables.UserIndividualView, error) {
	data := &tables.UserIndividualView{}
	err := executor.SelectOne(
		data,
		"SELECT u.id, ui.name, u.created_at, u.updated_at FROM users u "+
			"inner join users_individual ui on u.id = ui.user_id "+
			"WHERE id = $1",
		id,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

// 法人顧客エンティティを保存する
func (r *Repository) SaveUserCorporation(userEntity *user.UserCorporationEntity, executor gorp.SqlExecutor) error {
	now := time.Now()

	// userRecord作成
	userRecord := tables.NewUserRecordFromUserCorporationEntity(userEntity)
	userRecord.CreatedAt = now
	userRecord.UpdatedAt = now

	// 保存
	err := executor.Insert(userRecord)
	if err != nil {
		return errors.WithStack(err)
	}

	// userCorporationRecord作成
	userCorporationRecord := tables.NewUserCorporationRecordFromUserCorporationEntity(userEntity)
	userCorporationRecord.UserId = userRecord.Id
	userCorporationRecord.CreatedAt = now
	userCorporationRecord.UpdatedAt = now

	// 保存
	err = executor.Insert(userCorporationRecord)
	if err != nil {
		return errors.WithStack(err)
	}

	// todo 再読込する
	return nil
}

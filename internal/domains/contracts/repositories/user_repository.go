package repositories

import (
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type UserRepository struct {
}

// 個人顧客エンティティを保存する
func (r *UserRepository) SaveUserIndividual(userEntity *entities.UserIndividualEntity, transaction *gorp.Transaction) (*entities.UserIndividualEntity, error) {
	// db接続。
	conn, err := db_connection.GetConnectionIfNotTransaction(transaction)
	if err != nil {
		return nil, err
	}
	defer db_connection.CloseConnectionIfNotTransaction(conn)

	// エンティティからリポジトリ用構造体に値をセットし直す
	user := data_mappers.NewUserRecordFromUserIndividualEntity(userEntity)

	err = conn.Insert(user)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// individualを保存
	userIndividualDbMap := data_mappers.NewUserIndividualRecordFromUserIndividualEntity(userEntity)
	userIndividualDbMap.UserId = user.Id

	err = conn.Insert(userIndividualDbMap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// dbから再読込してentityに詰め直す
	return r.getUserIndividualEntityById(user.Id, userEntity, conn)
}

// Idで個人顧客情報を取得する。データがなければnilを返す
func (r *UserRepository) GetUserIndividualById(id int, transaction *gorp.Transaction) (*entities.UserIndividualEntity, error) {
	// db接続。
	conn, err := db_connection.GetConnectionIfNotTransaction(transaction)
	if err != nil {
		return nil, err
	}
	defer db_connection.CloseConnectionIfNotTransaction(conn)

	// dbからデータ取得
	return r.getUserIndividualEntityById(id, &entities.UserIndividualEntity{}, conn)
}

// dbからid指定で個人顧客情報を取得する
func (r *UserRepository) getUserIndividualEntityById(id int, entity *entities.UserIndividualEntity, executor gorp.SqlExecutor) (*entities.UserIndividualEntity, error) {
	userIndividualView := data_mappers.UserIndividualView{}
	noRow, err := selectOne(
		executor,
		&userIndividualView,
		entity,
		"SELECT u.id, ui.name, u.created_at, u.updated_at FROM users u "+
			"inner join users_individual ui on u.id = ui.user_id "+
			"WHERE id = $1", id,
	)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return entity, nil
}

// 法人顧客エンティティを保存する
func (r *UserRepository) SaveUserCorporation(userEntity *entities.UserCorporationEntity, transaction *gorp.Transaction) (*entities.UserCorporationEntity, error) {
	// db接続。
	conn, err := db_connection.GetConnectionIfNotTransaction(transaction)
	if err != nil {
		return nil, err
	}
	defer db_connection.CloseConnectionIfNotTransaction(conn)

	// userRecord作成
	userRecord := data_mappers.NewUserRecordFromUserCorporationEntity(userEntity)

	// 保存
	err = conn.Insert(userRecord)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// userCorporationRecord作成
	userCorporationRecord := data_mappers.NewUserCorporationRecordFromUserCorporationEntity(userEntity)
	userCorporationRecord.UserId = userRecord.Id

	// 保存
	err = conn.Insert(userCorporationRecord)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 再読込する
	return r.getUserCorporationEntityById(userRecord.Id, userEntity, conn)
}

// dbからid指定で法人顧客情報を取得する
func (r *UserRepository) getUserCorporationEntityById(id int, entity *entities.UserCorporationEntity, executor gorp.SqlExecutor) (*entities.UserCorporationEntity, error) {
	// dbからデータ取得
	record := data_mappers.UserCorporationView{}
	query := "SELECT u.id, uc.contact_person_name, uc.president_name, u.created_at, u.updated_at " +
		"FROM users u " +
		"inner join users_corporation uc on u.id = uc.user_id " +
		"WHERE id = $1"
	noRow, err := selectOne(executor, &record, entity, query, id)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return entity, nil
}

// Idで法人顧客情報を取得する。データがなければnilを返す
func (r *UserRepository) GetUserCorporationById(id int, transaction *gorp.Transaction) (*entities.UserCorporationEntity, error) {
	// db接続。
	conn, err := db_connection.GetConnectionIfNotTransaction(transaction)
	if err != nil {
		return nil, err
	}
	defer db_connection.CloseConnectionIfNotTransaction(conn)

	return r.getUserCorporationEntityById(id, &entities.UserCorporationEntity{}, conn)
}

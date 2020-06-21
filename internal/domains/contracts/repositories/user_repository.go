package repositories

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/tables"
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
	user := tables.NewUserRecordFromUserIndividualEntity(userEntity)

	err = conn.Insert(user)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// individualを保存
	userIndividualDbMap := tables.NewUserIndividualRecordFromUserIndividualEntity(userEntity)
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
	userIndividualView := tables.UserIndividualView{}
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
	userRecord := tables.NewUserRecordFromUserCorporationEntity(userEntity)

	// 保存
	err = conn.Insert(userRecord)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// userCorporationRecord作成
	userCorporationRecord := tables.NewUserCorporationRecordFromUserCorporationEntity(userEntity)
	userCorporationRecord.UserId = userRecord.Id

	// 保存
	err = conn.Insert(userCorporationRecord)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 再読込する
	data, err := r.getUserCorporationViewById(userRecord.Id, conn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userEntity.LoadData(data.Id, data.ContactPersonName, data.PresidentName, data.CreatedAt, data.UpdatedAt)
	return userEntity, nil
}

// dbからid指定で法人顧客情報を取得する
func (r *UserRepository) getUserCorporationEntityById(id int, executor gorp.SqlExecutor) (*entities.UserCorporationEntity, error) {
	data, err := r.getUserCorporationViewById(id, executor)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	entity, err := entities.NewUserCorporationEntityWithData(
		data.Id,
		data.ContactPersonName,
		data.PresidentName,
		data.CreatedAt,
		data.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// dbからid指定で法人顧客情報を取得する
func (r *UserRepository) getUserCorporationViewById(id int, executor gorp.SqlExecutor) (*tables.UserCorporationView, error) {
	data := &tables.UserCorporationView{}
	selectSql := "SELECT u.id, uc.contact_person_name, uc.president_name, u.created_at, u.updated_at " +
		"FROM users u " +
		"inner join users_corporation uc on u.id = uc.user_id " +
		"WHERE id = $1"
	err := executor.SelectOne(data, selectSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		} else {
			return nil, errors.WithStack(err)
		}
	}
	return data, nil
}

// Idで法人顧客情報を取得する。データがなければnilを返す
func (r *UserRepository) GetUserCorporationById(id int, transaction *gorp.Transaction) (*entities.UserCorporationEntity, error) {
	// db接続。
	conn, err := db_connection.GetConnectionIfNotTransaction(transaction)
	if err != nil {
		return nil, err
	}
	defer db_connection.CloseConnectionIfNotTransaction(conn)

	// dbからデータ取得
	record := tables.UserCorporationView{}
	entity := entities.UserCorporationEntity{}
	query := "SELECT u.id, uc.contact_person_name, uc.president_name, u.created_at, u.updated_at " +
		"FROM users u " +
		"inner join users_corporation uc on u.id = uc.user_id " +
		"WHERE id = $1"
	noRow, err := selectOne(conn, &record, &entity, query, id)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return &entity, nil
}

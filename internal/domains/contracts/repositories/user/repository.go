package user

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user/tables"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type Repository struct {
}

// 個人顧客エンティティを保存する
func (r *Repository) SaveUserIndividual(userEntity *user.UserIndividualEntity, transaction *gorp.Transaction) (*user.UserIndividualEntity, error) {
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
	userDbData, err := r.getUserIndividualViewById(user.Id, conn)
	if err != nil {
		return nil, err
	}
	userEntity.LoadData(
		userDbData.Id,
		userDbData.Name,
		userDbData.CreatedAt,
		userDbData.UpdatedAt,
	)

	return userEntity, nil
}

// Idで個人顧客情報を取得する。データがなければnilを返す
func (r *Repository) GetUserIndividualById(id int, transaction *gorp.Transaction) (*user.UserIndividualEntity, error) {
	// db接続。
	conn, err := db_connection.GetConnectionIfNotTransaction(transaction)
	if err != nil {
		return nil, err
	}
	defer db_connection.CloseConnectionIfNotTransaction(conn)

	// dbからデータ取得
	userData, err := r.getUserIndividualViewById(id, conn)
	if err == sql.ErrNoRows {
		// データがなかったら、nilを返す
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// entityに詰める
	userEntity, err := user.NewUserIndividualEntityWithData(
		userData.Id,
		userData.Name,
		userData.CreatedAt,
		userData.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
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
		if err == sql.ErrNoRows {
			return nil, err
		} else {
			return nil, errors.WithStack(err)
		}
	}

	return data, nil
}

// 法人顧客エンティティを保存する
func (r *Repository) SaveUserCorporation(userEntity *user.UserCorporationEntity, transaction *gorp.Transaction) (*user.UserCorporationEntity, error) {
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
func (r *Repository) getUserCorporationEntityById(id int, executor gorp.SqlExecutor) (*user.UserCorporationEntity, error) {
	data, err := r.getUserCorporationViewById(id, executor)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	entity, err := user.NewUserCorporationEntityWithData(
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
func (r *Repository) getUserCorporationViewById(id int, executor gorp.SqlExecutor) (*tables.UserCorporationView, error) {
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
func (r *Repository) GetUserCorporationById(id int, transaction *gorp.Transaction) (*user.UserCorporationEntity, error) {
	// db接続。
	conn, err := db_connection.GetConnectionIfNotTransaction(transaction)
	if err != nil {
		return nil, err
	}
	defer db_connection.CloseConnectionIfNotTransaction(conn)

	// dbからデータ取得
	userData, err := r.getUserCorporationViewById(id, conn)
	if err != nil {
		if err == sql.ErrNoRows {
			// 対象データがない。nilを返す
			return nil, nil
		} else {
			return nil, err
		}
	}

	// entityに詰める
	userEntity, err := user.NewUserCorporationEntityWithData(
		userData.Id,
		userData.ContactPersonName,
		userData.PresidentName,
		userData.CreatedAt,
		userData.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return userEntity, nil
}

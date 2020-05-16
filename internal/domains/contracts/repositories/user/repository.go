package user

import (
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/user/tables"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"time"
)

type Repository struct {
}

// 個人顧客エンティティを保存する
func (r *Repository) SaveUserIndividual(userEntity *user.UserIndividualEntity, transaction *gorp.Transaction) error {
	// db接続。
	var conn gorp.SqlExecutor
	if transaction != nil {
		// トランザクションが渡されていればそれを使う
		conn = transaction
	} else {
		// なければ、db接続を取得する
		dbMap, err := db_connection.GetConnectionIfNotTransaction()
		if err != nil {
			return err
		}
		defer dbMap.Db.Close()
		conn = dbMap
	}

	// エンティティからリポジトリ用構造体に値をセットし直す
	now := time.Now()

	user := tables.NewUserRecordFromUserIndividualEntity(userEntity)
	user.CreatedAt = now
	user.UpdatedAt = now

	err := conn.Insert(user)
	if err != nil {
		return errors.WithStack(err)
	}

	// individualを保存
	userIndividualDbMap := tables.NewUserIndividualRecordFromUserIndividualEntity(userEntity)
	userIndividualDbMap.UserId = user.Id
	userIndividualDbMap.CreatedAt = now
	userIndividualDbMap.UpdatedAt = now

	err = conn.Insert(userIndividualDbMap)
	if err != nil {
		return errors.WithStack(err)
	}

	// dbから再読込してentityに詰め直す
	userDbData, err := r.getUserIndividualViewById(user.Id, conn)
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

func (r *Repository) GetUserIndividualById(id int, transaction *gorp.Transaction) (*user.UserIndividualEntity, error) {
	// db接続。
	var conn gorp.SqlExecutor
	if transaction != nil {
		// トランザクションが渡されていればそれを使う
		conn = transaction
	} else {
		// なければ、db接続を取得する
		dbMap, err := db_connection.GetConnectionIfNotTransaction()
		if err != nil {
			return nil, err
		}
		defer dbMap.Db.Close()
		conn = dbMap
	}

	// dbからデータ取得
	userData, err := r.getUserIndividualViewById(id, conn)
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
func (r *Repository) SaveUserCorporation(userEntity *user.UserCorporationEntity, transaction *gorp.Transaction) error {
	// db接続。
	var conn gorp.SqlExecutor
	if transaction != nil {
		// トランザクションが渡されていればそれを使う
		conn = transaction
	} else {
		// なければ、db接続を取得する
		dbMap, err := db_connection.GetConnectionIfNotTransaction()
		if err != nil {
			return err
		}
		defer dbMap.Db.Close()
		conn = dbMap
	}

	now := time.Now()

	// userRecord作成
	userRecord := tables.NewUserRecordFromUserCorporationEntity(userEntity)
	userRecord.CreatedAt = now
	userRecord.UpdatedAt = now

	// 保存
	err := conn.Insert(userRecord)
	if err != nil {
		return errors.WithStack(err)
	}

	// userCorporationRecord作成
	userCorporationRecord := tables.NewUserCorporationRecordFromUserCorporationEntity(userEntity)
	userCorporationRecord.UserId = userRecord.Id
	userCorporationRecord.CreatedAt = now
	userCorporationRecord.UpdatedAt = now

	// 保存
	err = conn.Insert(userCorporationRecord)
	if err != nil {
		return errors.WithStack(err)
	}

	// 再読込する
	data, err := r.getUserCorporationViewById(userRecord.Id, conn)
	if err != nil {
		return errors.WithStack(err)
	}

	userEntity.LoadData(data.Id, data.ContactPersonName, data.PresidentName, data.CreatedAt, data.UpdatedAt)
	return nil
}

// dbからid指定で法人顧客情報を取得する
func (r *Repository) getUserCorporationEntityById(id int, executor gorp.SqlExecutor) (*user.UserCorporationEntity, error) {
	data, err := r.getUserCorporationViewById(id, executor)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	entity := user.NewUserCorporationEntityWithData(
		data.Id,
		data.ContactPersonName,
		data.PresidentName,
		data.CreatedAt,
		data.UpdatedAt,
	)

	return entity, nil
}

// dbからid指定で法人顧客情報を取得する
func (r *Repository) getUserCorporationViewById(id int, executor gorp.SqlExecutor) (*tables.UserCorporationView, error) {
	data := &tables.UserCorporationView{}
	sql := "SELECT u.id, uc.contact_person_name, uc.president_name, u.created_at, u.updated_at " +
		"FROM users u " +
		"inner join users_corporation uc on u.id = uc.user_id " +
		"WHERE id = $1"
	err := executor.SelectOne(data, sql, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

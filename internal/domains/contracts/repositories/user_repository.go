package repositories

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type UserRepository struct {
	*BaseRepository
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		&BaseRepository{},
	}
}

// 個人顧客エンティティを保存する
func (r *UserRepository) SaveUserIndividual(userEntity *entities.UserIndividualEntity, executor gorp.SqlExecutor) (savedId int, err error) {
	// エンティティからリポジトリ用構造体に値をセットし直す
	user := data_mappers.NewUserMapperFromUserIndividualEntity(userEntity)

	err = executor.Insert(user)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	// individualを保存
	userIndividualDbMap := data_mappers.NewUserIndividualMapperFromUserIndividualEntity(userEntity)
	userIndividualDbMap.UserId = user.Id

	err = executor.Insert(userIndividualDbMap)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return user.Id, nil
}

// Idで顧客情報を取得する。データがなければnilを返す
func (r *UserRepository) GetUserById(id int, executor gorp.SqlExecutor) (interface{}, error) {
	var mapper data_mappers.UserView
	query := `
select
   u.id as user_id,
   case
	   when ui.user_id IS NOT NULL then 'individual'
	   when uc.user_id IS NOT NULL then 'corporation'
	end as user_type,
   ui.name as user_individual_name,
   ui.created_at as user_individual_created_at,
   ui.updated_at as user_individual_updated_at,
   uc.company_name as user_corporation_corporation_name,
   uc.contact_person_name as user_corporation_contact_person_name,
   uc.president_name as user_corporation_president_name,
   uc.created_at as user_corporation_created_at,
   uc.updated_at as user_corporation_updated_at
from users u
left outer join users_individual ui on u.id = ui.user_id
left outer join users_corporation uc on u.id = uc.user_id
where u.id = $1
`
	err := executor.SelectOne(&mapper, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			// データがない
			return nil, nil
		} else {
			return nil, errors.Wrapf(err, "userデータ取得時エラー。id: %v", id)
		}
	}

	switch mapper.UserType {
	case "individual":
		retUser, err := entities.NewUserIndividualEntityWithData(mapper.UserId, mapper.UserIndividualName.String, mapper.UserIndividualCreatedAt.Time, mapper.UserIndividualUpdatedAt.Time)
		if err != nil {
			return nil, err
		}
		return retUser, nil
	case "corporation":
		retUser, err := entities.NewUserCorporationEntityWithData(mapper.UserId, mapper.UserCorporationCorporationName.String, mapper.UserCorporationContractPersonName.String, mapper.UserCorporationPresidentName.String, mapper.UserCorporationCreatedAt.Time, mapper.UserCorporationUpdatedAt.Time)
		if err != nil {
			return nil, err
		}
		return retUser, nil
	default:
		return nil, errors.Errorf("想定外のUserTypeが来た。mapper: %v", mapper)
	}
}

// Idで個人顧客情報を取得する。データがなければnilを返す
func (r *UserRepository) GetUserIndividualById(id int, executor gorp.SqlExecutor) (*entities.UserIndividualEntity, error) {
	// dbからデータ取得
	return r.getUserIndividualEntityById(id, &entities.UserIndividualEntity{}, executor)
}

// dbからid指定で個人顧客情報を取得する
func (r *UserRepository) getUserIndividualEntityById(id int, entity *entities.UserIndividualEntity, executor gorp.SqlExecutor) (*entities.UserIndividualEntity, error) {
	userIndividualView := data_mappers.UserIndividualView{}
	noRow, err := r.selectOne(
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
func (r *UserRepository) SaveUserCorporation(userEntity *entities.UserCorporationEntity, executor gorp.SqlExecutor) (savedId int, err error) {
	// userRecord作成
	userRecord := data_mappers.NewUserMapperFromUserCorporationEntity(userEntity)

	// 保存
	err = executor.Insert(userRecord)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	// userCorporationRecord作成
	userCorporationRecord := data_mappers.NewUserCorporationMapperFromUserCorporationEntity(userEntity)
	userCorporationRecord.UserId = userRecord.Id

	// 保存
	err = executor.Insert(userCorporationRecord)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return userRecord.Id, nil

}

// dbからid指定で法人顧客情報を取得する
func (r *UserRepository) getUserCorporationEntityById(id int, entity *entities.UserCorporationEntity, executor gorp.SqlExecutor) (*entities.UserCorporationEntity, error) {
	// dbからデータ取得
	record := data_mappers.UserCorporationView{}
	query := "SELECT u.id, uc.company_name, uc.contact_person_name, uc.president_name, u.created_at, u.updated_at " +
		"FROM users u " +
		"inner join users_corporation uc on u.id = uc.user_id " +
		"WHERE id = $1"
	noRow, err := r.selectOne(executor, &record, entity, query, id)
	if err != nil {
		return nil, err
	}
	if noRow {
		return nil, nil
	}
	return entity, nil
}

// Idで法人顧客情報を取得する。データがなければnilを返す
func (r *UserRepository) GetUserCorporationById(id int, executor gorp.SqlExecutor) (*entities.UserCorporationEntity, error) {
	return r.getUserCorporationEntityById(id, &entities.UserCorporationEntity{}, executor)
}

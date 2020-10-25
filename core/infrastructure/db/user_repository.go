package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	entities "github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/data_mappers"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"reflect"
	"time"
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
	user := NewUserMapperFromUserIndividualEntity(userEntity)

	err = executor.Insert(user)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	// individualを保存
	userIndividualDbMap := NewUserIndividualMapperFromUserIndividualEntity(userEntity)
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
   u.created_at as user_individual_created_at,
   u.updated_at as user_individual_updated_at,
   uc.corporation_name as user_corporation_corporation_name,
   uc.contact_person_name as user_corporation_contact_person_name,
   uc.president_name as user_corporation_president_name,
   u.created_at as user_corporation_created_at,
   u.updated_at as user_corporation_updated_at
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
	//case "corporation":
	//	retUser, err := entities.NewUserCorporationEntityWithData(mapper.UserId, mapper.UserCorporationCorporationName.String, mapper.UserCorporationContractPersonName.String, mapper.UserCorporationPresidentName.String, mapper.UserCorporationCreatedAt.Time, mapper.UserCorporationUpdatedAt.Time)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return retUser, nil
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
	userIndividualView := UserIndividualView{}
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
	userRecord := NewUserMapperFromUserCorporationEntity(userEntity)

	// 保存
	err = executor.Insert(userRecord)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	// userCorporationRecord作成
	userCorporationRecord := NewUserCorporationMapperFromUserCorporationEntity(userEntity)
	userCorporationRecord.UserId = userRecord.Id

	// 保存
	err = executor.Insert(userCorporationRecord)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return userRecord.Id, nil

}

// dbからid指定で法人顧客情報を取得する
//func (r *UserRepository) getUserCorporationEntityById(id int, entity *entities.UserCorporationEntity, executor gorp.SqlExecutor) (*entities.UserCorporationEntity, error) {
//	// dbからデータ取得
//	record := data_mappers.UserCorporationView{}
//	query := "SELECT u.id, uc.corporation_name, uc.contact_person_name, uc.president_name, u.created_at, u.updated_at " +
//		"FROM users u " +
//		"inner join users_corporation uc on u.id = uc.user_id " +
//		"WHERE id = $1"
//	noRow, err := r.selectOne(executor, &record, entity, query, id)
//	if err != nil {
//		return nil, err
//	}
//	if noRow {
//		return nil, nil
//	}
//	return entity, nil
//}

// Idで法人顧客情報を取得する。データがなければnilを返す
//func (r *UserRepository) GetUserCorporationById(id int, executor gorp.SqlExecutor) (*entities.UserCorporationEntity, error) {
//	return r.getUserCorporationEntityById(id, &entities.UserCorporationEntity{}, executor)
//}

type UserMapper struct {
	Id int `db:"id"`
	CreatedAtUpdatedAtMapper
}

func (u *UserMapper) SetDataToEntity(entity interface{}) error {
	userEntity, ok := entity.(*entities.UserEntity)
	if !ok {
		return errors.Errorf("entityが*entities.UserEntityではない。%v", reflect.TypeOf(entity))
	}
	userEntity.LoadData(u.Id, u.CreatedAt, u.UpdatedAt)
	return nil
}

// UserIndividualEntityからデータを読み込んでUser(DBマッピング用)を作成する
func NewUserMapperFromUserIndividualEntity(userIndividual *entities.UserIndividualEntity) *UserMapper {
	return &UserMapper{
		Id: userIndividual.Id(),
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: userIndividual.CreatedAt(),
			UpdatedAt: userIndividual.UpdatedAt(),
		},
	}
}

// UserCorporationEntityからデータを読み込んでUser(DBマッピング用)を作成する
func NewUserMapperFromUserCorporationEntity(userCorporation *entities.UserCorporationEntity) *UserMapper {
	return &UserMapper{
		Id: userCorporation.Id(),
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: userCorporation.CreatedAt(),
			UpdatedAt: userCorporation.UpdatedAt(),
		},
	}
}

type UserIndividualMapper struct {
	UserId int    `db:"user_id"`
	Name   string `db:"name"`
	CreatedAtUpdatedAtMapper
}

// UserIndividualEntity Entityからデータを読み込んでUserIndividual(DBマッピング用)を作成する
func NewUserIndividualMapperFromUserIndividualEntity(entity *entities.UserIndividualEntity) *UserIndividualMapper {
	return &UserIndividualMapper{
		UserId: entity.Id(),
		Name:   entity.Name(),
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: entity.CreatedAt(),
			UpdatedAt: entity.UpdatedAt(),
		},
	}
}

type UserView struct {
	UserId                  int            `db:"user_id"`
	UserType                string         `db:"user_type"`
	UserIndividualName      sql.NullString `db:"user_individual_name"`
	UserIndividualCreatedAt sql.NullTime   `db:"user_individual_created_at"`
	UserIndividualUpdatedAt sql.NullTime   `db:"user_individual_updated_at"`

	UserCorporationCorporationName    sql.NullString `db:"user_corporation_corporation_name"`
	UserCorporationContractPersonName sql.NullString `db:"user_corporation_contact_person_name"`
	UserCorporationPresidentName      sql.NullString `db:"user_corporation_president_name"`
	UserCorporationCreatedAt          sql.NullTime   `db:"user_corporation_created_at"`
	UserCorporationUpdatedAt          sql.NullTime   `db:"user_corporation_updated_at"`
}

type EntitySetter interface {
	SetDataToEntity(entity interface{}) error
}

type IBaseEntity interface {
	Id() int
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type UserIndividualView struct {
	UserMapper
	Name string
}

func (u *UserIndividualView) SetDataToEntity(entity interface{}) error {
	value, ok := entity.(*entities.UserIndividualEntity)
	if !ok {
		return errors.New(fmt.Sprintf("*entities.UserIndividualEntity型ではないものが渡ってきた。 %v", entity))
	}

	err := value.LoadData(u.Id, u.Name, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

type UserCorporationMapper struct {
	UserId            int    `db:"user_id"`
	ContactParsonName string `db:"contact_person_name"`
	PresidentName     string `db:"president_name"`
	CorporationName   string `db:"corporation_name"`
	CreatedAtUpdatedAtMapper
}

// UserIndividualEntity Entityからデータを読み込んでUserIndividual(DBマッピング用)を作成する
func NewUserCorporationMapperFromUserCorporationEntity(entity *entities.UserCorporationEntity) *UserCorporationMapper {
	return &UserCorporationMapper{
		UserId:            entity.Id(),
		CorporationName:   entity.CorporationName(),
		ContactParsonName: entity.ContactPersonName(),
		PresidentName:     entity.PresidentName(),
		CreatedAtUpdatedAtMapper: CreatedAtUpdatedAtMapper{
			CreatedAt: entity.CreatedAt(),
			UpdatedAt: entity.UpdatedAt(),
		},
	}
}

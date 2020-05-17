package tables

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user"
	"gopkg.in/gorp.v2"
	"time"
)

type UserRecord struct {
	Id        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// UserIndividualEntityからデータを読み込んでUser(DBマッピング用)を作成する
func NewUserRecordFromUserIndividualEntity(userIndividual *user.UserIndividualEntity) *UserRecord {
	return &UserRecord{
		Id:        userIndividual.Id(),
		CreatedAt: userIndividual.CreatedAt(),
		UpdatedAt: userIndividual.UpdatedAt(),
	}
}

// UserCorporationEntityからデータを読み込んでUser(DBマッピング用)を作成する
func NewUserRecordFromUserCorporationEntity(userCorporation *user.UserCorporationEntity) *UserRecord {
	return &UserRecord{
		Id:        userCorporation.Id(),
		CreatedAt: userCorporation.CreatedAt(),
		UpdatedAt: userCorporation.UpdatedAt(),
	}
}

// insert時に時刻をセットするhook
func (u *UserRecord) PreInsert(s gorp.SqlExecutor) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	return nil
}

// updateに時刻をセットするhook
func (u *UserRecord) PreUpdate(s gorp.SqlExecutor) error {
	u.UpdatedAt = time.Now()
	return nil
}

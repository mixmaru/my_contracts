package db

import (
	"database/sql"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type BaseRepository struct {
}

func (b *BaseRepository) selectOne(executor gorp.SqlExecutor, record EntitySetter, entity IBaseEntity, query string, args ...interface{}) (noRow bool, err error) {
	err = executor.SelectOne(record, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			// データがない
			return true, nil
		} else {
			return true, errors.WithStack(err)
		}
	}

	// エンティティに詰める
	err = record.SetDataToEntity(entity)
	if err != nil {
		return true, err
	}

	return false, nil
}

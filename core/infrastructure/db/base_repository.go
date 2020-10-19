package db

import (
	"database/sql"
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/data_mappers"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

type BaseRepository struct {
}

func (b *BaseRepository) selectOne(executor gorp.SqlExecutor, record data_mappers.EntitySetter, entity entities.IBaseEntity, query string, args ...interface{}) (noRow bool, err error) {
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

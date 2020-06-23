package repositories

import (
	"database/sql"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	tables2 "github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/tables"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
)

func selectOne(executor gorp.SqlExecutor, record tables2.IRecord, entity entities.IBaseEntity, query string, args ...interface{}) (noRow bool, err error) {
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

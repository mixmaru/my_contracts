package db

import (
	"database/sql"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"strconv"
	"strings"
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

/*
与えられた件数のIN QUERY用の"IN ($1, $2, $3)"の"$1, $2, $3"という文字列を生成する
*/
func CrateInStatement(num int) string {
	tmpSlice := make([]string, 0, num)
	for i := 1; i <= num; i++ {
		tmpSlice = append(tmpSlice, "$"+strconv.Itoa(i))
	}
	return strings.Join(tmpSlice, ", ")
}

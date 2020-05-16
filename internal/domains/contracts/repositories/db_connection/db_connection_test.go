package db_connection

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
)

func TestDbConnection_GetConnection(t *testing.T) {

	t.Run("トランザクションが渡されなかった場合 nilが渡された場合 dbMapが返る", func(t *testing.T) {
		conn, err := GetConnectionIfNotTransaction(nil)
		assert.NoError(t, err)
		assert.IsType(t, &gorp.DbMap{}, conn)
	})

	t.Run("トランザクションが渡された場合 トランザクションが返る", func(t *testing.T) {
		// トランザクション取得
		conn1, err := GetConnectionIfNotTransaction(nil)
		assert.NoError(t, err)
		dbMap, ok := conn1.(*gorp.DbMap)
		assert.True(t, ok)
		tran, err := dbMap.Begin()
		assert.NoError(t, err)

		// トランザクションを渡す
		conn2, err := GetConnectionIfNotTransaction(tran)
		assert.NoError(t, err)

		assert.IsType(t, &gorp.Transaction{}, conn2)

	})
}

func TestDbConnection_Close(t *testing.T) {
	t.Run("dbMapが渡された場合dbMapをcloseする", func(t *testing.T) {
		// dbConnection取得
		conn, err := GetConnectionIfNotTransaction(nil)
		assert.NoError(t, err)
		err = CloseConnectionIfNotTransaction(conn)
		assert.NoError(t, err)

		// 接続しようとしたらエラーになるはず
		dbMap, ok := conn.(*gorp.DbMap)
		assert.True(t, ok)
		err = dbMap.Db.Ping()
		assert.Error(t, err)
	})

	t.Run("トランザクションが渡された場合何もしない", func(t *testing.T) {
		// トランザクションを作成
		conn, err := GetConnectionIfNotTransaction(nil)
		assert.NoError(t, err)
		dbMap, ok := conn.(*gorp.DbMap)
		assert.True(t, ok)
		tran, err := dbMap.Begin()
		assert.NoError(t, err)

		// 実行
		err = CloseConnectionIfNotTransaction(tran)
		assert.NoError(t, err)

		// 接続できるはず
		_, err = tran.Query("select count(*) from users")
		assert.NoError(t, err)
	})
}

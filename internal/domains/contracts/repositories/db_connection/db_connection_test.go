package db_connection

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDbConnection_GetConnection(t *testing.T) {
	dbMap, err := GetConnection()
	assert.NoError(t, err)
	defer dbMap.Db.Close()
}

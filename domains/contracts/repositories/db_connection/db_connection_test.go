package db_connection

import (
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"os"
	"testing"
)

func TestDbConnection_GetConnection(t *testing.T) {
	dbMap, err := GetConnection()
	assert.NoError(t, err)
	assert.IsType(t, &gorp.DbMap{}, dbMap)
}

func TestDbConnection_getConnectionString(t *testing.T) {
	// 環境変数用意
	os.Setenv("DB_TEST_HOST", "test_host")
	os.Setenv("DB_TEST_NAME", "test_name")
	os.Setenv("DB_TEST_USER", "test_user")
	os.Setenv("DB_TEST_PASSWORD", "test_pass")
	os.Setenv("DB_DEVELOPMENT_HOST", "dev_host")
	os.Setenv("DB_DEVELOPMENT_NAME", "dev_name")
	os.Setenv("DB_DEVELOPMENT_USER", "dev_user")
	os.Setenv("DB_DEVELOPMENT_PASSWORD", "dev_pass")
	os.Setenv("DB_PRODUCTION_HOST", "prod_host")
	os.Setenv("DB_PRODUCTION_NAME", "prod_name")
	os.Setenv("DB_PRODUCTION_USER", "prod_user")
	os.Setenv("DB_PRODUCTION_PASSWORD", "prod_pass")

	t.Run("testの時", func(t *testing.T) {
		str, err := getConnectionString(utils.Test)
		assert.NoError(t, err)
		assert.Equal(t, "host=test_host user=test_user dbname=test_name password=test_pass sslmode=disable", str)
	})

	t.Run("developmentの時", func(t *testing.T) {
		str, err := getConnectionString(utils.Development)
		assert.NoError(t, err)
		assert.Equal(t, "host=dev_host user=dev_user dbname=dev_name password=dev_pass sslmode=disable", str)
	})

	t.Run("productionの時", func(t *testing.T) {
		str, err := getConnectionString(utils.Production)
		assert.NoError(t, err)
		assert.Equal(t, "host=prod_host user=prod_user dbname=prod_name password=prod_pass sslmode=disable", str)
	})

	t.Run("それ意外の時", func(t *testing.T) {
		str, err := getConnectionString("other")
		assert.Error(t, err)
		assert.Equal(t, "", str)
	})
}

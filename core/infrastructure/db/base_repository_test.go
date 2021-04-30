package db

import (
	"database/sql"
	"testing"

	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
)

func TestBaseRepository_selectOne(t *testing.T) {
	db, err := GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	baseRepository := &BaseRepository{}

	t.Run("データがある時_マッパーとエンティティとクエリを渡すとマッパーを使ってデータを取り出しエンティティにデータを詰めてくれる", func(t *testing.T) {
		// データ取得
		productMapper := productGetMapper{}
		productEntity := product.ProductEntity{}
		query := `
SELECT
       1 AS id,
       '商品名' AS name,
       to_timestamp('2020-01-01', 'YYYY-MM-DD') AS created_at,
       to_timestamp('2020-01-02', 'YYYY-MM-DD') AS updated_at,
       true AS exist_price_monthly,
       '200' AS price_monthly
;
`
		noRow, err := baseRepository.selectOne(db, &productMapper, &productEntity, query)
		assert.NoError(t, err)
		assert.False(t, noRow)

		assert.Equal(t, 1, productEntity.Id())
		assert.Equal(t, "商品名", productEntity.Name())
		price, exist := productEntity.MonthlyPrice()
		assert.True(t, exist)
		assert.Equal(t, "200", price.String())
		assert.NotZero(t, productEntity.CreatedAt())
		assert.NotZero(t, productEntity.UpdatedAt())
	})

	t.Run("データがない時_noRowがtrueで返る", func(t *testing.T) {
		productRecord := productGetMapper{}
		productEntity := product.ProductEntity{}
		noRow, err := baseRepository.selectOne(db, &productRecord, &productEntity, "select * from products where 1 = 2")
		assert.NoError(t, err)
		assert.True(t, noRow)
	})

	t.Run("渡すrecordとentityがアベコベだったときはエラーが返る_noRowもTrueで返る", func(t *testing.T) {
		productRecord := productGetMapper{}
		userCorporationEntity := user.UserCorporationEntity{}
		query := `
SELECT
       1 AS id,
       '商品名' AS name,
       to_timestamp('2020-01-01', 'YYYY-MM-DD') AS created_at,
       to_timestamp('2020-01-02', 'YYYY-MM-DD') AS updated_at,
       true AS exist_price_monthly,
       '200' AS price_monthly
;
`
		noRow, err := baseRepository.selectOne(db, &productRecord, &userCorporationEntity, query)
		assert.Error(t, err, "aaa")
		assert.True(t, noRow)
	})
}

func TestTest(t *testing.T) {
	//query := "select id from customers where id in (:ids);"
	//query := "select id from customers where id in ($1);"
	query := "select id from customers where id in (?);"

	executeMode, err := utils.GetExecuteMode()
    assert.NoError(t, err)

	connectionStr, err := getConnectionString(executeMode)
    assert.NoError(t, err)

	db, err := sql.Open("postgres", connectionStr)
	dbmap := &gorp.DbMap{
		Db:              db, // コネクション
		Dialect:         gorp.PostgresDialect{},
		ExpandSliceArgs: true,
	}
    mapper := []int{}
	//_, err = dbmap.Select(&mapper, query, map[string]interface{}{"ids": []int{3,4,7}})
	_, err = dbmap.Select(&mapper, query, []int{3,4,7})
    assert.NoError(t, err)
    assert.Equal(t, []int{3,4,7}, mapper)
}

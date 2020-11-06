package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/stretchr/testify/assert"
	"testing"
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

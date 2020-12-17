package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestCustomerTypeRepository_Create_And_Get(t *testing.T) {
	t.Run("CustomerTypeエンティティを渡すとDBへ保存される", func(t *testing.T) {
		conn, err := GetConnection()
		assert.NoError(t, err)
		tran, err := conn.Begin()
		assert.NoError(t, err)

		////// 準備
		customerType, err := createNewCustomerTypeEntity()
		assert.NoError(t, err)

		////// 実行
		r := NewCustomerTypeRepository()
		savedId, err := r.Create(customerType, tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		////// 検証 todo: あとで再取得してデータが取れるか確認する
		assert.NotZero(t, savedId)

		////// データの再取得
		reloadedEntity, err := r.GetById(savedId, conn)
		assert.NoError(t, err)

		////// 検証
		expected := customer.NewCustomerTypeEntityWithData(savedId, customerType.Name(), customerType.CustomerPropertyTypeIds())
		assert.True(t, reflect.DeepEqual(expected, reloadedEntity))
	})
}

func createNewCustomerTypeEntity() (*customer.CustomerTypeEntity, error) {
	// カスタマープロパティタイプの登録
	timestamp := time.Now().UnixNano()
	timestampstr := strconv.Itoa(int(timestamp))
	customerProperties := []*customer.CustomerPropertyTypeEntity{
		customer.NewCustomerPropertyTypeEntity("性別"+timestampstr, customer.PROPERTY_TYPE_STRING),
		customer.NewCustomerPropertyTypeEntity("年齢"+timestampstr, customer.PROPERTY_TYPE_NUMERIC),
		customer.NewCustomerPropertyTypeEntity("住所"+timestampstr, customer.PROPERTY_TYPE_STRING),
	}
	propertyTypeRep := NewCustomerPropertyTypeRepository()
	conn, err := GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Db.Close()

	savedPropertyIds, err := propertyTypeRep.Create(customerProperties, conn)
	if err != nil {
		return nil, err
	}
	// カスタマータイプエンティティの作成
	customerType := customer.NewCustomerTypeEntity("顧客名"+timestampstr, savedPropertyIds)
	return customerType, nil
}

func TestCustomerTypeRepository_GetById(t *testing.T) {
	t.Run("存在しないidを渡されるとnilを返す", func(t *testing.T) {
		conn, err := GetConnection()
		assert.NoError(t, err)

		////// 実行
		r := NewCustomerTypeRepository()
		reloadedEntity, err := r.GetById(-10000, conn)
		assert.NoError(t, err)

		////// 検証
		assert.Nil(t, reloadedEntity)
	})
}

func TestCustomerTypeRepository_GetByName(t *testing.T) {
	//// 準備
	preInsertCustomerType, err := createNewCustomerTypeEntity()

	assert.NoError(t, err)
	conn, err := GetConnection()
	assert.NoError(t, err)
	defer conn.Db.Close()
	tran, err := conn.Begin()
	assert.NoError(t, err)

	r := NewCustomerTypeRepository()
	preSavedId, err := r.Create(preInsertCustomerType, tran)
	assert.NoError(t, err)

	type args struct {
		name     string
		executor gorp.SqlExecutor
	}
	tests := []struct {
		name       string
		args       args
		wantEntity *customer.CustomerTypeEntity
		wantErr    bool
	}{
		{
			"nameを渡すとentityが返ってくる",
			args{
				preInsertCustomerType.Name(),
				tran,
			},
			customer.NewCustomerTypeEntityWithData(preSavedId, preInsertCustomerType.Name(), preInsertCustomerType.CustomerPropertyTypeIds()),
			false,
		},
		{
			"存在しないnameを渡すとnilが返ってくる",
			args{
				"notExistName_aaaaaaaaaabbbbbbbbbbcccccccccc",
				tran,
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CustomerTypeRepository{}
			gotEntity, err := r.GetByName(tt.args.name, tt.args.executor)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEntity, tt.wantEntity) {
				t.Errorf("GetByName() gotEntity = %v, want %v", gotEntity, tt.wantEntity)
			}
		})
	}
}

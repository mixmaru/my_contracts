package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// UserIndividualのインスタンス化をテスト
func TestUserIndividual_Instantiate(t *testing.T) {
	// インスタンス化
	userIndividual := &UserIndividual{}
	userIndividual.SetName(NewName("顧客太郎"))

	// テスト
	assert.Equal(t, NewName("顧客太郎"), userIndividual.Name())
}

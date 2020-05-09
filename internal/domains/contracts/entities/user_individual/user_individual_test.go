package user_individual

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// UserIndividualのインスタンス化をテスト
func TestUserIndividual_Instantiate(t *testing.T) {
	// インスタンス化
	userIndividual := &UserIndividual{}
	userIndividual.SetName("顧客太郎")

	// テスト
	assert.Equal(t, "顧客太郎", userIndividual.Name())
}

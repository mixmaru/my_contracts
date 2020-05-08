package user

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/common_values"
	"github.com/stretchr/testify/assert"
	"testing"
)

// UserIndividualのインスタンス化をテスト
func TestUserIndividual_Instantiate(t *testing.T) {
	// インスタンス化
	userIndividual := &UserIndividual{}
	userIndividual.SetName(common_values.NewName("顧客太郎"))

	// テスト
	assert.Equal(t, common_values.NewName("顧客太郎"), userIndividual.Name())
}

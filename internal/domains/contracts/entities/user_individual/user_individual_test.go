package user_individual

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user_individual/values"
	"github.com/stretchr/testify/assert"
	"testing"
)

// UserIndividualのインスタンス化をテスト
func TestUserIndividual_Instantiate(t *testing.T) {
	// インスタンス化
	userIndividual := &UserIndividual{}
	userIndividual.SetName(values.NewName("顧客太郎"))

	// テスト
	assert.Equal(t, values.NewName("顧客太郎"), userIndividual.Name())
}

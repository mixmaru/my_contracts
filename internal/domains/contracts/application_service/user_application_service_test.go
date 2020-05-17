package application_service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserApplication_NewUserApplicationService(t *testing.T) {
	// リポジトリモックを用意する
	userApp := NewUserApplicationService()
	assert.IsType(t, &UserApplicationService{}, userApp)
}

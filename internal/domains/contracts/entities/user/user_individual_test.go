package user

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/common_values"
	"testing"
)

func TestUserIndividual_Test222(t *testing.T) {
	userIndividual := &UserIndividual{}
	userIndividual.SetId(common_values.NewIdInt(1))
	userIndividual.SetName(NewName("何がしかの名前"))

	testArray := []IUser{
		&User{
			id: common_values.NewIdInt(1),
		},
		userIndividual,
	}
	for _, user := range testArray {
		fmt.Printf("%+v", user)
	}
}

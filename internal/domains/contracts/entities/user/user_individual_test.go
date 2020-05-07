package user

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/user/values"
	"testing"
)

func TestUserIndividual_Test222(t *testing.T) {
	userIndividual := &UserIndividual{}
	userIndividual.SetId(values.NewIdInt(1))
	userIndividual.SetName(NewName("何がしかの名前"))

	testArray := []IUser{
		&User{
			id: values.NewIdInt(1),
		},
		userIndividual,
	}
	for _, user := range testArray {
		fmt.Printf("%+v", user)
	}
}

package product

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductRepository_Save(t *testing.T) {
	r := ProductRepository{}
	_, err := r.Save(nil, nil)
	assert.NoError(t, err)
}

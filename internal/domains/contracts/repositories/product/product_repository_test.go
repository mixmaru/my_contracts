package product

import (
	"testing"
)

func TestProductRepository_Save(t *testing.T) {
	r := ProductRepository{}
	r.Save(nil, nil)
}

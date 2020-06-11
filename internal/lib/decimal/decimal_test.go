package decimal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecimal_New(t *testing.T) {
	decimal := New(100, 10)
	assert.NotNil(t, decimal)
}

func TestDecimal_NewFromFloat(t *testing.T) {
	decimal := NewFromFloat(100.01)
	assert.NotNil(t, decimal)
}

func TestDecimal_NewFromInt(t *testing.T) {
	decimal := NewFromInt(100)
	assert.NotNil(t, decimal)
}

type TestCase struct {
	Left   float64
	Right  float64
	Result float64
}

// 加算
func TestDecimal_Add(t *testing.T) {
	type Inp struct {
		a float64
		b float64
	}

	inputs := map[Inp]float64{
		Inp{2, 3}:                     5,
		Inp{2454495034, 3451204593}:   5905699627,
		Inp{24544.95034, .3451204593}: 24545.2954604593,
		Inp{.1, .1}:                   0.2,
		Inp{.1, -.1}:                  0,
		Inp{0, 1.001}:                 1.001,
	}

	for input, expect := range inputs {
		decimalA := NewFromFloat(input.a)
		decimalB := NewFromFloat(input.b)
		result := decimalA.Add(decimalB)
		expect := NewFromFloat(expect)
		if !result.Equal(expect) {
			t.Errorf("%v != %v", result.decimal, expect.decimal)
		}
	}
}

// 減産

// 積算

// 割り算

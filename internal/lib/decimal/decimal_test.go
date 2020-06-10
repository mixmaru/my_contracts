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
	testCases := []TestCase{
		{
			Left:   100,
			Right:  200,
			Result: 300,
		},
		{
			Left:   100,
			Right:  -200,
			Result: -100,
		},
		{
			Left:   100.1,
			Right:  100,
			Result: 200.1,
		},
	}

	for _, testCase := range testCases {
		decimal1 := NewFromFloat(testCase.Left)
		decimal2 := NewFromFloat(testCase.Right)
		result := decimal1.Add(decimal2)
		expect := NewFromFloat(testCase.Result)
		if !result.Equal(expect) {
			t.Errorf("%v != %v", result.decimal, expect.decimal)
		}
	}
}

// 減産

// 積算

// 割り算

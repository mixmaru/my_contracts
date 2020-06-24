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

func TestDecimal_String(t *testing.T) {
	decimal := NewFromFloat(100.01)
	assert.Equal(t, "100.01", decimal.String())
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

// 引き算
func TestDecimal_Sub(t *testing.T) {
	type Inp struct {
		a float64
		b float64
	}

	inputs := map[Inp]float64{
		Inp{2, 3}:                     -1,
		Inp{12, 3}:                    9,
		Inp{-2, 9}:                    -11,
		Inp{2454495034, 3451204593}:   -996709559,
		Inp{24544.95034, .3451204593}: 24544.6052195407,
		Inp{.1, -.1}:                  0.2,
		Inp{.1, .1}:                   0,
		Inp{0, 1.001}:                 -1.001,
		Inp{1.001, 0}:                 1.001,
		Inp{2.3, .3}:                  2,
	}

	for input, expect := range inputs {
		decimalA := NewFromFloat(input.a)
		decimalB := NewFromFloat(input.b)
		result := decimalA.Sub(decimalB)
		expect := NewFromFloat(expect)
		if !result.Equal(expect) {
			t.Errorf("%v != %v", result.decimal, expect.decimal)
		}
	}
}

// 積算
func TestDecimal_Mul(t *testing.T) {
	type Inp struct {
		a string
		b string
	}

	inputs := map[Inp]string{
		Inp{"2", "3"}:                     "6",
		Inp{"2454495034", "3451204593"}:   "8470964534836491162",
		Inp{"24544.95034", ".3451204593"}: "8470.964534836491162",
		Inp{".1", ".1"}:                   "0.01",
		Inp{"0", "1.001"}:                 "0",
	}

	for input, expect := range inputs {
		decimalA, err := NewFromString(input.a)
		assert.NoError(t, err)

		decimalB, err := NewFromString(input.b)
		assert.NoError(t, err)

		result := decimalA.Mul(decimalB)

		expect, err := NewFromString(expect)
		assert.NoError(t, err)

		if !result.Equal(expect) {
			t.Errorf("%v != %v", result.decimal, expect.decimal)
		}
	}
}

// 割り算
func TestDecimal_Div(t *testing.T) {
	type Inp struct {
		a string
		b string
	}

	inputs := map[Inp]string{
		Inp{"6", "3"}:                            "2",
		Inp{"10", "2"}:                           "5",
		Inp{"2.2", "1.1"}:                        "2",
		Inp{"-2.2", "-1.1"}:                      "2",
		Inp{"12.88", "5.6"}:                      "2.3",
		Inp{"1023427554493", "43432632"}:         "23563.5628642767953828", // rounded
		Inp{"1", "434324545566634"}:              "0.0000000000000023",
		Inp{"1", "3"}:                            "0.3333333333333333",
		Inp{"2", "3"}:                            "0.6666666666666667", // rounded
		Inp{"10000", "3"}:                        "3333.3333333333333333",
		Inp{"10234274355545544493", "-3"}:        "-3411424785181848164.3333333333333333",
		Inp{"-4612301402398.4753343454", "23.5"}: "-196268144782.9138440146978723",
	}

	for input, expect := range inputs {
		decimalA, err := NewFromString(input.a)
		assert.NoError(t, err)

		decimalB, err := NewFromString(input.b)
		assert.NoError(t, err)

		result := decimalA.Div(decimalB)

		expect, err := NewFromString(expect)
		assert.NoError(t, err)

		if !result.Equal(expect) {
			t.Errorf("%v != %v", result.decimal, expect.decimal)
		}
	}
}

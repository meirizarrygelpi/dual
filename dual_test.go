package dual

import (
	"fmt"
	"testing"
)

var (
	zero = New(0, 0)
	e0   = New(1, 0)
	e1   = New(0, 1)
)

func TestEquals(t *testing.T) {
	var tests = []struct {
		x    *Dual
		y    *Dual
		want bool
	}{
		{zero, zero, true},
		{e0, e0, true},
		{e1, e1, true},
		{e0, e1, false},
		{e1, e0, false},
		{New(2.03, 3), New(2.0299999999, 3), true},
		{New(1, 2), New(3, 4), false},
	}

	for _, test := range tests {
		if got := test.x.Equals(test.y); got != test.want {
			t.Errorf("Equals(%v, %v) = %v", test.x, test.y, got)
		}
	}
}

func TestSet(t *testing.T) {
	var tests = []struct {
		x    *Dual
		want *Dual
	}{
		{zero, zero},
		{New(1, 2), New(1, 2)},
	}

	for _, test := range tests {
		if got := new(Dual).Set(test.x); !got.Equals(test.want) {
			t.Errorf("Set(%v) = %v, want %v", test.x, got, test.want)
		}
	}
}

func TestString(t *testing.T) {
	var tests = []struct {
		x    *Dual
		want string
	}{
		{zero, "(0+0ε)"},
		{e0, "(1+0ε)"},
		{e1, "(0+1ε)"},
		{New(1, 1), "(1+1ε)"},
		{New(1, -1), "(1-1ε)"},
		{New(-1, 1), "(-1+1ε)"},
		{New(-1, -1), "(-1-1ε)"},
	}

	for _, test := range tests {
		if got := test.x.String(); got != test.want {
			t.Errorf("String(%v) = %v, want %v", test.x, got, test.want)
		}
	}
}

func ExampleNew() {
	fmt.Println(New(1, 0))
	fmt.Println(New(0, 1))
	fmt.Println(New(2, -3))
	fmt.Println(New(-4, 5))
	// Output:
	// (1+0ε)
	// (0+1ε)
	// (2-3ε)
	// (-4+5ε)
}

func TestScalar(t *testing.T) {}

func TestNeg(t *testing.T) {}

func TestConj(t *testing.T) {}

func TestAdd(t *testing.T) {}

func TestSub(t *testing.T) {}

func TestMul(t *testing.T) {}

func TestQuad(t *testing.T) {}

func TestIsZeroDiv(t *testing.T) {
	var tests = []struct {
		z    *Dual
		want bool
	}{
		{zero, true},
		{e0, false},
		{e1, true},
	}

	for _, test := range tests {
		if got := test.z.IsZeroDiv(); got != test.want {
			t.Errorf("IsZeroDiv(%v) = %v", test.z, got)
		}
	}
}

func TestQuo(t *testing.T) {}

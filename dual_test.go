package dual

import (
	"fmt"
	"math"
	"testing"
)

var (
	zero = New(0, 0)
	e0   = New(1, 0)
	e1   = New(0, 1)
)

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

func TestCopy(t *testing.T) {
	var tests = []struct {
		x    *Dual
		want *Dual
	}{
		{zero, zero},
		{New(1, 2), New(1, 2)},
	}
	for _, test := range tests {
		if got := new(Dual).Copy(test.x); !got.Equals(test.want) {
			t.Errorf("Copy(%v) = %v, want %v", test.x, got, test.want)
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

func TestScal(t *testing.T) {
	var tests = []struct {
		z    *Dual
		a    float64
		want *Dual
	}{
		{zero, 1, zero},
		{New(1, 2), 3, New(3, 6)},
		{New(1, 2), 0, zero},
	}
	for _, test := range tests {
		if got := new(Dual).Scal(test.z, test.a); !got.Equals(test.want) {
			t.Errorf("Scal(%v, %v) = %v, want %v",
				test.z, test.a, got, test.want)
		}
	}
}

func TestNeg(t *testing.T) {
	var tests = []struct {
		z    *Dual
		want *Dual
	}{
		{zero, zero},
		{e0, New(-1, 0)},
		{e1, New(0, -1)},
		{New(3, 4), New(-3, -4)},
	}
	for _, test := range tests {
		if got := new(Dual).Neg(test.z); !got.Equals(test.want) {
			t.Errorf("Neg(%v) = %v, want %v",
				test.z, got, test.want)
		}
	}
}

func TestConj(t *testing.T) {
	var tests = []struct {
		z    *Dual
		want *Dual
	}{
		{zero, zero},
		{e0, e0},
		{e1, New(0, -1)},
		{New(3, 4), New(3, -4)},
	}
	for _, test := range tests {
		if got := new(Dual).Conj(test.z); !got.Equals(test.want) {
			t.Errorf("Conj(%v) = %v, want %v",
				test.z, got, test.want)
		}
	}
}

func TestAdd(t *testing.T) {
	var tests = []struct {
		x    *Dual
		y    *Dual
		want *Dual
	}{
		{zero, zero, zero},
		{e0, e0, New(2, 0)},
		{e1, e1, New(0, 2)},
		{e0, e1, New(1, 1)},
		{e1, e0, New(1, 1)},
	}
	for _, test := range tests {
		if got := new(Dual).Add(test.x, test.y); !got.Equals(test.want) {
			t.Errorf("Add(%v, %v) = %v, want %v",
				test.x, test.y, got, test.want)
		}
	}
}

func TestSub(t *testing.T) {
	var tests = []struct {
		x    *Dual
		y    *Dual
		want *Dual
	}{
		{zero, zero, zero},
		{e0, e0, zero},
		{e1, e1, zero},
		{e0, e1, New(1, -1)},
		{e1, e0, New(-1, 1)},
	}
	for _, test := range tests {
		if got := new(Dual).Sub(test.x, test.y); !got.Equals(test.want) {
			t.Errorf("Sub(%v, %v) = %v, want %v",
				test.x, test.y, got, test.want)
		}
	}
}

func TestMul(t *testing.T) {
	var tests = []struct {
		x    *Dual
		y    *Dual
		want *Dual
	}{
		{zero, zero, zero},
		{e0, e0, e0},
		{e1, e1, zero},
		{e0, e1, e1},
		{e1, e0, e1},
	}
	for _, test := range tests {
		if got := new(Dual).Mul(test.x, test.y); !got.Equals(test.want) {
			t.Errorf("Mul(%v, %v) = %v, want %v",
				test.x, test.y, got, test.want)
		}
	}
}

func TestQuad(t *testing.T) {
	var tests = []struct {
		z    *Dual
		want float64
	}{
		{zero, 0},
		{e0, 1},
		{e1, 0},
		{New(-2, 1), 4},
	}
	for _, test := range tests {
		if got := test.z.Quad(); notEquals(got, test.want) {
			t.Errorf("Quad(%v) = %v, want %v",
				test.z, got, test.want)
		}
	}
}

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

func TestInv(t *testing.T) {
	var tests = []struct {
		x    *Dual
		want *Dual
	}{
		{e0, e0},
		{New(2, 0), New(0.5, 0)},
	}
	for _, test := range tests {
		if got := new(Dual).Inv(test.x); !got.Equals(test.want) {
			t.Errorf("Inv(%v) = %v, want %v",
				test.x, got, test.want)
		}
	}
}

func TestQuo(t *testing.T) {
	var tests = []struct {
		x    *Dual
		y    *Dual
		want *Dual
	}{
		{e0, e0, e0},
		{New(0.5, 0), New(2, 0), New(0.25, 0)},
	}
	for _, test := range tests {
		if got := new(Dual).Quo(test.x, test.y); !got.Equals(test.want) {
			t.Errorf("Quo(%v, %v) = %v, want %v",
				test.x, test.y, got, test.want)
		}
	}
}

func TestIsInf(t *testing.T) {
	var tests = []struct {
		z    *Dual
		want bool
	}{
		{zero, false},
		{e0, false},
		{e1, false},
		{New(math.Inf(0), 4), true},
	}
	for _, test := range tests {
		if got := test.z.IsInf(); got != test.want {
			t.Errorf("IsInf(%v) = %v", test.z, got)
		}
	}
}

func ExampleInf() {
	fmt.Println(Inf(+1, +1))
	fmt.Println(Inf(+1, -1))
	fmt.Println(Inf(-1, +1))
	fmt.Println(Inf(-1, -1))
	// Output:
	// (+Inf+Infε)
	// (+Inf-Infε)
	// (-Inf+Infε)
	// (-Inf-Infε)
}

func TestIsNaN(t *testing.T) {
	var tests = []struct {
		z    *Dual
		want bool
	}{
		{zero, false},
		{e0, false},
		{e1, false},
		{New(math.NaN(), 4), true},
		{New(math.Inf(0), math.NaN()), false},
	}
	for _, test := range tests {
		if got := test.z.IsNaN(); got != test.want {
			t.Errorf("IsNaN(%v) = %v", test.z, got)
		}
	}
}

func ExampleNaN() {
	fmt.Println(NaN())
	// Output:
	// (NaN+NaNε)
}

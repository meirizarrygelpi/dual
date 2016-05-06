// Copyright (c) 2016 Melvin Eloy Irizarry-Gelpí
// Licenced under the MIT License.

package dual

import (
	"fmt"
	"math"
	"testing"
)

var (
	zeroR = &Real{0, 0}
	oneR  = &Real{1, 0}
	epsiR = &Real{0, 1}
)

func TestRealString(t *testing.T) {
	var tests = []struct {
		x    *Real
		want string
	}{
		{zeroR, "(0+0ε)"},
		{oneR, "(1+0ε)"},
		{epsiR, "(0+1ε)"},
		{&Real{1, 1}, "(1+1ε)"},
		{&Real{1, -1}, "(1-1ε)"},
		{&Real{-1, 1}, "(-1+1ε)"},
		{&Real{-1, -1}, "(-1-1ε)"},
	}
	for _, test := range tests {
		if got := test.x.String(); got != test.want {
			t.Errorf("String(%v) = %v, want %v", test.x, got, test.want)
		}
	}
}

func TestRealEquals(t *testing.T) {
	var tests = []struct {
		x    *Real
		y    *Real
		want bool
	}{
		{zeroR, zeroR, true},
		{oneR, oneR, true},
		{epsiR, epsiR, true},
		{oneR, epsiR, false},
		{epsiR, oneR, false},
		{&Real{2.03, 3}, &Real{2.0299999999, 3}, true},
		{&Real{1, 2}, &Real{3, 4}, false},
	}
	for _, test := range tests {
		if got := test.x.Equals(test.y); got != test.want {
			t.Errorf("Equals(%v, %v) = %v", test.x, test.y, got)
		}
	}
}

func TestRealCopy(t *testing.T) {
	var tests = []struct {
		x    *Real
		want *Real
	}{
		{zeroR, zeroR},
		{&Real{1, 2}, &Real{1, 2}},
	}
	for _, test := range tests {
		if got := new(Real).Copy(test.x); !got.Equals(test.want) {
			t.Errorf("Copy(%v) = %v, want %v", test.x, got, test.want)
		}
	}
}

func ExampleNewReal() {
	fmt.Println(NewReal(1, 0))
	fmt.Println(NewReal(0, 1))
	fmt.Println(NewReal(2, -3))
	fmt.Println(NewReal(-4, 5))
	// Output:
	// (1+0ε)
	// (0+1ε)
	// (2-3ε)
	// (-4+5ε)
}

func TestRealScal(t *testing.T) {
	var tests = []struct {
		z    *Real
		a    float64
		want *Real
	}{
		{zeroR, 1, zeroR},
		{&Real{1, 2}, 3, &Real{3, 6}},
		{&Real{1, 2}, 0, zeroR},
	}
	for _, test := range tests {
		if got := new(Real).Scal(test.z, test.a); !got.Equals(test.want) {
			t.Errorf("Scal(%v, %v) = %v, want %v",
				test.z, test.a, got, test.want)
		}
	}
}

func TestRealNeg(t *testing.T) {
	var tests = []struct {
		z    *Real
		want *Real
	}{
		{zeroR, zeroR},
		{oneR, &Real{-1, 0}},
		{epsiR, &Real{0, -1}},
		{&Real{3, 4}, &Real{-3, -4}},
	}
	for _, test := range tests {
		if got := new(Real).Neg(test.z); !got.Equals(test.want) {
			t.Errorf("Neg(%v) = %v, want %v",
				test.z, got, test.want)
		}
	}
}

func TestRealConj(t *testing.T) {
	var tests = []struct {
		z    *Real
		want *Real
	}{
		{zeroR, zeroR},
		{oneR, oneR},
		{epsiR, &Real{0, -1}},
		{&Real{3, 4}, &Real{3, -4}},
	}
	for _, test := range tests {
		if got := new(Real).Conj(test.z); !got.Equals(test.want) {
			t.Errorf("Conj(%v) = %v, want %v",
				test.z, got, test.want)
		}
	}
}

func TestRealAdd(t *testing.T) {
	var tests = []struct {
		x    *Real
		y    *Real
		want *Real
	}{
		{zeroR, zeroR, zeroR},
		{oneR, oneR, &Real{2, 0}},
		{epsiR, epsiR, &Real{0, 2}},
		{oneR, epsiR, &Real{1, 1}},
		{epsiR, oneR, &Real{1, 1}},
	}
	for _, test := range tests {
		if got := new(Real).Add(test.x, test.y); !got.Equals(test.want) {
			t.Errorf("Add(%v, %v) = %v, want %v",
				test.x, test.y, got, test.want)
		}
	}
}

func TestRealSub(t *testing.T) {
	var tests = []struct {
		x    *Real
		y    *Real
		want *Real
	}{
		{zeroR, zeroR, zeroR},
		{oneR, oneR, zeroR},
		{epsiR, epsiR, zeroR},
		{oneR, epsiR, &Real{1, -1}},
		{epsiR, oneR, &Real{-1, 1}},
	}
	for _, test := range tests {
		if got := new(Real).Sub(test.x, test.y); !got.Equals(test.want) {
			t.Errorf("Sub(%v, %v) = %v, want %v",
				test.x, test.y, got, test.want)
		}
	}
}

func TestRealMul(t *testing.T) {
	var tests = []struct {
		x    *Real
		y    *Real
		want *Real
	}{
		{zeroR, zeroR, zeroR},
		{oneR, oneR, oneR},
		{epsiR, epsiR, zeroR},
		{oneR, epsiR, epsiR},
		{epsiR, oneR, epsiR},
	}
	for _, test := range tests {
		if got := new(Real).Mul(test.x, test.y); !got.Equals(test.want) {
			t.Errorf("Mul(%v, %v) = %v, want %v",
				test.x, test.y, got, test.want)
		}
	}
}

func TestRealQuad(t *testing.T) {
	var tests = []struct {
		z    *Real
		want float64
	}{
		{zeroR, 0},
		{oneR, 1},
		{epsiR, 0},
		{&Real{-2, 1}, 4},
	}
	for _, test := range tests {
		if got := test.z.Quad(); notEquals(got, test.want) {
			t.Errorf("Quad(%v) = %v, want %v",
				test.z, got, test.want)
		}
	}
}

func TestRealIsZeroDiv(t *testing.T) {
	var tests = []struct {
		z    *Real
		want bool
	}{
		{zeroR, true},
		{oneR, false},
		{epsiR, true},
	}
	for _, test := range tests {
		if got := test.z.IsZeroDiv(); got != test.want {
			t.Errorf("IsZeroDiv(%v) = %v", test.z, got)
		}
	}
}

func TestRealInv(t *testing.T) {
	var tests = []struct {
		x    *Real
		want *Real
	}{
		{oneR, oneR},
		{&Real{2, 0}, &Real{0.5, 0}},
	}
	for _, test := range tests {
		if got := new(Real).Inv(test.x); !got.Equals(test.want) {
			t.Errorf("Inv(%v) = %v, want %v",
				test.x, got, test.want)
		}
	}
}

func TestRealQuo(t *testing.T) {
	var tests = []struct {
		x    *Real
		y    *Real
		want *Real
	}{
		{oneR, oneR, oneR},
		{&Real{0.5, 0}, &Real{2, 0}, &Real{0.25, 0}},
	}
	for _, test := range tests {
		if got := new(Real).Quo(test.x, test.y); !got.Equals(test.want) {
			t.Errorf("Quo(%v, %v) = %v, want %v",
				test.x, test.y, got, test.want)
		}
	}
}

func TestRealIsInf(t *testing.T) {
	var tests = []struct {
		z    *Real
		want bool
	}{
		{zeroR, false},
		{oneR, false},
		{epsiR, false},
		{&Real{math.Inf(0), 4}, true},
	}
	for _, test := range tests {
		if got := test.z.IsInf(); got != test.want {
			t.Errorf("IsInf(%v) = %v", test.z, got)
		}
	}
}

func ExampleRealInf() {
	fmt.Println(RealInf(+1, +1))
	fmt.Println(RealInf(+1, -1))
	fmt.Println(RealInf(-1, +1))
	fmt.Println(RealInf(-1, -1))
	// Output:
	// (+Inf+Infε)
	// (+Inf-Infε)
	// (-Inf+Infε)
	// (-Inf-Infε)
}

func TestRealIsNaN(t *testing.T) {
	var tests = []struct {
		z    *Real
		want bool
	}{
		{zeroR, false},
		{oneR, false},
		{epsiR, false},
		{&Real{math.NaN(), 4}, true},
		{&Real{math.Inf(0), math.NaN()}, false},
	}
	for _, test := range tests {
		if got := test.z.IsNaN(); got != test.want {
			t.Errorf("IsNaN(%v) = %v", test.z, got)
		}
	}
}

func ExampleRealNaN() {
	fmt.Println(RealNaN())
	// Output:
	// (NaN+NaNε)
}

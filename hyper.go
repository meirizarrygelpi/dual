package dual

import (
	"fmt"
	"math"
	"strings"
)

// A Hyper represents a dual hyper number as an ordered array of two
// pointers to Real values.
type Hyper [2]*Real

var (
	// Symbols for the canonical dual hyper basis.
	symbHyper = [4]string{"", "ε", "η", "εη"}
)

// String returns the string representation of a Hyper value. If z
// corresponds to the dual hyper number a + bε + cη + dεη, then the string is
// "(a+bε+cη+dεη)", similar to complex128 values.
func (z *Hyper) String() string {
	v := make([]float64, 4)
	v[0], v[1] = (z[0])[0], (z[0])[1]
	v[2], v[3] = (z[1])[0], (z[1])[1]
	a := make([]string, 9)
	a[0] = "("
	a[1] = fmt.Sprintf("%g", v[0])
	i := 1
	for j := 2; j < 8; j = j + 2 {
		switch {
		case math.Signbit(v[i]):
			a[j] = fmt.Sprintf("%g", v[i])
		case math.IsInf(v[i], +1):
			a[j] = "+Inf"
		default:
			a[j] = fmt.Sprintf("+%g", v[i])
		}
		a[j+1] = symbHyper[i]
		i++
	}
	a[8] = ")"
	return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Hyper) Equals(y *Hyper) bool {
	for i := range z {
		if !z[i].Equals(y[i]) {
			return false
		}
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Hyper) Copy(y *Hyper) *Hyper {
	for i, v := range y {
		z[i] = new(Real).Copy(v)
	}
	return z
}

// NewHyper returns a pointer to a Hyper value made from four given float64
// values.
func NewHyper(a, b, c, d float64) *Hyper {
	z := new(Hyper)
	z[0] = NewReal(a, b)
	z[1] = NewReal(c, d)
	return z
}

// IsHyperInf returns true if any of the components of z are infinite.
func (z *Hyper) IsHyperInf() bool {
	for _, v := range z {
		if v.IsRealInf() {
			return true
		}
	}
	return false
}

// HyperInf returns a pointer to a dual hyper infinity value.
func HyperInf(a, b, c, d int) *Hyper {
	z := new(Hyper)
	z[0] = RealInf(a, b)
	z[1] = RealInf(c, d)
	return z
}

// IsHyperNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Hyper) IsPerplexNaN() bool {
	for _, v := range z {
		if v.IsRealInf() {
			return false
		}
	}
	for _, v := range z {
		if v.IsRealNaN() {
			return true
		}
	}
	return false
}

// HyperNaN returns a pointer to a dual hyper NaN value.
func HyperNaN() *Hyper {
	z := new(Hyper)
	z[0] = RealNaN()
	z[1] = RealNaN()
	return z
}

// Scal sets z equal to y scaled by a (with a being a Real pointer),
// and returns z.
//
// This is a special case of Mul:
// 		Scal(y, a) = Mul(y, Hyper{a, 0})
func (z *Hyper) Scal(y *Hyper, a *Real) *Hyper {
	for i, v := range y {
		z[i] = new(Real).Mul(v, a)
	}
	return z
}

// Dil sets z equal to the dilation of y by a, and returns z.
//
// This is a special case of Mul:
// 		Dil(y, a) = Mul(y, Hyper{Real{a, 0}, 0})
func (z *Hyper) Dil(y *Hyper, a float64) *Hyper {
	for i, v := range y {
		z[i] = new(Real).Scal(v, a)
	}
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Hyper) Neg(y *Hyper) *Hyper {
	return z.Dil(y, -1)
}

// DualConj sets z equal to the dual hyper conjugate of y, and returns z.
func (z *Hyper) DualConj(y *Hyper) *Hyper {
	z[0] = new(Real).Copy(y[0])
	z[1] = new(Real).Neg(y[1])
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Hyper) Add(x, y *Hyper) *Hyper {
	for i, v := range x {
		z[i] = new(Real).Add(v, y[i])
	}
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Hyper) Sub(x, y *Hyper) *Hyper {
	for i, v := range x {
		z[i] = new(Real).Sub(v, y[i])
	}
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The basic multiplication rules are:
//      ε * ε = η * η = 0
//      ε * η = η * ε = εη
//      εη * εη = 0
//      ε * εη = εη * ε = 0
//      η * εη = εη * η = 0
// This multiplication rule is commutative and associative.
func (z *Hyper) Mul(x, y *Hyper) *Hyper {
	p := new(Hyper).Copy(x)
	q := new(Hyper).Copy(y)
	z[0] = new(Real).Mul(p[0], q[0])
	z[1] = new(Real).Add(
		new(Real).Mul(p[0], q[1]),
		new(Real).Mul(p[1], q[0]),
	)
	return z
}

// DualQuad returns the dual quadrance of z, a Real value.
func (z *Hyper) DualQuad() *Real {
	return (new(Hyper).Mul(z, new(Hyper).DualConj(z)))[0]
}

package dual

import (
	"fmt"
	"math"
	"strings"

	"github.com/meirizarrygelpi/split"
)

// A Perplex represents a dual perplex number as an ordered array of two
// pointers to split.Complex values.
type Perplex [2]*split.Complex

var (
	// Symbols for the canonical dual perplex basis.
	symbPerplex = [4]string{"", "s", "ε", "εs"}
)

// String returns the string representation of a Perplex value.
//
// If z corresponds to the dual perplex number a + bs + cε + dεs, then the
// string is "(a+bs+cε+dεs)", similar to complex128 values.
func (z *Perplex) String() string {
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
		a[j+1] = symbPerplex[i]
		i++
	}
	a[8] = ")"
	return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Perplex) Equals(y *Perplex) bool {
	if !z[0].Equals(y[0]) {
		return false
	}
	if !z[1].Equals(y[1]) {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Perplex) Copy(y *Perplex) *Perplex {
	z[0] = new(split.Complex).Copy(y[0])
	z[1] = new(split.Complex).Copy(y[1])
	return z
}

// NewPerplex returns a pointer to a Perplex value made from four given float64
// values.
func NewPerplex(a, b, c, d float64) *Perplex {
	z := new(Perplex)
	z[0] = split.New(a, b)
	z[1] = split.New(c, d)
	return z
}

// IsPerplexInf returns true if any of the components of z are infinite.
func (z *Perplex) IsPerplexInf() bool {
	if z[0].IsInf() || z[1].IsInf() {
		return true
	}
	return false
}

// PerplexInf returns a pointer to a dual perplex infinity value.
func PerplexInf(a, b, c, d int) *Perplex {
	z := new(Perplex)
	z[0] = split.Inf(a, b)
	z[1] = split.Inf(c, d)
	return z
}

// IsPerplexNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Perplex) IsPerplexNaN() bool {
	if z[0].IsInf() || z[1].IsInf() {
		return false
	}
	if z[0].IsNaN() || z[1].IsNaN() {
		return true
	}
	return false
}

// PerplexNaN returns a pointer to a dual perplex NaN value.
func PerplexNaN() *Perplex {
	z := new(Perplex)
	z[0] = split.NaN()
	z[1] = split.NaN()
	return z
}

// Scal sets z equal to y scaled by a (with a being a split.Complex pointer),
// and returns z.
//
// This is a special case of Mul:
// 		Scal(y, a) = Mul(y, Perplex{a, 0})
func (z *Perplex) Scal(y *Perplex, a *split.Complex) *Perplex {
	z[0] = new(split.Complex).Mul(y[0], a)
	z[1] = new(split.Complex).Mul(y[1], a)
	return z
}

// Dil sets z equal to the dilation of y by a, and returns z.
//
// This is a special case of Mul:
// 		Dil(y, a) = Mul(y, Perplex{split.Complex{a, 0}, 0})
func (z *Perplex) Dil(y *Perplex, a float64) *Perplex {
	z[0] = new(split.Complex).Scal(y[0], a)
	z[1] = new(split.Complex).Scal(y[1], a)
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Perplex) Neg(y *Perplex) *Perplex {
	return z.Dil(y, -1)
}

// Conj sets z equal to the split-complex conjugate of y, and returns z.
func (z *Perplex) Conj(y *Perplex) *Perplex {
	z[0] = new(split.Complex).Conj(y[0])
	z[1] = new(split.Complex).Conj(y[1])
	return z
}

// DualConj sets z equal to the dual conjugate of y, and returns z.
func (z *Perplex) DualConj(y *Perplex) *Perplex {
	z[0] = new(split.Complex).Copy(y[0])
	z[1] = new(split.Complex).Neg(y[1])
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Perplex) Add(x, y *Perplex) *Perplex {
	z[0] = new(split.Complex).Add(x[0], y[0])
	z[1] = new(split.Complex).Add(x[1], y[1])
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Perplex) Sub(x, y *Perplex) *Perplex {
	z[0] = new(split.Complex).Sub(x[0], y[0])
	z[1] = new(split.Complex).Sub(x[1], y[1])
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The basic multiplication rules are:
//      ε * ε = 0
//      s * s = +1
//      εs * εs = 0
//      ε * s = s * ε = εs
//      εs * s = s * εs = +ε
//      ε * εs = εs * ε = 0
// This multiplication rule is commutative and associative.
func (z *Perplex) Mul(x, y *Perplex) *Perplex {
	p := new(Perplex).Copy(x)
	q := new(Perplex).Copy(y)
	z[0] = new(split.Complex).Mul(p[0], q[0])
	z[1] = new(split.Complex).Add(
		new(split.Complex).Mul(p[0], q[1]),
		new(split.Complex).Mul(p[1], q[0]),
	)
	return z
}

// Quad returns the quadrance of z, a pointer to a Real value.
func (z *Perplex) Quad() *Real {
	p := new(Perplex).Mul(z, new(Perplex).Conj(z))
	r := new(Real)
	r[0] = (p[0])[0]
	r[1] = (p[1])[0]
	return r
}

// DualQuad returns the dual quadrance of z, a split.Complex value.
func (z *Perplex) DualQuad() *split.Complex {
	return (new(Perplex).Mul(z, new(Perplex).DualConj(z)))[0]
}

// IsZeroDiv returns true if z is a zero divisor. This is equivalent to
// z being nilpotent (i.e. z² = 0).
func (z *Perplex) IsZeroDiv() bool {
	return !z[0].Equals(&split.Complex{0, 0})
}

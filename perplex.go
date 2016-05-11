// Copyright (c) 2016 Melvin Eloy Irizarry-Gelpí
// Licenced under the MIT License.

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

// Real returns the real part of z, a pointer to a split.Complex value.
func (z *Perplex) Real() *split.Complex {
	return z[0]
}

// Dual returns the dual part of z, a pointer to a split.Complex value.
func (z *Perplex) Dual() *split.Complex {
	return z[1]
}

// SetReal sets the real part of z equal to a.
func (z *Perplex) SetReal(a *split.Complex) {
	z[0] = a
}

// SetDual sets the dual part of z equal to b.
func (z *Perplex) SetDual(b *split.Complex) {
	z[1] = b
}

// Cartesian returns the four Cartesian components of z.
func (z *Perplex) Cartesian() (a, b, c, d float64) {
	a, b = z.Real().Cartesian()
	c, d = z.Dual().Cartesian()
	return
}

// String returns the string representation of a Perplex value.
//
// If z corresponds to the dual perplex number a + bs + cε + dεs, then the
// string is "(a+bs+cε+dεs)", similar to complex128 values.
func (z *Perplex) String() string {
	v := make([]float64, 4)
	v[0], v[1], v[2], v[3] = z.Cartesian()
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
	if !z.Real().Equals(y.Real()) || !z.Dual().Equals(y.Dual()) {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Perplex) Copy(y *Perplex) *Perplex {
	z.SetReal(new(split.Complex).Copy(y.Real()))
	z.SetDual(new(split.Complex).Copy(y.Dual()))
	return z
}

// NewPerplex returns a pointer to a Perplex value made from four given float64
// values.
func NewPerplex(a, b, c, d float64) *Perplex {
	z := new(Perplex)
	z.SetReal(split.New(a, b))
	z.SetDual(split.New(c, d))
	return z
}

// IsInf returns true if any of the components of z are infinite.
func (z *Perplex) IsInf() bool {
	if z.Real().IsInf() || z.Dual().IsInf() {
		return true
	}
	return false
}

// Inf sets z equal to a dual perplex infinity value.
func (z *Perplex) Inf(a, b, c, d int) *Perplex {
	z.Real().Inf(a, b)
	z.Dual().Inf(c, d)
	return z
}

// IsNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Perplex) IsNaN() bool {
	if z.Real().IsInf() || z.Dual().IsInf() {
		return false
	}
	if z.Real().IsNaN() || z.Dual().IsNaN() {
		return true
	}
	return false
}

// NaN sets z equal to a dual perplex NaN value.
func (z *Perplex) NaN() *Perplex {
	z.Real().NaN()
	z.Dual().NaN()
	return z
}

// Scal sets z equal to y scaled by a (with a being a split.Complex pointer),
// and returns z.
//
// This is a special case of Mul:
// 		Scal(y, a) = Mul(y, Perplex{a, 0})
func (z *Perplex) Scal(y *Perplex, a *split.Complex) *Perplex {
	z.Real().Mul(y.Real(), a)
	z.Dual().Mul(y.Dual(), a)
	return z
}

// Dil sets z equal to the dilation of y by a, and returns z.
//
// This is a special case of Mul:
// 		Dil(y, a) = Mul(y, Perplex{split.Complex{a, 0}, 0})
func (z *Perplex) Dil(y *Perplex, a float64) *Perplex {
	z.Real().Scal(y.Real(), a)
	z.Dual().Scal(y.Dual(), a)
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Perplex) Neg(y *Perplex) *Perplex {
	return z.Dil(y, -1)
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Perplex) Conj(y *Perplex) *Perplex {
	z.Real().Conj(y.Real())
	z.Dual().Neg(y.Dual())
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Perplex) Add(x, y *Perplex) *Perplex {
	z.Real().Add(x.Real(), y.Real())
	z.Dual().Add(x.Dual(), y.Dual())
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Perplex) Sub(x, y *Perplex) *Perplex {
	z.Real().Sub(x.Real(), y.Real())
	z.Dual().Sub(x.Dual(), y.Dual())
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
// This multiplication rule is noncommutative but associative.
func (z *Perplex) Mul(x, y *Perplex) *Perplex {
	p := new(Perplex).Copy(x)
	q := new(Perplex).Copy(y)
	z.Real().Mul(p.Real(), q.Real())
	z.Dual().Mul(q.Dual(), p.Real())
	z.Dual().Add(z.Dual(), q.Dual().Mul(p.Dual(), q.Real().Conj(q.Real())))
	return z
}

// Commutator sets z equal to the commutator of x and y, and returns z.
func (z *Perplex) Commutator(x, y *Perplex) *Perplex {
	z.Mul(x, y)
	return z.Sub(z, new(Perplex).Mul(y, x))
}

// Quad returns the quadrance of z, a float64 value.
func (z *Perplex) Quad() float64 {
	return z.Real().Quad()
}

// Copyright (c) 2016 Melvin Eloy Irizarry-Gelpí
// Licenced under the MIT License.

package dual

import (
	"fmt"
	"math"
	"strings"
)

// A Real represents a dual real number.
type Real [2]float64

// Real returns the real part of z, a float64 value.
func (z *Real) Real() float64 {
	return z[0]
}

// Dual returns the dual part of z, a float64 value.
func (z *Real) Dual() float64 {
	return z[1]
}

// SetReal sets the real part of z equal to a.
func (z *Real) SetReal(a float64) {
	z[0] = a
}

// SetDual sets the dual part of z equal to b.
func (z *Real) SetDual(b float64) {
	z[1] = b
}

// Cartesian returns the two Cartesian components of z.
func (z *Real) Cartesian() (a, b float64) {
	a = z.Real()
	b = z.Dual()
	return
}

// String returns the string version of a Real value.
//
// If z = a + bε, then the string is "(a+bε)", similar to complex128 values.
func (z *Real) String() string {
	a := make([]string, 5)
	a[0] = "("
	a[1] = fmt.Sprintf("%g", z.Real())
	switch {
	case math.Signbit(z.Dual()):
		a[2] = fmt.Sprintf("%g", z.Dual())
	case math.IsInf(z.Dual(), +1):
		a[2] = "+Inf"
	default:
		a[2] = fmt.Sprintf("+%g", z.Dual())
	}
	a[3] = "ε"
	a[4] = ")"
	return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Real) Equals(y *Real) bool {
	if notEquals(z.Real(), y.Real()) || notEquals(z.Dual(), y.Dual()) {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Real) Copy(y *Real) *Real {
	z.SetReal(y.Real())
	z.SetDual(y.Dual())
	return z
}

// NewReal returns a pointer to a Real value made from two given float64 values.
func NewReal(a, b float64) *Real {
	z := new(Real)
	z.SetReal(a)
	z.SetDual(b)
	return z
}

// IsInf returns true if any of the components of z are infinite.
func (z *Real) IsInf() bool {
	if math.IsInf(z.Real(), 0) || math.IsInf(z.Dual(), 0) {
		return true
	}
	return false
}

// RealInf returns a pointer to a dual real infinity value.
func RealInf(a, b int) *Real {
	z := new(Real)
	z.SetReal(math.Inf(a))
	z.SetDual(math.Inf(b))
	return z
}

// IsNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Real) IsNaN() bool {
	if math.IsInf(z.Real(), 0) || math.IsInf(z.Dual(), 0) {
		return false
	}
	if math.IsNaN(z.Real()) || math.IsNaN(z.Dual()) {
		return true
	}
	return false
}

// RealNaN returns a pointer to a dual real NaN value.
func RealNaN() *Real {
	nan := math.NaN()
	z := new(Real)
	z.SetReal(nan)
	z.SetDual(nan)
	return z
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Real) Scal(y *Real, a float64) *Real {
	z.SetReal(y.Real() * a)
	z.SetDual(y.Dual() * a)
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Real) Neg(y *Real) *Real {
	return z.Scal(y, -1)
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Real) Conj(y *Real) *Real {
	z.SetReal(y.Real())
	z.SetDual(y.Dual() * -1)
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Real) Add(x, y *Real) *Real {
	z.SetReal(x.Real() + y.Real())
	z.SetDual(x.Dual() + y.Dual())
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Real) Sub(x, y *Real) *Real {
	z.SetReal(x.Real() - y.Real())
	z.SetDual(x.Dual() - y.Dual())
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The basic rule is:
// 		ε * ε = 0
// This multiplication operation is commutative and associative.
func (z *Real) Mul(x, y *Real) *Real {
	p := new(Real).Copy(x)
	q := new(Real).Copy(y)
	z.SetReal(p.Real() * q.Real())
	z.SetDual((p.Real() * q.Dual()) + (p.Dual() * q.Real()))
	return z
}

// Quad returns the non-negative dual quadrance of z, a float64 value.
func (z *Real) Quad() float64 {
	return z.Real() * z.Real()
}

// IsZeroDiv returns true if z is a zero divisor. This is equivalent to
// z being nilpotent (i.e. z² = 0).
func (z *Real) IsZeroDiv() bool {
	return !notEquals(z.Real(), 0)
}

// Inv sets z equal to the inverse of y, and returns z. If y is a zero divisor,
// then Inv panics.
func (z *Real) Inv(y *Real) *Real {
	if y.IsZeroDiv() {
		panic("zero divisor")
	}
	return z.Scal(new(Real).Conj(y), 1/y.Quad())
}

// Quo sets z equal to the quotient of x and y, and returns z. If y is a zero
// divisor, then Quo panics.
func (z *Real) Quo(x, y *Real) *Real {
	if y.IsZeroDiv() {
		panic("zero divisor denominator")
	}
	return z.Scal(new(Real).Mul(x, new(Real).Conj(y)), 1/y.Quad())
}

// Sin sets z equal to the dual sine of y, and returns z.
func (z *Real) Sin(y *Real) *Real {
	s, c := math.Sincos(y.Real())
	z.SetReal(s)
	z.SetDual(y.Dual() * c)
	return z
}

// Cos sets z equal to the dual cosine of y, and returns z.
func (z *Real) Cos(y *Real) *Real {
	s, c := math.Sincos(y.Real())
	z.SetReal(c)
	z.SetDual(y.Dual() * s * -1)
	return z
}

// Exp sets z equal to the dual exponential of y, and returns z.
func (z *Real) Exp(y *Real) *Real {
	e := math.Exp(y.Real())
	z.SetReal(e)
	z.SetDual(y.Dual() * e)
	return z
}

// Sinh sets z equal to the dual hyperbolic sine of y, and returns z.
func (z *Real) Sinh(y *Real) *Real {
	z.SetReal(math.Sinh(y.Real()))
	z.SetDual(y.Dual() * math.Cosh(y.Real()))
	return z
}

// Cosh sets z equal to the dual hyperbolic cosine of y, and returns z.
func (z *Real) Cosh(y *Real) *Real {
	z.SetReal(math.Cosh(y.Real()))
	z.SetDual(y.Dual() * math.Sinh(y.Real()))
	return z
}

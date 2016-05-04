package dual

import (
	"fmt"
	"math"
	"strings"
)

// A Real represents a dual real number as an ordered array of two float64
// values.
type Real [2]float64

// String returns the string version of a Real value.
//
// If z = a + bε, then the string is "(a+bε)", similar to complex128 values.
func (z *Real) String() string {
	a := make([]string, 5)
	a[0] = "("
	a[1] = fmt.Sprintf("%g", z[0])
	switch {
	case math.Signbit(z[1]):
		a[2] = fmt.Sprintf("%g", z[1])
	case math.IsInf(z[1], +1):
		a[2] = "+Inf"
	default:
		a[2] = fmt.Sprintf("+%g", z[1])
	}
	a[3] = "ε"
	a[4] = ")"
	return strings.Join(a, "")
}

// Equals returns true if z and y are equal.
func (z *Real) Equals(y *Real) bool {
	if notEquals(z[0], y[0]) || notEquals(z[1], y[1]) {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Real) Copy(y *Real) *Real {
	z[0] = y[0]
	z[1] = y[1]
	return z
}

// NewReal returns a pointer to a Real value made from two given float64 values.
func NewReal(a, b float64) *Real {
	z := new(Real)
	z[0] = a
	z[1] = b
	return z
}

// IsRealInf returns true if any of the components of z are infinite.
func (z *Real) IsRealInf() bool {
	if math.IsInf(z[0], 0) || math.IsInf(z[1], 0) {
		return true
	}
	return false
}

// RealInf returns a pointer to a dual real infinity value.
func RealInf(a, b int) *Real {
	z := new(Real)
	z[0] = math.Inf(a)
	z[1] = math.Inf(b)
	return z
}

// IsRealNaN returns true if any component of z is NaN and neither is an
// infinity.
func (z *Real) IsRealNaN() bool {
	if math.IsInf(z[0], 0) || math.IsInf(z[1], 0) {
		return false
	}
	if math.IsNaN(z[0]) || math.IsNaN(z[1]) {
		return true
	}
	return false
}

// RealNaN returns a pointer to a dual real NaN value.
func RealNaN() *Real {
	nan := math.NaN()
	return &Real{nan, nan}
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Real) Scal(y *Real, a float64) *Real {
	z[0] = y[0] * a
	z[1] = y[1] * a
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Real) Neg(y *Real) *Real {
	return z.Scal(y, -1)
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Real) Conj(y *Real) *Real {
	z[0] = +y[0]
	z[1] = -y[1]
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Real) Add(x, y *Real) *Real {
	z[0] = x[0] + y[0]
	z[1] = x[1] + y[1]
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Real) Sub(x, y *Real) *Real {
	z[0] = x[0] - y[0]
	z[1] = x[1] - y[1]
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
	z[0] = p[0] * q[0]
	z[1] = (p[0] * q[1]) + (p[1] * q[0])
	return z
}

// Quad returns the non-negative dual quadrance of z, a float64 value.
func (z *Real) Quad() float64 {
	return z[0] * z[0]
}

// IsZeroDiv returns true if z is a zero divisor. This is equivalent to
// z being nilpotent (i.e. z² = 0).
func (z *Real) IsZeroDiv() bool {
	return !notEquals(z[0], 0)
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
	s, c := math.Sincos(y[0])
	z[0] = s
	z[1] = y[1] * c
	return z
}

// Cos sets z equal to the dual cosine of y, and returns z.
func (z *Real) Cos(y *Real) *Real {
	s, c := math.Sincos(y[0])
	z[0] = c
	z[1] = -y[1] * s
	return z
}

// Exp sets z equal to the dual exponential of y, and returns z.
func (z *Real) Exp(y *Real) *Real {
	e := math.Exp(y[0])
	z[0] = e
	z[1] = y[1] * e
	return z
}

// Sinh sets z equal to the dual hyperbolic sine of y, and returns z.
func (z *Real) Sinh(y *Real) *Real {
	z[0] = math.Sinh(y[0])
	z[1] = y[1] * math.Cosh(y[0])
	return z
}

// Cosh sets z equal to the dual hyperbolic cosine of y, and returns z.
func (z *Real) Cosh(y *Real) *Real {
	z[0] = math.Cosh(y[0])
	z[1] = y[1] * math.Sinh(y[0])
	return z
}

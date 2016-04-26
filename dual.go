package dual

import (
	"fmt"
	"math"
	"strings"
)

// A Dual represents a dual number as an ordered array of two float64 values.
type Dual [2]float64

// String returns the string version of a Dual value. If z = a + bε, then the
// string is "(a+bε)", similar to complex128 values.
func (z *Dual) String() string {
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
func (z *Dual) Equals(y *Dual) bool {
	for i := range z {
		if notEquals(z[i], y[i]) {
			return false
		}
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Dual) Copy(y *Dual) *Dual {
	for i, v := range y {
		z[i] = v
	}
	return z
}

// New returns a pointer to a Dual value made from two given float64 values.
func New(a, b float64) *Dual {
	z := new(Dual)
	z[0] = a
	z[1] = b
	return z
}

// IsInf returns true if any of the components of z are infinite.
func (z *Dual) IsInf() bool {
	for _, v := range z {
		if math.IsInf(v, 0) {
			return true
		}
	}
	return false
}

// Inf returns a pointer to a dual infinity value.
func Inf(a, b int) *Dual {
	return New(math.Inf(a), math.Inf(b))
}

// IsNaN returns true if any component of z is NaN and neither is an infinity.
func (z *Dual) IsNaN() bool {
	for _, v := range z {
		if math.IsInf(v, 0) {
			return false
		}
	}
	for _, v := range z {
		if math.IsNaN(v) {
			return true
		}
	}
	return false
}

// NaN returns a pointer to a dual NaN value.
func NaN() *Dual {
	nan := math.NaN()
	return New(nan, nan)
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Dual) Scal(y *Dual, a float64) *Dual {
	for i, v := range y {
		z[i] = a * v
	}
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Dual) Neg(y *Dual) *Dual {
	return z.Scal(y, -1)
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Dual) Conj(y *Dual) *Dual {
	z[0] = +y[0]
	z[1] = -y[1]
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Dual) Add(x, y *Dual) *Dual {
	for i, v := range x {
		z[i] = v + y[i]
	}
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Dual) Sub(x, y *Dual) *Dual {
	for i, v := range x {
		z[i] = v - y[i]
	}
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
func (z *Dual) Mul(x, y *Dual) *Dual {
	p := new(Dual).Copy(x)
	q := new(Dual).Copy(y)
	z[0] = p[0] * q[0]
	z[1] = (p[0] * q[1]) + (p[1] * q[0])
	return z
}

// Quad returns the non-negative quadrance of z.
func (z *Dual) Quad() float64 {
	return (new(Dual).Mul(z, new(Dual).Conj(z)))[0]
}

// IsZeroDiv returns true if z is a zero divisor. This is equivalent to
// z being nilpotent (i.e. z² = 0).
func (z *Dual) IsZeroDiv() bool {
	return !notEquals(z[0], 0)
}

// Inv sets z equal to the inverse of y, and returns z. If y is a
// zero divisor, then Inv panics.
func (z *Dual) Inv(y *Dual) *Dual {
	if y.IsZeroDiv() {
		panic("zero divisor")
	}
	return z.Scal(new(Dual).Conj(y), 1/y.Quad())
}

// Quo sets z equal to the quotient of x and y, and returns z. If y is a
// zero divisor, then Quo panics.
func (z *Dual) Quo(x, y *Dual) *Dual {
	if y.IsZeroDiv() {
		panic("denominator is a zero divisor")
	}
	return z.Scal(new(Dual).Mul(x, new(Dual).Conj(y)), 1/y.Quad())
}

// Sin sets z equal to the dual sine of y, and returns z.
func (z *Dual) Sin(y *Dual) *Dual {
	s, c := math.Sincos(y[0])
	z[0] = s
	z[1] = y[1] * c
	return z
}

// Cos sets z equal to the dual cosine of y, and returns z.
func (z *Dual) Cos(y *Dual) *Dual {
	s, c := math.Sincos(y[0])
	z[0] = c
	z[1] = -y[1] * s
	return z
}

// Exp sets z equal to the dual exponential of y, and returns z.
func (z *Dual) Exp(y *Dual) *Dual {
	e := math.Exp(y[0])
	z[0] = e
	z[1] = y[1] * e
	return z
}

// Sinh sets z equal to the dual hyperbolic sine of y, and returns z.
func (z *Dual) Sinh(y *Dual) *Dual {
	z[0] = math.Sinh(y[0])
	z[1] = y[1] * math.Cosh(y[0])
	return z
}

// Cosh sets z equal to the dual hyperbolic cosine of y, and returns z.
func (z *Dual) Cosh(y *Dual) *Dual {
	z[0] = math.Cosh(y[0])
	z[1] = y[1] * math.Sinh(y[0])
	return z
}

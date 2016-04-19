package dual

import (
	"fmt"
	"math"
	"strings"
)

// Dual type represents a dual number a + bε over the real numbers, with
// ε² = 0.
type Dual [2]float64

// String method returns the string version of a Dual value. If z = a + bε,
// then the string is "(a+bε)", similar to complex128 values.
func (z *Dual) String() string {
	a := make([]string, 5)
	a[0] = "("
	a[1] = fmt.Sprintf("%g", z[0])
	switch {
	case math.IsInf(z[1], +1):
		a[2] = "+Inf"
	case z[1] < 0:
		a[2] = fmt.Sprintf("%g", z[1])
	default:
		a[2] = fmt.Sprintf("+%g", z[1])
	}
	a[3] = "ε"
	a[4] = ")"
	return strings.Join(a, "")
}

// Equals method returns true if z and x are equal.
func (z *Dual) Equals(x *Dual) bool {
	for i := range z {
		if notEquals(z[i], x[i]) {
			return false
		}
	}
	return true
}

// Copy method copies x onto z, and returns z.
func (z *Dual) Copy(x *Dual) *Dual {
	for i, v := range x {
		z[i] = v
	}
	return z
}

// New function returns a pointer to a Dual value made from two given real
// numbers (i.e. float64s).
func New(a, b float64) *Dual {
	z := new(Dual)
	z[0] = a
	z[1] = b
	return z
}

// Scal method sets z equal to x scaled by a, and returns z.
func (z *Dual) Scal(x *Dual, a float64) *Dual {
	for i, v := range x {
		z[i] = a * v
	}
	return z
}

// Neg method sets z equal to the negative of x, and returns z.
func (z *Dual) Neg(x *Dual) *Dual {
	return z.Scal(x, -1)
}

// Conj method sets z equal to the conjugate of x, and returns z.
func (z *Dual) Conj(x *Dual) *Dual {
	z[0] = +x[0]
	z[1] = -x[1]
	return z
}

// Add method sets z equal to the sum of x and y, and returns z.
func (z *Dual) Add(x, y *Dual) *Dual {
	for i, v := range x {
		z[i] = v + y[i]
	}
	return z
}

// Sub method sets z equal to the difference of x and y, and returns z.
func (z *Dual) Sub(x, y *Dual) *Dual {
	for i, v := range x {
		z[i] = v - y[i]
	}
	return z
}

// Mul method sets z equal to the product of x and y, and returns z.
func (z *Dual) Mul(x, y *Dual) *Dual {
	p := new(Dual).Copy(x)
	q := new(Dual).Copy(y)
	z[0] = p[0] * q[0]
	z[1] = (p[0] * q[1]) + (p[1] * q[0])
	return z
}

// Quad method returns the non-negative quadrance of z.
func (z *Dual) Quad() float64 {
	return (new(Dual).Mul(z, new(Dual).Conj(z)))[0]
}

// IsZeroDiv method returns true if z is a zero divisor. This is equivalent to
// z being nilpotent (i.e. z² = 0).
func (z *Dual) IsZeroDiv() bool {
	return !notEquals(z[0], 0)
}

// Inv method sets z equal to the inverse of x, and returns z. If x is a
// zero divisor, then Inv panics.
func (z *Dual) Inv(x *Dual) *Dual {
	if x.IsZeroDiv() {
		panic("zero divisor")
	}
	return z.Scal(new(Dual).Conj(x), 1/x.Quad())
}

// Quo method sets z equal to the quotient of x and y, and returns z. If y is a
// zero divisor, then Quo panics.
func (z *Dual) Quo(x, y *Dual) *Dual {
	if y.IsZeroDiv() {
		panic("denominator is a zero divisor")
	}
	return z.Scal(new(Dual).Mul(x, new(Dual).Conj(y)), 1/y.Quad())
}

// IsInf method returns true if any of the components of z are infinite.
func (z *Dual) IsInf() bool {
	for _, v := range z {
		if math.IsInf(v, 0) {
			return true
		}
	}
	return false
}

// Inf function returns a pointer to a dual infinity value.
func Inf(a, b int) *Dual {
	return New(math.Inf(a), math.Inf(b))
}

// IsNaN method returns true if any component of z is NaN and neither is an
// infinity.
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

// NaN function returns a pointer to a dual NaN value.
func NaN() *Dual {
	nan := math.NaN()
	return New(nan, nan)
}
